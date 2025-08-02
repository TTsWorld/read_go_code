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
	"errors" // errors包：错误处理
	"sort"   // sort包：排序工具
	"time"   // time包：时间处理

	"go.uber.org/zap/zapcore" // zapcore包：核心接口和实现
)

// SamplingConfig sets a sampling strategy for the logger. Sampling caps the
// global CPU and I/O load that logging puts on your process while attempting
// to preserve a representative subset of your logs.
//
// If specified, the Sampler will invoke the Hook after each decision.
//
// Values configured here are per-second. See zapcore.NewSamplerWithOptions for
// details.
// SamplingConfig为logger设置采样策略。采样限制了日志记录对进程造成的
// 全局CPU和I/O负载，同时尝试保留日志的代表性子集。
//
// 如果指定，Sampler将在每次决策后调用Hook。
//
// 这里配置的值是每秒的。详细信息请参见zapcore.NewSamplerWithOptions。
type SamplingConfig struct {
	Initial    int                                           `json:"initial" yaml:"initial"`       // 每秒初始允许的日志数量
	Thereafter int                                           `json:"thereafter" yaml:"thereafter"` // 超过Initial后，每N条消息采样一次
	Hook       func(zapcore.Entry, zapcore.SamplingDecision) `json:"-" yaml:"-"`                   // 采样决策后的回调钩子
}

// Config offers a declarative way to construct a logger. It doesn't do
// anything that can't be done with New, Options, and the various
// zapcore.WriteSyncer and zapcore.Core wrappers, but it's a simpler way to
// toggle common options.
//
// Note that Config intentionally supports only the most common options. More
// unusual logging setups (logging to network connections or message queues,
// splitting output between multiple files, etc.) are possible, but require
// direct use of the zapcore package. For sample code, see the package-level
// BasicConfiguration and AdvancedConfiguration examples.
//
// For an example showing runtime log level changes, see the documentation for
// AtomicLevel.
// Config提供声明式构造logger的方式。它不能做任何用New、Options
// 和各种zapcore.WriteSyncer及zapcore.Core包装器做不到的事情，
// 但它是切换常见选项的更简单方式。
//
// 注意Config有意只支持最常见的选项。更不寻常的日志设置
// （记录到网络连接或消息队列，在多个文件间分割输出等）是可能的，
// 但需要直接使用zapcore包。示例代码请参见包级别的
// BasicConfiguration和AdvancedConfiguration示例。
//
// 有关运行时日志级别更改的示例，请参见AtomicLevel的文档。
type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	// Level是最小启用的日志级别。注意这是动态级别，
	// 调用Config.Level.SetLevel将原子地更改从此配置派生的所有logger的日志级别。
	Level AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	// Development将logger置于开发模式，这会改变DPanicLevel的行为
	// 并更自由地获取堆栈跟踪。
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	// DisableCaller停止用调用函数的文件名和行号注释日志。
	// 默认情况下，所有日志都被注释。
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	// DisableStacktrace完全禁用自动堆栈跟踪捕获。
	// 默认情况下，开发模式下WarnLevel及以上级别、生产模式下ErrorLevel及以上级别
	// 会捕获堆栈跟踪。
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	// Sampling设置采样策略。nil SamplingConfig禁用采样。
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	// Encoding设置logger的编码。有效值为"json"和"console"，
	// 以及通过RegisterEncoder注册的任何第三方编码。
	Encoding string `json:"encoding" yaml:"encoding"`
	// EncoderConfig sets options for the chosen encoder. See
	// zapcore.EncoderConfig for details.
	// EncoderConfig为选择的编码器设置选项。详细信息请参见zapcore.EncoderConfig。
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	// OutputPaths是写入日志输出的URL或文件路径列表。详细信息请参见Open。
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	// ErrorOutputPaths是写入内部logger错误的URL列表。默认是标准错误。
	//
	// 注意此设置只影响内部错误；有关将错误级别日志发送到与
	// info和debug级别日志不同位置的示例代码，
	// 请参见包级别的AdvancedConfiguration示例。
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	// InitialFields is a collection of fields to add to the root logger.
	// InitialFields是要添加到根logger的字段集合。
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

