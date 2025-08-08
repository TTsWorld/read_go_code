// 文件功能: 提供 HTTP 日志中间件，包括简要日志(LogHandler)与详细日志(DetailedLogHandler)，记录请求/响应、耗时、慢调用、错误等信息。
// 关键技术点:
// - 中间件模式: 将日志逻辑以 `http.Handler` 链式插入请求处理流程。
// - 响应包装: 通过 `response.WithCodeResponseWriter` 捕获状态码与响应体。
// - 请求Body复制: 使用 `iox.LimitDupReadCloser` 限制复制大小，避免大包对内存/日志的影响。
// - 慢调用统计: 使用原子可变 `syncx.ForAtomicDuration` 设置阈值，动态调整慢调用标准。
// - 详细日志: 通过自定义 `detailLoggedResponseWriter` 缓存响应体以便记录。
// - 彩色输出: 根据 HTTP 方法与状态码着色便于扫描。
// - 工具集成: 使用 `httputil.DumpRequest`、`timex.ReprOfDuration` 等辅助函数。
// 适用场景: 生产环境服务访问日志、问题定位、行为审计。
// 包声明: 当前文件属于 `handler` 包，聚合 HTTP 相关中间件。
package handler

// 导入依赖包列表
import (
	// 缓冲、字节缓冲、错误、格式化、IO 等基础能力
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	// 网络与 HTTP 协议相关类型
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	// go-zero 核心能力: 颜色、IO 扩展、日志、原子、时间、工具
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/iox"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/core/utils"

	// rest 子系统: request 上下文日志、内部类型、响应包装
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/internal"
	"github.com/zeromicro/go-zero/rest/internal/response"
)

const (
	// 限制简要日志中请求体复制的最大字节数，避免过大 payload 影响性能
	limitBodyBytes = 1024
	// 限制详细日志中请求体复制的最大字节数
	limitDetailedBodyBytes = 4096
	// 默认慢调用阈值为 500ms，可动态修改
	defaultSlowThreshold = time.Millisecond * 500
)

// 原子可变的慢调用阈值，支持运行时调整
var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

// LogHandler returns a middleware that logs http request and response.
func LogHandler(next http.Handler) http.Handler {
	// 返回标准 `http.HandlerFunc` 作为中间件
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 计时器: 统计请求处理耗时
		timer := utils.NewElapsedTimer()
		// 收集器: 收集在请求上下文中追加的日志片段
		logs := new(internal.LogCollector)
		// 包装 ResponseWriter 以捕获状态码
		lrw := response.NewWithCodeResponseWriter(w)

		// 复制并限制请求体用于日志打印，避免读取一次后 body 丢失
		var dup io.ReadCloser
		r.Body, dup = iox.LimitDupReadCloser(r.Body, limitBodyBytes)
		// 将 LogCollector 注入到请求上下文，供后续 Handler/中间件添加日志
		next.ServeHTTP(lrw, r.WithContext(internal.WithLogCollector(r.Context(), logs)))
		// 复原原始 Body，便于后续可能的再次读取
		r.Body = dup
		// 打印简要日志
		logBrief(r, lrw.Code, timer, logs)
	})
}

// 详细日志记录时使用的 ResponseWriter 包装器
type detailLoggedResponseWriter struct {
	writer *response.WithCodeResponseWriter
	buf    *bytes.Buffer
}

// 创建 `detailLoggedResponseWriter`
func newDetailLoggedResponseWriter(writer *response.WithCodeResponseWriter,
	buf *bytes.Buffer) *detailLoggedResponseWriter {
	return &detailLoggedResponseWriter{
		writer: writer,
		buf:    buf,
	}
}

// 确保底层 ResponseWriter 的 Flush 被调用（若实现）
func (w *detailLoggedResponseWriter) Flush() {
	w.writer.Flush()
}

// 透传 Header
func (w *detailLoggedResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *detailLoggedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.writer.Writer.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

// 拦截 Write，既写入响应，也写入缓存以便日志
func (w *detailLoggedResponseWriter) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

// 透传 WriteHeader
func (w *detailLoggedResponseWriter) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}

