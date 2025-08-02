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

func convertAttrToField(attr slog.Attr) zapcore.Field {
	if attr.Equal(slog.Attr{}) {
		// Ignore empty attrs.
		return zap.Skip()
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindTime:
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		if attr.Key == "" {
			// Inlines recursively.
			return zap.Inline(groupObject(attr.Value.Group()))
		}
		return zap.Object(attr.Key, groupObject(attr.Value.Group()))
	case slog.KindLogValuer:
		return convertAttrToField(slog.Attr{
			Key: attr.Key,
			// TODO: resolve the value in a lazy way.
			// This probably needs a new Zap field type
			// that can be resolved lazily.
			Value: attr.Value.Resolve(),
		})
	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}

// convertSlogLevel maps slog Levels to zap Levels.
// Note that there is some room between slog levels while zap levels are continuous, so we can't 1:1 map them.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func convertSlogLevel(l slog.Level) zapcore.Level {
	switch {
	case l >= slog.LevelError:
		return zapcore.ErrorLevel
	case l >= slog.LevelWarn:
		return zapcore.WarnLevel
	case l >= slog.LevelInfo:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.core.Enabled(convertSlogLevel(level))
}

// Handle handles the Record.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	ent := zapcore.Entry{
		Level:      convertSlogLevel(record.Level),
		Time:       record.Time,
		Message:    record.Message,
		LoggerName: h.name,
	}
	ce := h.core.Check(ent, nil)
	if ce == nil {
		return nil
	}

	if h.addCaller && record.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
		if frame.PC != 0 {
			ce.Caller = zapcore.EntryCaller{
				Defined:  true,
				PC:       frame.PC,
				File:     frame.File,
				Line:     frame.Line,
				Function: frame.Function,
			}
		}
	}

	if record.Level >= h.addStackAt {
		// Skipping 3:
		// zapslog/handler log/slog.(*Logger).log
		// slog/logger log/slog.(*Logger).log
		// slog/logger log/slog.(*Logger).<level>
		ce.Stack = stacktrace.Take(3 + h.callerSkip)
	}

	fields := make([]zapcore.Field, 0, record.NumAttrs()+len(h.groups))

	var addedNamespace bool
	record.Attrs(func(attr slog.Attr) bool {
		f := convertAttrToField(attr)
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() {
			// Namespaces are added only if at least one field is present
			// to avoid creating empty groups.
			fields = h.appendGroups(fields)
			addedNamespace = true
		}
		fields = append(fields, f)
		return true
	})

	ce.Write(fields...)
	return nil
}

func (h *Handler) appendGroups(fields []zapcore.Field) []zapcore.Field {
	for _, g := range h.groups {
		fields = append(fields, zap.Namespace(g))
	}
	return fields
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zapcore.Field, 0, len(attrs)+len(h.groups))
	var addedNamespace bool
	for _, attr := range attrs {
		f := convertAttrToField(attr)
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() {
			// Namespaces are added only if at least one field is present
			// to avoid creating empty groups.
			fields = h.appendGroups(fields)
			addedNamespace = true
		}
		fields = append(fields, f)
	}

	cloned := *h
	cloned.core = h.core.With(fields)
	if addedNamespace {
		// These groups have been applied so we can clear them.
		cloned.groups = nil
	}
	return &cloned
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *Handler) WithGroup(group string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = group

	cloned := *h
	cloned.groups = newGroups
	return &cloned
}
