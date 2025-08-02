// Copyright (c) 2016 Uber Technologies, Inc.
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

package zap // zap包：提供快速、结构化、分级日志记录

import (
	"fmt"    // fmt包：格式化I/O
	"io"     // io包：基本I/O原语
	"os"     // os包：操作系统接口
	"strings" // strings包：字符串操作

	"go.uber.org/zap/internal/bufferpool" // 内部缓冲池
	"go.uber.org/zap/internal/stacktrace" // 内部堆栈跟踪
	"go.uber.org/zap/zapcore"             // zapcore包：核心接口和实现
)

// A Logger provides fast, leveled, structured logging. All methods are safe
// for concurrent use.
//
// The Logger is designed for contexts in which every microsecond and every
// allocation matters, so its API intentionally favors performance and type
// safety over brevity. For most applications, the SugaredLogger strikes a
// better balance between performance and ergonomics.
// Logger提供快速、分级、结构化的日志记录。所有方法都是并发安全的。
//
// Logger设计用于每微秒和每次内存分配都很重要的场景，
// 因此它的API有意优先考虑性能和类型安全，而不是简洁性。
// 对于大多数应用程序，SugaredLogger在性能和易用性之间取得了更好的平衡。
type Logger struct {
	core zapcore.Core // 核心Core实例，负责实际的日志处理

	development bool                      // 开发模式标志
	addCaller   bool                      // 是否添加调用者信息
	onPanic     zapcore.CheckWriteHook    // Panic时的钩子函数，默认是WriteThenPanic
	onFatal     zapcore.CheckWriteHook    // Fatal时的钩子函数，默认是WriteThenFatal

	name        string                    // Logger的名称
	errorOutput zapcore.WriteSyncer       // 错误输出写入器

	addStack zapcore.LevelEnabler         // 堆栈添加级别启用器

	callerSkip int                        // 调用者跳过的栈帧数

	clock zapcore.Clock                   // 时钟接口，用于获取时间戳
}

// New constructs a new Logger from the provided zapcore.Core and Options. If
// the passed zapcore.Core is nil, it falls back to using a no-op
// implementation.
//
// This is the most flexible way to construct a Logger, but also the most
// verbose. For typical use cases, the highly-opinionated presets
// (NewProduction, NewDevelopment, and NewExample) or the Config struct are
// more convenient.
//
// For sample code, see the package-level AdvancedConfiguration example.
// New从提供的zapcore.Core和Options构造一个新的Logger。
// 如果传入的zapcore.Core为nil，它会回退到使用无操作实现。
//
// 这是构造Logger最灵活的方式，但也是最冗长的方式。
// 对于典型用例，高度集成的预设（NewProduction、NewDevelopment和NewExample）
// 或Config结构体更加方便。
//
// 示例代码请参见包级别的AdvancedConfiguration示例。
func New(core zapcore.Core, options ...Option) *Logger {
	if core == nil {                    // 如果core为nil
		return NewNop()                 // 返回无操作Logger
	}
	log := &Logger{                     // 创建新的Logger实例
		core:        core,              // 设置核心Core
		errorOutput: zapcore.Lock(os.Stderr), // 设置错误输出为标准错误（带锁）
		addStack:    zapcore.FatalLevel + 1,  // 设置堆栈添加级别（Fatal级别以上）
		clock:       zapcore.DefaultClock,    // 设置默认时钟
	}
	return log.WithOptions(options...)  // 应用选项并返回
}

// NewNop returns a no-op Logger. It never writes out logs or internal errors,
// and it never runs user-defined hooks.
//
// Using WithOptions to replace the Core or error output of a no-op Logger can
// re-enable logging.
// NewNop返回一个无操作的Logger。它从不写出日志或内部错误，
// 也从不运行用户定义的钩子函数。
//
// 使用WithOptions替换无操作Logger的Core或错误输出可以重新启用日志记录。
func NewNop() *Logger {
	return &Logger{
		core:        zapcore.NewNopCore(),         // 设置无操作Core
		errorOutput: zapcore.AddSync(io.Discard), // 设置错误输出为丢弃器
		addStack:    zapcore.FatalLevel + 1,      // 设置堆栈添加级别
		clock:       zapcore.DefaultClock,        // 设置默认时钟
	}
}