// NewProductionEncoderConfig returns an opinionated EncoderConfig for
// production environments.
//
// Messages encoded with this configuration will be JSON-formatted
// and will have the following keys by default:
//
//   - "level": The logging level (e.g. "info", "error").
//   - "ts": The current time in number of seconds since the Unix epoch.
//   - "msg": The message passed to the log statement.
//   - "caller": If available, a short path to the file and line number
//     where the log statement was issued.
//     The logger configuration determines whether this field is captured.
//   - "stacktrace": If available, a stack trace from the line
//     where the log statement was issued.
//     The logger configuration determines whether this field is captured.
//
// By default, the following formats are used for different types:
//
//   - Time is formatted as floating-point number of seconds since the Unix
//     epoch.
//   - Duration is formatted as floating-point number of seconds.
//
// You may change these by setting the appropriate fields in the returned
// object.
// For example, use the following to change the time encoding format:
//
//	cfg := zap.NewProductionEncoderConfig()
//	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewProductionConfig builds a reasonable default production logging
// configuration.
// Logging is enabled at InfoLevel and above, and uses a JSON encoder.
// Logs are written to standard error.
// Stacktraces are included on logs of ErrorLevel and above.
// DPanicLevel logs will not panic, but will write a stacktrace.
//
// Sampling is enabled at 100:100 by default,
// meaning that after the first 100 log entries
// with the same level and message in the same second,
// it will log every 100th entry
// with the same level and message in the same second.
// You may disable this behavior by setting Sampling to nil.
//
// See [NewProductionEncoderConfig] for information
// on the default encoder configuration.
func NewProductionConfig() Config {
	return Config{
		Level:       NewAtomicLevelAt(InfoLevel),
		Development: false,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// NewDevelopmentEncoderConfig returns an opinionated EncoderConfig for
// development environments.
//
// Messages encoded with this configuration will use Zap's console encoder
// intended to print human-readable output.
// It will print log messages with the following information:
//
//   - The log level (e.g. "INFO", "ERROR").
//   - The time in ISO8601 format (e.g. "2017-01-01T12:00:00Z").
//   - The message passed to the log statement.
//   - If available, a short path to the file and line number
//     where the log statement was issued.
//     The logger configuration determines whether this field is captured.
//   - If available, a stacktrace from the line
//     where the log statement was issued.
//     The logger configuration determines whether this field is captured.
//
// By default, the following formats are used for different types:
//
//   - Time is formatted in ISO8601 format (e.g. "2017-01-01T12:00:00Z").
//   - Duration is formatted as a string (e.g. "1.234s").
//
// You may change these by setting the appropriate fields in the returned
// object.
// For example, use the following to change the time encoding format:
//
//	cfg := zap.NewDevelopmentEncoderConfig()
//	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewDevelopmentConfig builds a reasonable default development logging
// configuration.
// Logging is enabled at DebugLevel and above, and uses a console encoder.
// Logs are written to standard error.
// Stacktraces are included on logs of WarnLevel and above.
// DPanicLevel logs will panic.
//
// See [NewDevelopmentEncoderConfig] for information
// on the default encoder configuration.
func NewDevelopmentConfig() Config {
	return Config{
		Level:            NewAtomicLevelAt(DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// Build constructs a logger from the Config and Options.
func (cfg Config) Build(opts ...Option) (*Logger, error) {
	enc, err := cfg.buildEncoder()
	if err != nil {
		return nil, err
	}

	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	if cfg.Level == (AtomicLevel{}) {
		return nil, errors.New("missing Level")
	}

	log := New(
		zapcore.NewCore(enc, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

func (cfg Config) buildOptions(errSink zapcore.WriteSyncer) []Option { // 构建配置选项
	opts := []Option{ErrorOutput(errSink)} // 初始化选项列表，设置错误输出

	if cfg.Development { // 如果是开发模式
		opts = append(opts, Development()) // 添加开发选项
	}

	if !cfg.DisableCaller { // 如果未禁用调用者信息
		opts = append(opts, AddCaller()) // 添加调用者选项
	}

	stackLevel := ErrorLevel // 默认堆栈跟踪级别为Error
	if cfg.Development {     // 开发模式下
		stackLevel = WarnLevel // 降低到Warn级别
	}
	if !cfg.DisableStacktrace { // 如果未禁用堆栈跟踪
		opts = append(opts, AddStacktrace(stackLevel)) // 添加堆栈跟踪选项
	}

	if scfg := cfg.Sampling; scfg != nil { // 如果配置了采样
		opts = append(opts, WrapCore(func(core zapcore.Core) zapcore.Core { // 包装Core以添加采样功能
			var samplerOpts []zapcore.SamplerOption // 采样器选项
			if scfg.Hook != nil {                   // 如果有采样钩子
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(scfg.Hook)) // 添加钩子选项
			}
			return zapcore.NewSamplerWithOptions( // 创建带选项的采样器
				core,                    // 底层Core
				time.Second,             // 采样间隔
				cfg.Sampling.Initial,    // 初始数量
				cfg.Sampling.Thereafter, // 之后数量
				samplerOpts...,          // 其他选项
			)
		})) // 添加包装的Core选项
	}

	if len(cfg.InitialFields) > 0 { // 如果有初始字段
		fs := make([]Field, 0, len(cfg.InitialFields))    // 初始化字段列表
		keys := make([]string, 0, len(cfg.InitialFields)) // 初始化键名列表
		for k := range cfg.InitialFields {                // 遍历初始字段
			keys = append(keys, k) // 收集所有键名
		}
		sort.Strings(keys)       // 对键名排序以保证确定性顺序
		for _, k := range keys { // 遍历排序后的键名
			fs = append(fs, Any(k, cfg.InitialFields[k])) // 添加字段
		}
		opts = append(opts, Fields(fs...)) // 添加字段选项
	}

	return opts // 返回所有选项
}

// openSinks opens the configured output and error sinks.
// openSinks打开配置的输出和错误接收器。
func (cfg Config) openSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := Open(cfg.OutputPaths...) // 打开输出路径
	if err != nil {                                 // 如果有错误
		return nil, nil, err // 返回错误
	}
	errSink, _, err := Open(cfg.ErrorOutputPaths...) // 打开错误输出路径
	if err != nil {                                  // 如果有错误
		closeOut()           // 关闭已打开的输出
		return nil, nil, err // 返回错误
	}
	return sink, errSink, nil // 返回两个接收器
}

// buildEncoder builds the configured encoder.
// buildEncoder构建配置的编码器。
func (cfg Config) buildEncoder() (zapcore.Encoder, error) {
	return newEncoder(cfg.Encoding, cfg.EncoderConfig) // 使用配置创建编码器
}
