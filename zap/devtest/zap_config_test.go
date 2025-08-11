package log

import (
	"strconv"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func TestConfigCustom2(t *testing.T) {

	logger := NewLogger(
		WithDebug(true),
		WithEncoding("json"),
		WithAppName("read_go_code"),
		WithDir("logs"),
		WithFileName("./logs/test1.log"),

		WithMaxSize(1),
		WithMaxBackups(2),
		WithMaxAge(1),
		WithCompress(true),
		//
		WithRotationTime(1),
		WithRotationSize(1),
		WithRotationCount(2),
		WithRotationMaxAge(3),
	)
	for i := 0; i < 100000000; i++ {
		logger.Sugar().Debug("debug " + strconv.Itoa(i))
		logger.Sugar().Info("info " + strconv.Itoa(i))
		logger.Sugar().Warn("warn " + strconv.Itoa(i))
		logger.Sugar().Error("error " + strconv.Itoa(i))
		logger.Sugar().Info("================================================")
		time.Sleep(1 * time.Second)
	}
	// logger.Sugar().Fatal("fatal")
	// logger.Sugar().Panic("panic")

}

func main() {
	logger := NewLogger(
	// WithDebug(false),
	// WithEncoding("json"),
	// WithAppName("read_go_code"),
	// WithDir("logs"),
	// WithFileName("./logs/test1.log"),

	// WithMaxSize(1),
	// WithMaxBackups(2),
	// WithMaxAge(1),
	// WithCompress(true),
	// //
	// WithRotationTime(1),
	// WithRotationSize(1),
	// WithRotationCount(2),
	// WithRotationMaxAge(3),
	)
	for i := 0; i < 100000000; i++ {
		logger.Sugar().Debug("debug " + strconv.Itoa(i))
		logger.Sugar().Info("info " + strconv.Itoa(i))
		logger.Sugar().Warn("warn " + strconv.Itoa(i))
		logger.Sugar().Error("error " + strconv.Itoa(i))
		logger.Sugar().Info("================================================")
		time.Sleep(1 * time.Second)
	}
	// logger.Sugar().Fatal("fatal")
	// logger.Sugar().Panic("panic")
}
