// Copyright (c) 2023 Uber Technologies, Inc.
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

//go:build go1.21

// 构建约束：需要Go 1.21或更高版本

package zapslog // zapslog包：slog与Zap的集成

import "log/slog" // slog包：结构化日志记录

// A HandlerOption configures a slog Handler.
// HandlerOption配置slog Handler。
type HandlerOption interface {
	apply(*Handler) // 应用配置到Handler
}

// handlerOptionFunc wraps a func so it satisfies the Option interface.
// handlerOptionFunc包装函数使其满足Option接口。
type handlerOptionFunc func(*Handler)

func (f handlerOptionFunc) apply(handler *Handler) { // 实现HandlerOption接口
	f(handler) // 调用包装的函数
}

// WithName configures the Logger to annotate each message with the logger name.
// WithName配置Logger为每个消息添加日志器名称注释。
func WithName(name string) HandlerOption {
	return handlerOptionFunc(func(h *Handler) { // 返回选项函数
		h.name = name // 设置日志器名称
	})
}

// WithCaller configures the Logger to include the filename and line number
// of the caller in log messages--if available.
// WithCaller配置Logger在日志消息中包含调用者的文件名和行号（如果可用）。
func WithCaller(enabled bool) HandlerOption {
	return handlerOptionFunc(func(handler *Handler) { // 返回选项函数
		handler.addCaller = enabled // 设置是否添加调用者信息
	})
}

// WithCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the [WithCaller] option).
//
// When building wrappers around the Logger,
// supplying this Option prevents Zap from always reporting
// the wrapper code as the caller.
// WithCallerSkip增加调用者注释跳过的调用者数量（由[WithCaller]选项启用）。
//
// 在构建Logger的包装器时，提供此选项可防止Zap始终将包装器代码报告为调用者。
func WithCallerSkip(skip int) HandlerOption {
	return handlerOptionFunc(func(log *Handler) { // 返回选项函数
		log.callerSkip += skip // 增加跳过的调用者数量
	})
}

// AddStacktraceAt configures the Logger to record a stack trace
// for all messages at or above a given level.
// AddStacktraceAt配置Logger为所有处于或高于给定级别的消息记录堆栈跟踪。
func AddStacktraceAt(lvl slog.Level) HandlerOption {
	return handlerOptionFunc(func(log *Handler) { // 返回选项函数
		log.addStackAt = lvl // 设置堆栈跟踪的级别
	})
}
