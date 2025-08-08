package log

import (
	"fmt"
	"io"
	"testing"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/petermattis/goid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 测试zap 的 ProductionConfig 配置
func TestZapProductionConfig(t *testing.T) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.Encoding = "json"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, _ := config.Build()
	logger.Info("info")
	logger.Debug("debug")
	logger.Warn("warn")
	logger.Error("error")
	logger.Fatal("fatal")
	logger.Panic("panic")
}

// 测试zap 的各种配置参数，看看如何影响日志输出
func TestConfigCustom(t *testing.T) {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		SkipLineEnding: false, // 改为false，确保不跳过行结束符
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		NewReflectedEncoder: func(io.Writer) zapcore.ReflectedEncoder {
			panic("TODO")
		},
	}
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{},
	}

	logger, _ := config.Build()
	logger.Sugar().Info("info")
	logger.Sugar().Debug("debug")
	logger.Sugar().Warn("warn")
	logger.Sugar().Error("error")
	logger.Sugar().Fatal("fatal")
	logger.Sugar().Panic("panic")
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
		enc.AppendString(fmt.Sprintf("[%d]%s", goid.Get(), caller.TrimmedPath()))
	}
	if opt.Encoding == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

// 普通业务日志使用lumberjack来按大小压缩即可
func getLumberJackWriteSyncer(opt Option) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   opt.FileName,   // 日志文件名
		MaxSize:    opt.MaxSize,    // 文件最大大小（MB），超过后切割
		MaxBackups: opt.MaxBackups, // 保留的备份文件数量，0表示保留所有
		MaxAge:     opt.MaxAge,     // 文件最大保留天数（约20年）
		Compress:   opt.Compress,   // 是否压缩旧的日志文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// lumberjack 按时间压缩有问题，所以使用file-rotatelogs来按时间压缩
// 按时间压缩的某些场景比如 data log 等需要按天统计，所以使用file-rotatelogs更合适
func getFileRotateWriteSyncer(opt Option) zapcore.WriteSyncer {
	logs, err := rotatelogs.New(
		opt.FileName,
		rotatelogs.WithLinkName(opt.FileName),
		rotatelogs.WithMaxAge(time.Duration(opt.RotationMaxAge)*time.Hour*24),
		rotatelogs.WithRotationTime(time.Duration(opt.RotationTime)*time.Hour),
		rotatelogs.WithRotationSize(opt.RotationSize),
		rotatelogs.WithRotationCount(uint(opt.RotationCount)),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(logs)
}

func getZapCore(opt Option) zap.Option {
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.WarnLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})

	infocore := zapcore.NewCore(getZapEncoder(opt), getLumberJackWriteSyncer(opt), opt.Level)
	errorcore := zapcore.NewCore(getZapEncoder(opt), getLumberJackWriteSyncer(opt), errorPriority)
	warncore := zapcore.NewCore(getZapEncoder(opt), getLumberJackWriteSyncer(opt), warnPriority)
	tee := zapcore.NewTee(infocore, errorcore, warncore)

	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(tee)
	})
}

func getLogger(opt Option) *zap.Logger {
	core := getZapCore(opt)
	config := getZapConfig(opt)

	logger, err := config.Build(core, zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	return logger
}

func TestConfigCustom2(t *testing.T) {

	logger := getLogger(Option{
		Level: zap.InfoLevel,
	})
	logger.Sugar().Info("info")
	logger.Sugar().Debug("debug")
	logger.Sugar().Warn("warn")
	logger.Sugar().Error("error")
	logger.Sugar().Fatal("fatal")
	logger.Sugar().Panic("panic")

}
