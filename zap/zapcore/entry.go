// Copyright (c) 2016 Uber Technologies, Inc.
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

package zapcore // zapcore包：定义zap日志库的核心低层接口

import (
	"fmt"     // fmt包：格式化I/O
	"runtime" // runtime包：运行时信息
	"strings" // strings包：字符串操作
	"time"    // time包：时间处理

	"go.uber.org/multierr"                // multierr包：多重错误处理
	"go.uber.org/zap/internal/bufferpool" // 内部缓冲池
	"go.uber.org/zap/internal/exit"       // 内部退出处理
	"go.uber.org/zap/internal/pool"       // 内部对象池
)

// CheckedEntry对象池，用于复用CheckedEntry对象以提高性能
var _cePool = pool.New(func() *CheckedEntry { // CheckedEntry对象池
	// Pre-allocate some space for cores.
	// 为core预分配一些空间。
	return &CheckedEntry{
		cores: make([]Core, 4), // 预分配4个core的空间
	}
})

// getCheckedEntry 从对象池获取CheckedEntry实例
func getCheckedEntry() *CheckedEntry { // 从对象池获取CheckedEntry
	ce := _cePool.Get() // 从池中获取
	ce.reset()          // 重置状态
	return ce           // 返回实例
}

// putCheckedEntry 将CheckedEntry归还到对象池
func putCheckedEntry(ce *CheckedEntry) { // 将CheckedEntry归还到对象池
	if ce == nil { // 如果为nil
		return // 直接返回
	}
	_cePool.Put(ce) // 归还到池中
}

// NewEntryCaller makes an EntryCaller from the return signature of
// runtime.Caller.
// NewEntryCaller根据runtime.Caller的返回签名创建EntryCaller。
func NewEntryCaller(pc uintptr, file string, line int, ok bool) EntryCaller {
	if !ok { // 如果调用失败
		return EntryCaller{} // 返回空的EntryCaller
	}
	return EntryCaller{ // 返回填充的EntryCaller
		PC:      pc,   // 程序计数器
		File:    file, // 文件名
		Line:    line, // 行号
		Defined: true, // 标记为已定义
	}
}

// EntryCaller represents the caller of a logging function.
// EntryCaller表示日志函数的调用者。
// EnterCaller存储有关生成日志的代码位置的信息
type EntryCaller struct {
	Defined  bool    // 是否已定义
	PC       uintptr // 程序计数器
	File     string  // 文件路径
	Line     int     // 行号
	Function string  // 函数名
}

// String returns the full path and line number of the caller.
// String返回调用者的完整路径和行号。
func (ec EntryCaller) String() string {
	return ec.FullPath() // 调用FullPath方法
}

// FullPath returns a /full/path/to/package/file:line description of the
// caller.
// FullPath返回调用者的/full/path/to/package/file:line描述。
func (ec EntryCaller) FullPath() string {
	if !ec.Defined { // 如果未定义
		return "undefined" // 返回"undefined"
	}
	buf := bufferpool.Get()       // 从池中获取缓冲区
	buf.AppendString(ec.File)     // 追加文件路径
	buf.AppendByte(':')           // 追加冒号
	buf.AppendInt(int64(ec.Line)) // 追加行号
	caller := buf.String()        // 获取字符串
	buf.Free()                    // 释放缓冲区
	return caller                 // 返回调用者信息
}