// NewProduction builds a sensible production Logger that writes InfoLevel and
// above logs to standard error as JSON.
//
// It's a shortcut for NewProductionConfig().Build(...Option).
// NewProduction构建一个合理的生产环境Logger，它将InfoLevel及以上级别的
// 日志以JSON格式写入标准错误。
//
// 这是NewProductionConfig().Build(...Option)的快捷方式。
func NewProduction(options ...Option) (*Logger, error) {
	return NewProductionConfig().Build(options...) // 使用生产配置构建Logger
}

// NewDevelopment builds a development Logger that writes DebugLevel and above
// logs to standard error in a human-friendly format.
//
// It's a shortcut for NewDevelopmentConfig().Build(...Option).
// NewDevelopment构建一个开发环境Logger，它将DebugLevel及以上级别的
// 日志以人类友好的格式写入标准错误。
//
// 这是NewDevelopmentConfig().Build(...Option)的快捷方式。
func NewDevelopment(options ...Option) (*Logger, error) {
	return NewDevelopmentConfig().Build(options...) // 使用开发配置构建Logger
}

// Must is a helper that wraps a call to a function returning (*Logger, error)
// and panics if the error is non-nil. It is intended for use in variable
// initialization such as:
//
//	var logger = zap.Must(zap.NewProduction())
// Must是一个辅助函数，它包装对返回(*Logger, error)的函数的调用，
// 如果错误非nil则触发panic。它旨在用于变量初始化，例如：
//
//	var logger = zap.Must(zap.NewProduction())
func Must(logger *Logger, err error) *Logger {
	if err != nil {    // 如果存在错误
		panic(err)     // 触发panic
	}

	return logger      // 返回logger
}

