package log

import (
	"testing"

	"go.uber.org/zap"
)

/*
zap 默认Logger API测试

zap.NewDevelopment() 创建一个开发环境的日志记录器，它会在控制台输出日志，并且包含时间戳、调用者信息等调试信息。
zap.NewProduction() 创建一个生产环境的日志记录器，它会在控制台输出日志，并且包含时间戳、调用者信息等调试信息。
zap.NewExample() 创建一个示例日志记录器，它会在控制台输出日志，并且包含时间戳、调用者信息等调试信息。

*/

// JSON输出到STDERR，Infovel及以上，包括时间戳和caller信息
func TestZapDevelopment(t *testing.T) {
	logger, _ := zap.NewDevelopment() // console 风格输出
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")   // warn 会打印堆栈信息
	logger.Error("error") // error 会打印堆栈信息
}

// 控制台输出到stderr，debuglevel及以上，包括时间戳和caller信息
func TestZapProduction(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Debug("debug") // debug 不会输出
	logger.Info("info")
	logger.Warn("warn")   // warn 不会打印堆栈
	logger.Error("error") // error 会打印堆栈信息
}

// JSON输出到Stdout，可用于示例和测试
// json 风格输出
// 没有时间
func TestZapExample(t *testing.T) {
	logger := zap.NewExample() // json
	logger.Info("info")
	logger.Debug("debug")
	logger.Warn("warn")
	logger.Error("error")
}