// TrimmedPath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
// TrimmedPath返回调用者的package/file:line描述，只保留叶子目录名和文件名
func (ec EntryCaller) TrimmedPath() string {
	if !ec.Defined {
		return "undefined"
	}
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	//
	// Find the last separator.
	//
	// 注意：为了确保在Windows上也能正确修剪路径，我们需要使用'/'而不是os.PathSeparator，
	// 因为路径来自Go标准库，特别是runtime.Caller()（截至2017年3月）即使在Windows上也返回正斜杠。
	//
	// 查找最后一个分隔符
	idx := strings.LastIndexByte(ec.File, '/')
	if idx == -1 {
		return ec.FullPath()
	}
	// Find the penultimate separator.
	// 查找倒数第二个分隔符
	idx = strings.LastIndexByte(ec.File[:idx], '/')
	if idx == -1 {
		return ec.FullPath()
	}
	buf := bufferpool.Get()
	// Keep everything after the penultimate separator.
	// 保留倒数第二个分隔符之后的所有内容
	buf.AppendString(ec.File[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}

// An Entry represents a complete log message. The entry's structured context
// is already serialized, but the log level, time, message, and call site
// information are available for inspection and modification. Any fields left
// empty will be omitted when encoding.
//
// Entries are pooled, so any functions that accept them MUST be careful not to
// retain references to them.
// Entry表示一个完整的日志消息。条目的结构化上下文已经序列化，
// 但日志级别、时间、消息和调用位置信息可用于检查和修改。
// 编码时会省略任何留空的字段。
//
// Entry是池化的，所以任何接受它们的函数必须小心不要保留对它们的引用。
// ENtry 包含用于日志消息的元数据（级别，时间，消息等）
type Entry struct {
	Level      Level       // 日志级别
	Time       time.Time   // 时间戳
	LoggerName string      // 日志器名称
	Message    string      // 日志消息
	Caller     EntryCaller // 调用者信息
	Stack      string      // 堆栈信息
}

// CheckWriteHook is a custom action that may be executed after an entry is
// written.
//
// Register one on a CheckedEntry with the After method.
//
//	if ce := logger.Check(...); ce != nil {
//	  ce = ce.After(hook)
//	  ce.Write(...)
//	}
//
// You can configure the hook for Fatal log statements at the logger level with
// the zap.WithFatalHook option.
// CheckWriteHook是在条目写入后可能执行的自定义操作。
//
// 使用After方法在CheckedEntry上注册一个钩子。
//
// 您可以使用zap.WithFatalHook选项在日志器级别为Fatal日志语句配置钩子。
type CheckWriteHook interface {
	// OnWrite is invoked with the CheckedEntry that was written and a list
	// of fields added with that entry.
	//
	// The list of fields DOES NOT include fields that were already added
	// to the logger with the With method.
	// OnWrite在写入CheckedEntry时被调用，并传入与该条目一起添加的字段列表。
	//
	// 字段列表不包括已经使用With方法添加到日志器的字段。
	OnWrite(*CheckedEntry, []Field)
}

// CheckWriteAction indicates what action to take after a log entry is
// processed. Actions are ordered in increasing severity.
// CheckWriteAction指示在处理日志条目后要采取的操作。操作按严重性递增排序。
type CheckWriteAction uint8

const (
	// WriteThenNoop indicates that nothing special needs to be done. It's the
	// default behavior.
	// WriteThenNoop表示不需要做任何特殊操作。这是默认行为。
	WriteThenNoop CheckWriteAction = iota
	// WriteThenGoexit runs runtime.Goexit after Write.
	// WriteThenGoexit在Write后运行runtime.Goexit。
	WriteThenGoexit
	// WriteThenPanic causes a panic after Write.
	// WriteThenPanic在Write后触发panic。
	WriteThenPanic
	// WriteThenFatal causes an os.Exit(1) after Write.
	// WriteThenFatal在Write后调用os.Exit(1)。
	WriteThenFatal
)

// OnWrite implements the OnWrite method to keep CheckWriteAction compatible
// with the new CheckWriteHook interface which deprecates CheckWriteAction.
// OnWrite实现OnWrite方法以保持CheckWriteAction与新的CheckWriteHook接口兼容，
// 该接口已弃用CheckWriteAction。
func (a CheckWriteAction) OnWrite(ce *CheckedEntry, _ []Field) {
	switch a {
	case WriteThenGoexit:
		runtime.Goexit() // 退出当前goroutine
	case WriteThenPanic:
		panic(ce.Message) // 触发panic
	case WriteThenFatal:
		exit.With(1) // 退出程序
	}
}

// 确保CheckWriteAction实现了CheckWriteHook接口
var _ CheckWriteHook = CheckWriteAction(0)

// CheckedEntry is an Entry together with a collection of Cores that have
// already agreed to log it.
//
// CheckedEntry references should be created by calling AddCore or After on a
// nil *CheckedEntry. References are returned to a pool after Write, and MUST
// NOT be retained after calling their Write method.
// CheckedEntry是一个Entry和已经同意记录它的Core集合的组合。
//
// CheckedEntry引用应该通过在nil *CheckedEntry上调用AddCore或After来创建。
// 引用在Write后返回到池中，在调用其Write方法后绝不能保留。
// checkedentry是一种汇总资源，只有在需要编写日志时才会分配
type CheckedEntry struct {
	Entry                      // 嵌入Entry结构体
	ErrorOutput WriteSyncer    // 错误输出
	dirty       bool           // 尽力检测池误用的标志
	after       CheckWriteHook // 写入后的钩子
	cores       []Core         // Core集合
}

// reset 重置CheckedEntry到初始状态
func (ce *CheckedEntry) reset() {
	ce.Entry = Entry{}   // 重置Entry
	ce.ErrorOutput = nil // 清空错误输出
	ce.dirty = false     // 重置dirty标志
	ce.after = nil       // 清空钩子
	for i := range ce.cores {
		// don't keep references to cores
		// 不保留对cores的引用
		ce.cores[i] = nil
	}
	ce.cores = ce.cores[:0] // 清空cores切片
}

// Write writes the entry to the stored Cores, returns any errors, and returns
// the CheckedEntry reference to a pool for immediate re-use. Finally, it
// executes any required CheckWriteAction.
// Write将条目写入存储的Cores，返回任何错误，并将CheckedEntry引用返回到池中以便立即重用。
// 最后，它执行任何必需的CheckWriteAction。
func (ce *CheckedEntry) Write(fields ...Field) {
	if ce == nil {
		return
	}

	if ce.dirty {
		if ce.ErrorOutput != nil {
			// Make a best effort to detect unsafe re-use of this CheckedEntry.
			// If the entry is dirty, log an internal error; because the
			// CheckedEntry is being used after it was returned to the pool,
			// the message may be an amalgamation from multiple call sites.
			// 尽力检测这个CheckedEntry的不安全重用。
			// 如果条目是脏的，记录内部错误；因为CheckedEntry在返回到池后被使用，
			// 消息可能是来自多个调用点的混合。
			_, _ = fmt.Fprintf(
				ce.ErrorOutput,
				"%v Unsafe CheckedEntry re-use near Entry %+v.\n",
				ce.Time,
				ce.Entry,
			)
			_ = ce.ErrorOutput.Sync() // ignore error
		}
		return
	}
	ce.dirty = true

	var err error
	for i := range ce.cores {
		err = multierr.Append(err, ce.cores[i].Write(ce.Entry, fields))
	}
	if err != nil && ce.ErrorOutput != nil {
		_, _ = fmt.Fprintf(
			ce.ErrorOutput,
			"%v write error: %v\n",
			ce.Time,
			err,
		)
		_ = ce.ErrorOutput.Sync() // ignore error
	}

	hook := ce.after
	if hook != nil {
		hook.OnWrite(ce, fields)
	}
	putCheckedEntry(ce)
}

// AddCore adds a Core that has agreed to log this CheckedEntry. It's intended to be
// used by Core.Check implementations, and is safe to call on nil CheckedEntry
// references.
// AddCore添加一个已经同意记录这个CheckedEntry的Core。它旨在被Core.Check实现使用，
// 在nil CheckedEntry引用上调用是安全的。
func (ce *CheckedEntry) AddCore(ent Entry, core Core) *CheckedEntry {
	if ce == nil {
		ce = getCheckedEntry()
		ce.Entry = ent
	}
	ce.cores = append(ce.cores, core)
	return ce
}

// Should sets this CheckedEntry's CheckWriteAction, which controls whether a
// Core will panic or fatal after writing this log entry. Like AddCore, it's
// safe to call on nil CheckedEntry references.
//
// Deprecated: Use [CheckedEntry.After] instead.
// Should设置这个CheckedEntry的CheckWriteAction，它控制Core在写入此日志条目后是否会panic或fatal。
// 像AddCore一样，在nil CheckedEntry引用上调用是安全的。
//
// 已弃用：请使用[CheckedEntry.After]代替。
func (ce *CheckedEntry) Should(ent Entry, should CheckWriteAction) *CheckedEntry {
	return ce.After(ent, should)
}

// After sets this CheckEntry's CheckWriteHook, which will be called after this
// log entry has been written. It's safe to call this on nil CheckedEntry
// references.
// After设置这个CheckEntry的CheckWriteHook，它将在写入此日志条目后被调用。
// 在nil CheckedEntry引用上调用是安全的。
func (ce *CheckedEntry) After(ent Entry, hook CheckWriteHook) *CheckedEntry {
	if ce == nil {
		ce = getCheckedEntry()
		ce.Entry = ent
	}
	ce.after = hook
	return ce
}
