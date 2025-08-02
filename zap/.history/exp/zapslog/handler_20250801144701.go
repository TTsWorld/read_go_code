// Copyright (c) 2023 Uber Technologies, Inc.
// 版权归Uber Technologies公司所有
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 特此免费授予任何获得本软件副本的人使用、复制、修改、分发等权利
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// 上述版权声明和许可声明应包含在软件的所有副本中
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
// 本软件"按原样"提供，不提供任何明示或暗示的保证

//go:build go1.21
// 构建约束：需要Go 1.21或更高版本

package zapslog // zapslog包：slog与Zap的集成

import (
	"context"    // context包：上下文管理
	"log/slog"   // slog包：结构化日志记录
	"runtime"    // runtime包：运行时服务

	"go.uber.org/zap"                   // zap包：高性能日志库
	"go.uber.org/zap/internal/stacktrace" // 内部堆栈跟踪包
	"go.uber.org/zap/zapcore"           // zapcore包：核心接口和实现
)

// Handler implements the slog.Handler by writing to a zap Core.
// Handler通过写入到zap Core来实现slog.Handler。
type Handler struct {
	core       zapcore.Core // 底层zap核心
	name       string       // 日志器名称
	addCaller  bool         // 是否添加调用者信息
	addStackAt slog.Level   // 添加堆栈跟踪的级别
	callerSkip int          // 跳过的调用者数量

	// List of unapplied groups.
	//
	// These are applied only if we encounter a real field
	// to avoid creating empty namespaces -- which is disallowed by slog's
	// usage contract.
	// 未应用的组列表。
	//
	// 这些只有在遇到真实字段时才会应用，以避免创建空命名空间
	// -- slog的使用合约不允许这样做。
	groups []string // 组名列表
}

// NewHandler builds a [Handler] that writes to the supplied [zapcore.Core]
// with options.
// NewHandler构建一个写入到提供的[zapcore.Core]的[Handler]，带有选项。
func NewHandler(core zapcore.Core, opts ...HandlerOption) *Handler {
	h := &Handler{                  // 初始化Handler
		core:       core,            // 设置核心
		addStackAt: slog.LevelError, // 默认在Error级别添加堆栈跟踪
	}
	for _, v := range opts { // 应用所有选项
		v.apply(h) // 应用选项
	}
	return h // 返回Handler
}

var _ slog.Handler = (*Handler)(nil) // 编译时检查Handler实现了slog.Handler接口

// groupObject holds all the Attrs saved in a slog.GroupValue.
// groupObject保存在slog.GroupValue中的所有Attrs。
type groupObject []slog.Attr

func (gs groupObject) MarshalLogObject(enc zapcore.ObjectEncoder) error { // 实现ObjectMarshaler接口
	for _, attr := range gs { // 遍历所有属性
		convertAttrToField(attr).AddTo(enc) // 转换属性为字段并添加到编码器
	}
	return nil // 返回nil表示成功
}

func convertAttrToField(attr slog.Attr) zapcore.Field { // 将slog属性转换为zap字段
	if attr.Equal(slog.Attr{}) { // 如果是空属性
		// Ignore empty attrs.
		// 忽略空属性。
		return zap.Skip() // 返回跳过字段
	}

	switch attr.Value.Kind() { // 根据属性值的类型进行转换
	case slog.KindBool: // 布尔类型
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration: // 持续时间类型
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64: // 64位浮点数类型
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64: // 64位整数类型
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString: // 字符串类型
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindTime: // 时间类型
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64: // 64位无符号整数类型
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup: // 组类型
		if attr.Key == "" { // 如果没有键名
			// Inlines recursively.
			// 递归内联。
			return zap.Inline(groupObject(attr.Value.Group())) // 内联组对象
		}
		return zap.Object(attr.Key, groupObject(attr.Value.Group())) // 创建组对象
	case slog.KindLogValuer: // LogValuer类型
		return convertAttrToField(slog.Attr{ // 递归转换解析后的值
			Key: attr.Key, // 保持相同的键
			// TODO: resolve the value in a lazy way.
			// This probably needs a new Zap field type
			// that can be resolved lazily.
			// TODO：以懒惰方式解析值。
			// 这可能需要一个新的Zap字段类型，可以懒惰解析。
			Value: attr.Value.Resolve(), // 解析值
		})
	default: // 默认情况
		return zap.Any(attr.Key, attr.Value.Any()) // 使用Any类型
	}
}

// convertSlogLevel maps slog Levels to zap Levels.
// Note that there is some room between slog levels while zap levels are continuous, so we can't 1:1 map them.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
// convertSlogLevel将slog级别映射到zap级别。
// 注意slog级别之间有一些间隙，而zap级别是连续的，所以我们不能1:1映射它们。
// 参见 https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func convertSlogLevel(l slog.Level) zapcore.Level {
	switch { // 根据slog级别进行分类
	case l >= slog.LevelError: // 错误级别及以上
		return zapcore.ErrorLevel // 返回错误级别
	case l >= slog.LevelWarn: // 警告级别及以上
		return zapcore.WarnLevel // 返回警告级别
	case l >= slog.LevelInfo: // 信息级别及以上
		return zapcore.InfoLevel // 返回信息级别
	default: // 其他情况（调试级别）
		return zapcore.DebugLevel // 返回调试级别
	}
}

