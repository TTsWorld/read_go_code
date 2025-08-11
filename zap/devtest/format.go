package log

import (
	"context"
	"fmt"
	"os"
)

// ========================== 日志函数 ==========================

// CRITICAL 记录严重错误日志但不退出程序
// 内部调用 Panicf 但通过 recover 捕获 panic，同时记录 ERROR 级别日志
// format: 格式化字符串
// a: 格式化参数
func CRITICAL(format string, a ...interface{}) {
	// 使用 defer + recover 捕获 Panicf 产生的 panic
	defer func() {
		if err := recover(); err != nil {
			// 忽略 panic，继续执行
		}
		// 格式化消息并记录为 ERROR 级别
		datastr := fmt.Sprintf(format, a...)
		Logger.Error(datastr)
	}()
	// 调用 Panicf（会产生 panic，但被上面的 recover 捕获）
	Logger.Panicf(format, a...)
}

func Debug(a ...interface{}) {
	Logger.Debug(a...)
}

func Debugf(format string, a ...interface{}) {
	Logger.Debugf(format, a...)
}
func Info(a ...interface{}) {
	Logger.Info(a...)
}

func Infof(format string, a ...interface{}) {
	Logger.Infof(format, a...)
}

func Warn(a ...interface{}) {
	Logger.Warn(a...)
}

func Warnf(format string, a ...interface{}) {
	Logger.Warnf(format, a...)
}

func Error(a ...interface{}) {
	Logger.Error(a...)
}

func Errorf(format string, a ...interface{}) {
	Logger.Errorf(format, a...)
}

func Fatal(a ...interface{}) {
	Logger.Fatal(a...)
}

func Fatalf(format string, a ...interface{}) {
	Logger.Fatalf(format, a...)
}

func FatalfWithExit(format string, a ...interface{}) {
	Logger.Fatalf(format, a...)
	os.Exit(1)
}

func Panic(a ...interface{}) {
	Logger.Panic(a...)
}

func Panicf(format string, a ...interface{}) {
	Logger.Panicf(format, a...)
}

func Debugx(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		Debugf(format, a...)
		return
	}
	format = "[TRACEID:%v] " + format
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	Logger.Debugf(format, a...)
}

func Infox(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		Infof(format, a...)
		return
	}
	format = "[TRACEID:%v] " + format
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	Logger.Infof(format, a...)
}

func Warnx(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		Warnf(format, a...)
		return
	}
	format = "[TRACEID:%v] " + format
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	Logger.Warnf(format, a...)
}

func Errorx(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		Errorf(format, a...)
		return
	}
	format = "[TRACEID:%v] " + format
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	Logger.Errorf(format, a...)
}
