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

// Package stacktrace provides support for gathering stack traces
// efficiently.
// Package stacktrace提供高效收集堆栈跟踪的支持。
package stacktrace

import (
	"runtime" // runtime包：运行时服务

	"go.uber.org/zap/buffer"              // buffer包：高性能缓冲区
	"go.uber.org/zap/internal/bufferpool" // 内部缓冲池
	"go.uber.org/zap/internal/pool"       // 内部对象池
)

var _stackPool = pool.New(func() *Stack { // 堆栈对象池
	return &Stack{
		storage: make([]uintptr, 64), // 预分配64个程序计数器的存储空间
	}
})

// Stack is a captured stack trace.
// Stack是捕获的堆栈跟踪。
type Stack struct {
	pcs    []uintptr       // 程序计数器；始终是storage的子切片
	frames *runtime.Frames // 堆栈帧

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	// pcs的大小根据需求而定：如果只请求第一帧，它将是1，
	// 否则它将反映调用堆栈的深度。
	//
	// storage将我们需要的切片(pcs)与我们池化的切片分离。
	// 我们将始终分配一个合理大的存储，但只使用所需的部分。
	storage []uintptr // 存储空间
}

// Depth specifies how deep of a stack trace should be captured.
// Depth指定应捕获多深的堆栈跟踪。
type Depth int

const (
	// First captures only the first frame.
	// First仅捕获第一帧。
	First Depth = iota

	// Full captures the entire call stack, allocating more
	// storage for it if needed.
	// Full捕获整个调用堆栈，如需要则为其分配更多存储。
	Full
)

// Capture captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// Capture.
//
// The caller must call Free on the returned stacktrace after using it.
// Capture捕获指定深度的堆栈跟踪，跳过提供的帧数。skip=0标识Capture的调用者。
//
// 调用者在使用后必须在返回的堆栈跟踪上调用Free。
func Capture(skip int, depth Depth) *Stack {
	stack := _stackPool.Get() // 从对象池获取堆栈

	switch depth { // 根据深度选择策略
	case First: // 仅第一帧
		stack.pcs = stack.storage[:1] // 使用存储的第一个元素
	case Full: // 完整堆栈
		stack.pcs = stack.storage // 使用全部存储
	}

	// Unlike other "skip"-based APIs, skip=0 identifies runtime.Callers
	// itself. +2 to skip captureStacktrace and runtime.Callers.
	// 与其他基于"skip"的API不同，skip=0标识runtime.Callers本身。
	// +2跳过captureStacktrace和runtime.Callers。
	numFrames := runtime.Callers( // 获取调用者信息
		skip+2,    // 跳过的帧数
		stack.pcs, // 存储程序计数器的切片
	)

	// runtime.Callers truncates the recorded stacktrace if there is no
	// room in the provided slice. For the full stack trace, keep expanding
	// storage until there are fewer frames than there is room.
	// runtime.Callers如果提供的切片中没有空间，会截断记录的堆栈跟踪。
	// 对于完整的堆栈跟踪，继续扩展存储直到帧数少于空间数。
	if depth == Full { // 如果是完整深度
		pcs := stack.pcs            // 获取当前程序计数器切片
		for numFrames == len(pcs) { // 如果帧数等于切片长度（需要扩展）
			pcs = make([]uintptr, len(pcs)*2)        // 创建两倍大小的新切片
			numFrames = runtime.Callers(skip+2, pcs) // 重新获取调用者信息
		}

		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		// 丢弃旧存储而不是将其返回到池中。
		// 如果堆栈跟踪一直很深，这将随时间调整池大小。
		stack.storage = pcs         // 更新存储
		stack.pcs = pcs[:numFrames] // 设置有效的程序计数器切片
	} else {
		stack.pcs = stack.pcs[:numFrames] // 设置有效的程序计数器切片
	}

	stack.frames = runtime.CallersFrames(stack.pcs) // 创建帧迭代器
	return stack                                    // 返回堆栈
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
// Free释放与此堆栈跟踪相关的资源并将其返回到池中。
func (st *Stack) Free() {
	st.frames = nil    // 清空帧迭代器
	st.pcs = nil       // 清空程序计数器切片
	_stackPool.Put(st) // 返回到对象池
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
// Count报告此堆栈跟踪中的总帧数。调用Next时Count不会改变。
func (st *Stack) Count() int {
	return len(st.pcs) // 返回程序计数器切片的长度
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
// Next返回堆栈跟踪中的下一帧，以及指示之后是否还有更多帧的布尔值。
func (st *Stack) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next() // 调用帧迭代器的Next方法
}

// Take returns a string representation of the current stacktrace.
//
// skip is the number of frames to skip before recording the stack trace.
// skip=0 identifies the caller of Take.
// Take返回当前堆栈跟踪的字符串表示。
//
// skip是在记录堆栈跟踪之前要跳过的帧数。skip=0标识Take的调用者。
func Take(skip int) string {
	stack := Capture(skip+1, Full) // 捕获完整堆栈跟踪
	defer stack.Free()             // 延迟释放资源

	buffer := bufferpool.Get() // 从缓冲池获取缓冲区
	defer buffer.Free()        // 延迟释放缓冲区

	stackfmt := NewFormatter(buffer) // 创建格式化器
	stackfmt.FormatStack(stack)      // 格式化堆栈
	return buffer.String()           // 返回字符串
}

// Formatter formats a stack trace into a readable string representation.
// Formatter将堆栈跟踪格式化为可读的字符串表示。
type Formatter struct {
	b        *buffer.Buffer // 缓冲区
	nonEmpty bool           // 是否已经写入至少一帧
}

// NewFormatter builds a new Formatter.
// NewFormatter构建新的格式化器。
func NewFormatter(b *buffer.Buffer) Formatter {
	return Formatter{b: b} // 返回新的格式化器
}

// FormatStack formats all remaining frames in the provided stacktrace -- minus
// the final runtime.main/runtime.goexit frame.
// FormatStack格式化提供的堆栈跟踪中的所有剩余帧，除了最后的runtime.main/runtime.goexit帧。
func (sf *Formatter) FormatStack(stack *Stack) {
	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	// 注意：在最后一次迭代中，frames.Next()返回false，但有一个有效帧，
	// 但我们忽略这个帧。最后一帧是运行时帧，会增加噪声，
	// 因为它只是runtime.main或runtime.goexit。
	for frame, more := stack.Next(); more; frame, more = stack.Next() { // 迭代所有帧
		sf.FormatFrame(frame) // 格式化帧
	}
}

// FormatFrame formats the given frame.
// FormatFrame格式化给定的帧。
func (sf *Formatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty { // 如果不是第一帧
		sf.b.AppendByte('\n') // 添加换行符
	}
	sf.nonEmpty = true                // 标记为非空
	sf.b.AppendString(frame.Function) // 添加函数名
	sf.b.AppendByte('\n')             // 添加换行符
	sf.b.AppendByte('\t')             // 添加制表符
	sf.b.AppendString(frame.File)     // 添加文件名
	sf.b.AppendByte(':')              // 添加冒号
	sf.b.AppendInt(int64(frame.Line)) // 添加行号
}