// Enabled reports whether the handler handles records at the given level.
// Enabled报告处理程序是否处理给定级别的记录。
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.core.Enabled(convertSlogLevel(level)) // 检查转换后的级别是否启用
}

// Handle handles the Record.
// Handle处理记录。
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	ent := zapcore.Entry{ // 创建日志条目
		Level:      convertSlogLevel(record.Level), // 转换级别
		Time:       record.Time,                   // 设置时间
		Message:    record.Message,                // 设置消息
		LoggerName: h.name,                       // 设置日志器名称
	}
	ce := h.core.Check(ent, nil) // 检查是否应该记录此条目
	if ce == nil {               // 如果不应该记录
		return nil // 直接返回
	}

	if h.addCaller && record.PC != 0 { // 如果需要添加调用者信息且有PC
		frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next() // 获取调用帧信息
		if frame.PC != 0 {                                            // 如果有有效的PC
			ce.Caller = zapcore.EntryCaller{ // 设置调用者信息
				Defined:  true,          // 标记为已定义
				PC:       frame.PC,      // 程序计数器
				File:     frame.File,    // 文件名
				Line:     frame.Line,    // 行号
				Function: frame.Function, // 函数名
			}
		}
	}

	if record.Level >= h.addStackAt { // 如果级别达到需要添加堆栈的级别
		// Skipping 3:
		// zapslog/handler log/slog.(*Logger).log
		// slog/logger log/slog.(*Logger).log
		// slog/logger log/slog.(*Logger).<level>
		// 跳过3层调用：
		// zapslog/handler log/slog.(*Logger).log
		// slog/logger log/slog.(*Logger).log
		// slog/logger log/slog.(*Logger).<level>
		ce.Stack = stacktrace.Take(3 + h.callerSkip) // 捕获堆栈跟踪
	}

	fields := make([]zapcore.Field, 0, record.NumAttrs()+len(h.groups)) // 初始化字段列表

	var addedNamespace bool                              // 是否已添加命名空间
	record.Attrs(func(attr slog.Attr) bool {             // 遍历记录的所有属性
		f := convertAttrToField(attr)                     // 转换属性为字段
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() { // 如果还没添加命名空间且有组且字段非空
			// Namespaces are added only if at least one field is present
			// to avoid creating empty groups.
			// 只有在至少有一个字段存在时才添加命名空间，以避免创建空组。
			fields = h.appendGroups(fields) // 添加组
			addedNamespace = true           // 标记已添加
		}
		fields = append(fields, f) // 添加字段
		return true               // 继续迭代
	})

	ce.Write(fields...) // 写入所有字段
	return nil         // 返回nil表示成功
}

func (h *Handler) appendGroups(fields []zapcore.Field) []zapcore.Field { // 添加组到字段列表
	for _, g := range h.groups { // 遍历所有组
		fields = append(fields, zap.Namespace(g)) // 添加命名空间字段
	}
	return fields // 返回更新后的字段列表
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// WithAttrs返回一个新的Handler，其属性由接收者的属性和参数组成。
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zapcore.Field, 0, len(attrs)+len(h.groups)) // 初始化字段列表
	var addedNamespace bool                                      // 是否已添加命名空间
	for _, attr := range attrs {                                 // 遍历所有属性
		f := convertAttrToField(attr)                             // 转换属性为字段
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() { // 如果还没添加命名空间且有组且字段非空
			// Namespaces are added only if at least one field is present
			// to avoid creating empty groups.
			// 只有在至少有一个字段存在时才添加命名空间，以避免创建空组。
			fields = h.appendGroups(fields) // 添加组
			addedNamespace = true           // 标记已添加
		}
		fields = append(fields, f) // 添加字段
	}

	cloned := *h                      // 克隆Handler
	cloned.core = h.core.With(fields) // 使用新字段创建新的Core
	if addedNamespace {               // 如果已添加命名空间
		// These groups have been applied so we can clear them.
		// 这些组已经应用，所以可以清除它们。
		cloned.groups = nil // 清除组
	}
	return &cloned // 返回克隆的Handler
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// WithGroup返回一个新的Handler，将给定的组附加到接收者的现有组中。
func (h *Handler) WithGroup(group string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1) // 创建新的组列表
	copy(newGroups, h.groups)                   // 复制现有组
	newGroups[len(h.groups)] = group             // 添加新组

	cloned := *h              // 克隆Handler
	cloned.groups = newGroups // 设置新的组列表
	return &cloned            // 返回克隆的Handler
}