// DetailedLogHandler returns a middleware that logs http request and response in details.
func DetailedLogHandler(next http.Handler) http.Handler {
	// 返回标准 `http.HandlerFunc` 作为中间件
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 计时器: 统计请求处理耗时
		timer := utils.NewElapsedTimer()
		// 缓存响应体以便记录
		var buf bytes.Buffer
		rw := response.NewWithCodeResponseWriter(w)
		lrw := newDetailLoggedResponseWriter(rw, &buf)

		var dup io.ReadCloser
		// https://github.com/zeromicro/go-zero/issues/3564
		// 复制并限制请求体用于详细日志
		r.Body, dup = iox.LimitDupReadCloser(r.Body, limitDetailedBodyBytes)
		// 注入日志收集器
		logs := new(internal.LogCollector)
		next.ServeHTTP(lrw, r.WithContext(internal.WithLogCollector(r.Context(), logs)))
		// 复原 Body
		r.Body = dup
		// 打印详细日志
		logDetails(r, lrw, timer, logs)
	})
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	// 运行时动态调整慢调用阈值
	slowThreshold.Set(threshold)
}

// 将请求转储为字符串，包含起始行、头、可选的 body
func dumpRequest(r *http.Request) string {
	reqContent, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err.Error()
	}

	return string(reqContent)
}

// 判断响应是否为非 5xx
func isOkResponse(code int) bool {
	// not server error
	return code < http.StatusInternalServerError
}

// 打印简要日志：方法、路径、远端、UA、耗时、状态码、上下文日志等
func logBrief(r *http.Request, code int, timer *utils.ElapsedTimer, logs *internal.LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	logger := logx.WithContext(r.Context()).WithDuration(duration)
	buf.WriteString(fmt.Sprintf("[HTTP] %s - %s %s - %s - %s",
		wrapStatusCode(code), wrapMethod(r.Method), r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent()))
	if duration > slowThreshold.Load() {
		logger.Slowf("[HTTP] %s - %s %s - %s - %s - slowcall(%s)",
			wrapStatusCode(code), wrapMethod(r.Method), r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent(),
			timex.ReprOfDuration(duration))
	}

	ok := isOkResponse(code)
	if !ok {
		buf.WriteString(fmt.Sprintf("\n%s", dumpRequest(r)))
	}

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("\n%s", body))
	}

	if ok {
		logger.Info(buf.String())
	} else {
		logger.Error(buf.String())
	}
}

// 打印详细日志：包含请求/响应完整内容与上下文日志
func logDetails(r *http.Request, response *detailLoggedResponseWriter, timer *utils.ElapsedTimer,
	logs *internal.LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	code := response.writer.Code
	logger := logx.WithContext(r.Context())
	buf.WriteString(fmt.Sprintf("[HTTP] %s - %d - %s - %s\n=> %s\n",
		r.Method, code, r.RemoteAddr, timex.ReprOfDuration(duration), dumpRequest(r)))
	if duration > slowThreshold.Load() {
		logger.Slowf("[HTTP] %s - %d - %s - slowcall(%s)\n=> %s\n", r.Method, code, r.RemoteAddr,
			timex.ReprOfDuration(duration), dumpRequest(r))
	}

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("%s\n", body))
	}

	respBuf := response.buf.Bytes()
	if len(respBuf) > 0 {
		buf.WriteString(fmt.Sprintf("<= %s", respBuf))
	}

	if isOkResponse(code) {
		logger.Info(buf.String())
	} else {
		logger.Error(buf.String())
	}
}

// 根据方法选择彩色输出以便于识别
func wrapMethod(method string) string {
	var colour color.Color
	switch method {
	case http.MethodGet:
		colour = color.BgBlue
	case http.MethodPost:
		colour = color.BgCyan
	case http.MethodPut:
		colour = color.BgYellow
	case http.MethodDelete:
		colour = color.BgRed
	case http.MethodPatch:
		colour = color.BgGreen
	case http.MethodHead:
		colour = color.BgMagenta
	case http.MethodOptions:
		colour = color.BgWhite
	}

	if colour == color.NoColor {
		return method
	}

	return logx.WithColorPadding(method, colour)
}

// 根据状态码区间选择彩色输出
func wrapStatusCode(code int) string {
	var colour color.Color
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		colour = color.BgGreen
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		colour = color.BgBlue
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		colour = color.BgMagenta
	default:
		colour = color.BgYellow
	}

	return logx.WithColorPadding(strconv.Itoa(code), colour)
}
