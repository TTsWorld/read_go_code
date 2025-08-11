package log

import (
	"go.uber.org/zap/zapcore"
)

type Option struct {
	//zap
	AppName           string        //日志文件前缀
	Level             zapcore.Level //日志等级
	Development       bool          //是否是开发模式
	Encoding          string        // 日志编码, required
	DisableCaller     bool
	DisableStacktrace bool
	CallerSkip        int
	//lumberjack
	FileName   string //文件保存地方
	MaxSize    int    //日志文件小大（M）
	MaxBackups int    // 最多存在多少个切片文件
	MaxAge     int    //保存的最大天数
	Compress   bool
	//file-rotatelogs
	RotationTime   int   // 日志切割时间间隔（小时）
	RotationSize   int64 // 日志切割大小（MB）
	RotationCount  int   // 日志切割数量
	RotationMaxAge int   // 日志切割最大天数
	// 其他
	Dir   string
	Debug bool
}

type LogOption func(opts *Option)

func WithDir(dir string) LogOption {
	return func(opts *Option) {
		opts.Dir = dir
	}
}

func WithFileName(fileName string) LogOption {
	return func(opts *Option) {
		opts.FileName = fileName
	}
}

func WithMaxSize(maxSize int) LogOption {
	return func(opts *Option) {
		opts.MaxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) LogOption {
	return func(opts *Option) {
		opts.MaxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) LogOption {
	return func(opts *Option) {
		opts.MaxAge = maxAge
	}
}

func WithCompress(compress bool) LogOption {
	return func(opts *Option) {
		opts.Compress = compress
	}
}

func WithRotationTime(rotationTime int) LogOption {
	return func(opts *Option) {
		opts.RotationTime = rotationTime
	}
}

func WithRotationSize(rotationSize int64) LogOption {
	return func(opts *Option) {
		opts.RotationSize = rotationSize
	}
}

func WithRotationCount(rotationCount int) LogOption {
	return func(opts *Option) {
		opts.RotationCount = rotationCount
	}
}

func WithRotationMaxAge(rotationMaxAge int) LogOption {
	return func(opts *Option) {
		opts.RotationMaxAge = rotationMaxAge
	}
}

func WithLevel(level zapcore.Level) LogOption {
	return func(opts *Option) {
		opts.Level = level
	}
}

func WithDevelopment(development bool) LogOption {
	return func(opts *Option) {
		opts.Development = development
	}
}

func WithEncoding(encoding string) LogOption {
	return func(opts *Option) {
		opts.Encoding = encoding
	}
}

func WithDisableCaller(disableCaller bool) LogOption {
	return func(opts *Option) {
		opts.DisableCaller = disableCaller
	}
}

func WithDisableStacktrace(disableStacktrace bool) LogOption {
	return func(opts *Option) {
		opts.DisableStacktrace = disableStacktrace
	}
}

func WithCallerSkip(callerSkip int) LogOption {
	return func(opts *Option) {
		opts.CallerSkip = callerSkip
	}
}

func WithAppName(appName string) LogOption {
	return func(opts *Option) {
		opts.AppName = appName
	}
}

func WithDebug(debug bool) LogOption {
	return func(opts *Option) {
		opts.Debug = debug
	}
}
