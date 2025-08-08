package log

import (
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option struct {
	//zap
	AppName           string        //日志文件前缀
	ErrorFileName     string        // error 级别日志文件名
	WarnFileName      string        // error 级别日志文件名
	NormalFileName    string        // 非 error 级别日志文件名
	Level             zapcore.Level //日志等级
	Development       bool          //是否是开发模式
	Encoding          string        // 日志编码
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
	Dir string
}

type LogOption func(opts *Option)

func WithDir(dir string) LogOption {
	return func(opts *Option) {
		opts.Dir = dir
	}
}

type Logger struct {
	*zap.Logger
	sync.RWMutex
	initialized bool
}

func NewLogger(opts ...LogOption) *zap.Logger {
	l := &Logger{}
	l.Lock()
	defer l.Unlock()
	if l.initialized {
		return l.Logger
	}

	opt := Option{}
	for _, o := range opts {
		o(&opt)
	}

	config := zap.NewProductionConfig()
	logger, _ := config.Build()
	l.Logger = logger
	l.initialized = true
	return logger
}

func TestLogger(t *testing.T) {

}