// NewExample builds a Logger that's designed for use in zap's testable
// examples. It writes DebugLevel and above logs to standard out as JSON, but
// omits the timestamp and calling function to keep example output
// short and deterministic.
// NewExample构建一个专为zap的可测试示例而设计的Logger。
// 它将DebugLevel及以上级别的日志以JSON格式写入标准输出，
// 但省略时间戳和调用函数以保持示例输出简短且确定性。
func NewExample(options ...Option) *Logger {
	encoderCfg := zapcore.EncoderConfig{                    // 编码器配置
		MessageKey:     "msg",                              // 消息键名
		LevelKey:       "level",                            // 级别键名
		NameKey:        "logger",                           // 名称键名
		EncodeLevel:    zapcore.LowercaseLevelEncoder,      // 小写级别编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,         // ISO8601时间编码器
		EncodeDuration: zapcore.StringDurationEncoder,      // 字符串持续时间编码器
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, DebugLevel) // 创建Core
	return New(core).WithOptions(options...)                                           // 创建Logger并应用选项
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
// Sugar包装Logger以提供更符合人体工程学但稍微较慢的API。
// 给Logger加糖的成本很低，因此单个应用程序同时使用Logger和SugaredLogger
// 是合理的，在性能敏感代码的边界之间进行转换。
func (log *Logger) Sugar() *SugaredLogger {
	core := log.clone()        // 克隆当前logger
	core.callerSkip += 2       // 增加调用者跳过数（因为多了一层包装）
	return &SugaredLogger{core} // 返回SugaredLogger实例
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
// Named向logger的名称添加新的路径段。段之间用句点连接。
// 默认情况下，Logger是未命名的。
func (log *Logger) Named(s string) *Logger {
	if s == "" {               // 如果名称为空
		return log             // 直接返回原logger
	}
	l := log.clone()           // 克隆logger
	if log.name == "" {        // 如果原名称为空
		l.name = s             // 直接设置名称
	} else {                   // 如果原名称不为空
		l.name = strings.Join([]string{l.name, s}, ".") // 用句点连接名称
	}
	return l                   // 返回新logger
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
// WithOptions克隆当前Logger，应用提供的Options，
// 并返回结果Logger。并发使用是安全的。
func (log *Logger) WithOptions(opts ...Option) *Logger {
	c := log.clone()           // 克隆logger
	for _, opt := range opts { // 遍历选项
		opt.apply(c)           // 应用选项到克隆的logger
	}
	return c                   // 返回修改后的logger
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa. Any fields that
// require evaluation (such as Objects) are evaluated upon invocation of With.
// With创建一个子logger并向其添加结构化上下文。
// 添加到子logger的字段不会影响父logger，反之亦然。
// 任何需要求值的字段（如Objects）在调用With时被求值。
func (log *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {      // 如果没有字段
		return log             // 直接返回原logger
	}
	l := log.clone()           // 克隆logger
	l.core = l.core.With(fields) // 将字段添加到core
	return l                   // 返回新logger
}

// WithLazy creates a child logger and adds structured context to it lazily.
//
// The fields are evaluated only if the logger is further chained with [With]
// or is written to with any of the log level methods.
// Until that occurs, the logger may retain references to objects inside the fields,
// and logging will reflect the state of an object at the time of logging,
// not the time of WithLazy().
//
// WithLazy provides a worthwhile performance optimization for contextual loggers
// when the likelihood of using the child logger is low,
// such as error paths and rarely taken branches.
//
// Similar to [With], fields added to the child don't affect the parent, and vice versa.
// WithLazy创建一个子logger并懒惰地向其添加结构化上下文。
//
// 字段只有在logger进一步与[With]链式调用或使用任何日志级别方法写入时才会被求值。
// 在此之前，logger可能保留对字段内对象的引用，
// 日志记录将反映对象在记录时的状态，而不是WithLazy()时的状态。
//
// 当使用子logger的可能性较低时（如错误路径和很少采用的分支），
// WithLazy为上下文logger提供了有价值的性能优化。
//
// 与[With]类似，添加到子logger的字段不会影响父logger，反之亦然。
func (log *Logger) WithLazy(fields ...Field) *Logger {
	if len(fields) == 0 {      // 如果没有字段
		return log             // 直接返回原logger
	}
	return log.WithOptions(WrapCore(func(core zapcore.Core) zapcore.Core { // 包装Core
		return zapcore.NewLazyWith(core, fields) // 创建懒惰求值的Core
	}))
}

// Level reports the minimum enabled level for this logger.
//
// For NopLoggers, this is [zapcore.InvalidLevel].
// Level报告此logger的最小启用级别。
//
// 对于NopLogger，这是[zapcore.InvalidLevel]。
func (log *Logger) Level() zapcore.Level {
	return zapcore.LevelOf(log.core) // 从core获取级别
}

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
// Check如果在指定级别记录消息被启用，则返回CheckedEntry。
// 这是一个完全可选的优化；在高性能应用程序中，
// Check可以帮助避免分配切片来保存字段。
func (log *Logger) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return log.check(lvl, msg) // 调用内部check方法
}

// Log logs a message at the specified level. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
// Any Fields that require  evaluation (such as Objects) are evaluated upon
// invocation of Log.
// Log在指定级别记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
// 任何需要求值的字段（如Objects）在调用Log时被求值。
func (log *Logger) Log(lvl zapcore.Level, msg string, fields ...Field) {
	if ce := log.check(lvl, msg); ce != nil { // 检查是否应该记录
		ce.Write(fields...)                   // 写入字段
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
// Debug在DebugLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
func (log *Logger) Debug(msg string, fields ...Field) {
	if ce := log.check(DebugLevel, msg); ce != nil { // 检查Debug级别是否启用
		ce.Write(fields...)                         // 写入字段
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
// Info在InfoLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
func (log *Logger) Info(msg string, fields ...Field) {
	if ce := log.check(InfoLevel, msg); ce != nil { // 检查Info级别是否启用
		ce.Write(fields...)                        // 写入字段
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
// Warn在WarnLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
func (log *Logger) Warn(msg string, fields ...Field) {
	if ce := log.check(WarnLevel, msg); ce != nil { // 检查Warn级别是否启用
		ce.Write(fields...)                        // 写入字段
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
// Error在ErrorLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
func (log *Logger) Error(msg string, fields ...Field) {
	if ce := log.check(ErrorLevel, msg); ce != nil { // 检查Error级别是否启用
		ce.Write(fields...)                         // 写入字段
	}
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
// DPanic在DPanicLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
//
// 如果logger处于开发模式，它会触发panic（DPanic意味着"开发panic"）。
// 这对于捕获可恢复但不应该发生的错误很有用。
func (log *Logger) DPanic(msg string, fields ...Field) {
	if ce := log.check(DPanicLevel, msg); ce != nil { // 检查DPanic级别是否启用
		ce.Write(fields...)                          // 写入字段
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
// Panic在PanicLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
//
// 然后logger会触发panic，即使PanicLevel的日志记录被禁用。
func (log *Logger) Panic(msg string, fields ...Field) {
	if ce := log.check(PanicLevel, msg); ce != nil { // 检查Panic级别是否启用
		ce.Write(fields...)                         // 写入字段
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
// Fatal在FatalLevel记录消息。消息包括在日志记录点传递的任何字段，
// 以及在logger上累积的任何字段。
//
// 然后logger会调用os.Exit(1)，即使FatalLevel的日志记录被禁用。
func (log *Logger) Fatal(msg string, fields ...Field) {
	if ce := log.check(FatalLevel, msg); ce != nil { // 检查Fatal级别是否启用
		ce.Write(fields...)                         // 写入字段
	}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
// Sync调用底层Core的Sync方法，刷新任何缓冲的日志条目。
// 应用程序应该注意在退出前调用Sync。
func (log *Logger) Sync() error {
	return log.core.Sync() // 调用core的Sync方法
}

// Core returns the Logger's underlying zapcore.Core.
// Core返回Logger的底层zapcore.Core。
func (log *Logger) Core() zapcore.Core {
	return log.core // 返回core实例
}

// Name returns the Logger's underlying name,
// or an empty string if the logger is unnamed.
// Name返回Logger的底层名称，如果logger未命名则返回空字符串。
func (log *Logger) Name() string {
	return log.name // 返回logger名称
}

func (log *Logger) clone() *Logger { // clone方法创建Logger的副本
	clone := *log    // 值拷贝Logger结构体
	return &clone    // 返回新实例的指针
}

func (log *Logger) check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	// Logger.check must always be called directly by a method in the
	// Logger interface (e.g., Check, Info, Fatal).
	// This skips Logger.check and the Info/Fatal/Check/etc. method that
	// called it.
	const callerSkipOffset = 2

	// Check the level first to reduce the cost of disabled log calls.
	// Since Panic and higher may exit, we skip the optimization for those levels.
	if lvl < zapcore.DPanicLevel && !log.core.Enabled(lvl) {
		return nil
	}

	// Create basic checked entry thru the core; this will be non-nil if the
	// log message will actually be written somewhere.
	ent := zapcore.Entry{
		LoggerName: log.name,
		Time:       log.clock.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce := log.core.Check(ent, nil)
	willWrite := ce != nil

	// Set up any required terminal behavior.
	switch ent.Level {
	case zapcore.PanicLevel:
		ce = ce.After(ent, terminalHookOverride(zapcore.WriteThenPanic, log.onPanic))
	case zapcore.FatalLevel:
		ce = ce.After(ent, terminalHookOverride(zapcore.WriteThenFatal, log.onFatal))
	case zapcore.DPanicLevel:
		if log.development {
			ce = ce.After(ent, terminalHookOverride(zapcore.WriteThenPanic, log.onPanic))
		}
	}

	// Only do further annotation if we're going to write this message; checked
	// entries that exist only for terminal behavior don't benefit from
	// annotation.
	if !willWrite {
		return ce
	}

	// Thread the error output through to the CheckedEntry.
	ce.ErrorOutput = log.errorOutput

	addStack := log.addStack.Enabled(ce.Level)
	if !log.addCaller && !addStack {
		return ce
	}

	// Adding the caller or stack trace requires capturing the callers of
	// this function. We'll share information between these two.
	stackDepth := stacktrace.First
	if addStack {
		stackDepth = stacktrace.Full
	}
	stack := stacktrace.Capture(log.callerSkip+callerSkipOffset, stackDepth)
	defer stack.Free()

	if stack.Count() == 0 {
		if log.addCaller {
			_, _ = fmt.Fprintf(
				log.errorOutput,
				"%v Logger.check error: failed to get caller\n",
				ent.Time.UTC(),
			)
			_ = log.errorOutput.Sync()
		}
		return ce
	}

	frame, more := stack.Next()

	if log.addCaller {
		ce.Caller = zapcore.EntryCaller{
			Defined:  frame.PC != 0,
			PC:       frame.PC,
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		}
	}

	if addStack {
		buffer := bufferpool.Get()
		defer buffer.Free()

		stackfmt := stacktrace.NewFormatter(buffer)

		// We've already extracted the first frame, so format that
		// separately and defer to stackfmt for the rest.
		stackfmt.FormatFrame(frame)
		if more {
			stackfmt.FormatStack(stack)
		}
		ce.Stack = buffer.String()
	}

	return ce
}

func terminalHookOverride(defaultHook, override zapcore.CheckWriteHook) zapcore.CheckWriteHook {
	// A nil or WriteThenNoop hook will lead to continued execution after
	// a Panic or Fatal log entry, which is unexpected. For example,
	//
	//   f, err := os.Open(..)
	//   if err != nil {
	//     log.Fatal("cannot open", zap.Error(err))
	//   }
	//   fmt.Println(f.Name())
	//
	// The f.Name() will panic if we continue execution after the log.Fatal.
	if override == nil || override == zapcore.WriteThenNoop {
		return defaultHook
	}
	return override
}
