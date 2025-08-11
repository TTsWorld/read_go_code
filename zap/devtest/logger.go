package log

import (
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	TRACEID = "trace_id" // TraceID 在 context 中的 key，用于请求链路追踪
)

var (
	Logger     = NewLogger(WithDebug(true)).Sugar()
	DataLogger = NewDataLogger(WithDebug(true)).Sugar()
)

func NewLogger(opts ...LogOption) *zap.Logger {

	opt := Option{}
	for _, o := range opts {
		o(&opt)
	}

	validateOption(&opt)

	config := getZapConfig(opt)
	core := getZapCore(opt)

	logger, err := config.Build(core, zap.AddCallerSkip(opt.CallerSkip))
	if err != nil {
		panic(err)
	}
	return logger
}

func NewDataLogger(opts ...LogOption) *zap.Logger {

	opt := Option{}
	for _, o := range opts {
		o(&opt)
	}
	opt.DisableCaller = true
	opt.DisableStacktrace = true

	validateOption(&opt)

	config := getZapConfig(opt)
	core := getZapDataCore(opt)
	logger, err := config.Build(core)
	if err != nil {
		panic(err)
	}

	return logger
}

func getZapConfig(opt Option) zap.Config {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(opt.Level)
	config.Development = opt.Development
	config.Encoding = opt.Encoding
	config.DisableStacktrace = opt.DisableStacktrace
	config.DisableCaller = opt.DisableCaller

	return config
}

func getZapEncoder(opt Option) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("[%d]", goid.Get()) + caller.TrimmedPath())
	}
	if opt.Debug {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if opt.Encoding == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

// 普通业务日志使用lumberjack来按大小压缩即可
func getLumberJackWriteSyncer(opt Option, level zapcore.Level) zapcore.WriteSyncer {
	if opt.Debug {
		return zapcore.AddSync(os.Stdout)
	}
	var fileName string = opt.FileName

	switch level {
	case zap.DebugLevel:
		fileName = opt.Dir + "/" + opt.AppName + "." + zap.InfoLevel.CapitalString()
	case zap.InfoLevel:
		fileName = opt.Dir + "/" + opt.AppName + "." + zap.InfoLevel.CapitalString()
	case zap.WarnLevel:
		fileName = opt.Dir + "/" + opt.AppName + "." + zap.WarnLevel.CapitalString()
	case zap.ErrorLevel:
		fileName = opt.Dir + "/" + opt.AppName + "." + zap.ErrorLevel.CapitalString()
	default:
		fileName = opt.Dir + "/" + opt.AppName + "." + zap.InfoLevel.CapitalString()
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,       // 日志文件名
		MaxSize:    opt.MaxSize,    // 文件最大大小（MB），超过后切割
		MaxBackups: opt.MaxBackups, // 保留的备份文件数量，0表示保留所有
		MaxAge:     opt.MaxAge,     // 文件最大保留天数（约20年）
		Compress:   opt.Compress,   // 是否压缩旧的日志文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// getFileRotateWriteSyncer lumberjack 按时间压缩有问题，所以使用file-rotatelogs来按时间压缩
// 按时间压缩的某些场景比如 data log 等需要按天统计，所以使用file-rotatelogs更合适
func getFileRotateWriteSyncer(opt Option) zapcore.WriteSyncer {
	if opt.Debug {
		return zapcore.AddSync(os.Stdout)
	}
	fileName := opt.Dir + "/" + opt.AppName + ".data"

	logs, err := rotatelogs.New(
		fileName,
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(opt.RotationMaxAge)*time.Hour*24),
		rotatelogs.WithRotationTime(time.Duration(opt.RotationTime)*time.Hour),
		rotatelogs.WithRotationSize(opt.RotationSize*1024*1024),
		rotatelogs.WithRotationCount(uint(opt.RotationCount)),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(logs)
}

// 获取zap core
// 1. 根据opt.Level 获取当前级别允许的级别
// 2. 根据级别获取对应日志写入到不同文件
func getZapCore(opt Option) zap.Option {
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		// 需要同时满足：1) 当前级别允许 2) 是WARN及以上级别
		return opt.Level.Enabled(lvl) && lvl >= zap.WarnLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return opt.Level.Enabled(lvl) && lvl >= zap.ErrorLevel
	})
	if opt.Debug {
		core := zapcore.NewCore(getZapEncoder(opt), zapcore.AddSync(os.Stdout), zap.DebugLevel)
		return zap.WrapCore(func(_ zapcore.Core) zapcore.Core { return core })
	}

	infocore := zapcore.NewCore(getZapEncoder(opt), getLumberJackWriteSyncer(opt, zap.InfoLevel), opt.Level)
	errorcore := zapcore.NewCore(getZapEncoder(opt), getLumberJackWriteSyncer(opt, zap.ErrorLevel), errorPriority)
	warncore := zapcore.NewCore(getZapEncoder(opt), getLumberJackWriteSyncer(opt, zap.WarnLevel), warnPriority)

	tee := zapcore.NewTee(infocore, errorcore, warncore)

	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return tee
	})
}

// 获取zap core
// 1. 根据opt.Level 获取当前级别允许的级别
// 2. 根据级别获取对应日志写入到不同文件
func getZapDataCore(opt Option) zap.Option {
	// 创建生产环境的编码器配置
	encoderConf := zap.NewProductionEncoderConfig()
	// 设置时间编码器为空函数（不输出时间）
	encoderConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
	// 设置级别编码器为空函数（不输出级别）
	encoderConf.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {}
	// 设置调用者编码器为空函数（不输出文件行号）
	encoderConf.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {}
	// 使用控制台编码器创建编码器
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(encoder, getFileRotateWriteSyncer(opt), zapcore.DebugLevel)
	})
}

func validateOption(opt *Option) {
	if opt.Dir == "" {
		opt.Dir = "./"
	}
	if opt.Debug {
		opt.Level = zap.DebugLevel
		opt.Encoding = "console"
	}
	if opt.Encoding == "" {
		opt.Encoding = "json"
	}
	if opt.AppName == "" {
		opt.AppName = "log"
	}
}
