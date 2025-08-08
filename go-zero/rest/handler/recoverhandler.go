// 文件功能: 提供 HTTP 中间件 `RecoverHandler`，在处理请求时拦截可能发生的 panic，记录堆栈并统一返回 500。
// 关键技术点:
// - 中间件模式: 通过返回 `http.Handler`/`http.HandlerFunc` 将逻辑串联到处理链中。
// - panic/recover: 使用 `defer` + `recover()` 捕获运行时 panic，避免进程崩溃。
// - 调试堆栈: 借助 `runtime/debug.Stack()` 获取当前 goroutine 的调用栈以便问题排查。
// - 日志记录: 调用框架内部 `internal.Error` 记录错误与上下文信息。
// - HTTP 协议: 通过 `w.WriteHeader(http.StatusInternalServerError)` 返回标准 500 状态码。
// 适用场景: 统一兜底异常，保护线上稳定性，避免 panic 传播导致整个服务中断。
// 包声明: 当前文件属于 `handler` 包，聚合 HTTP 相关中间件。
package handler

// 导入依赖包列表
import (
	// `fmt` 用于格式化字符串，构造错误信息
	"fmt"
	// `net/http` 提供 HTTP 服务器与请求/响应相关类型
	"net/http"
	// `runtime/debug` 用于获取调试信息，如当前 goroutine 的堆栈
	"runtime/debug"

	// 框架内部模块，封装了请求上下文中的日志处理等功能
	"github.com/zeromicro/go-zero/rest/internal"
)

// RecoverHandler returns a middleware that recovers if panic happens.
func RecoverHandler(next http.Handler) http.Handler {
	// 将闭包转为 `http.HandlerFunc`，以便作为中间件插入处理链
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 延迟执行的收尾逻辑，确保后续代码发生 panic 也能被捕获
		defer func() {
			// `recover()` 捕获 panic 值，非 nil 表示发生了 panic
			if result := recover(); result != nil {
				// 记录错误：panic 内容与堆栈，便于排查根因
				internal.Error(r, fmt.Sprintf("%v\n%s", result, debug.Stack()))
				// 返回 500，避免将内部错误细节暴露给客户端
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		// 交给下一个处理器继续处理请求
		next.ServeHTTP(w, r)
	})
}
