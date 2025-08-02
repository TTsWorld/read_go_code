// Copyright (c) 2016-2022 Uber Technologies, Inc.
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

// Package observer provides a zapcore.Core that keeps an in-memory,
// encoding-agnostic representation of log entries. It's useful for
// applications that want to unit test their log output without tying their
// tests to a particular output encoding.
// Package observer提供了一个zapcore.Core，它保持日志条目的内存中、编码无关的表示。
// 对于想要对其日志输出进行单元测试而不将测试绑定到特定输出编码的应用程序很有用。
package observer // import "go.uber.org/zap/zaptest/observer"

import (
	"strings" // strings包：字符串操作
	"sync"    // sync包：同步原语
	"time"    // time包：时间功能

	"go.uber.org/zap/internal" // zap内部包
	"go.uber.org/zap/zapcore"  // zapcore包：核心接口和实现
)

// ObservedLogs is a concurrency-safe, ordered collection of observed logs.
// ObservedLogs是一个并发安全的、有序的观察日志集合。
type ObservedLogs struct {
	mu   sync.RWMutex   // 读写互斥锁，保护logs字段
	logs []LoggedEntry  // 日志条目列表
}

// Len returns the number of items in the collection.
// Len返回集合中项目的数量。
func (o *ObservedLogs) Len() int {
	o.mu.RLock()        // 获取读锁
	n := len(o.logs)    // 获取日志数量
	o.mu.RUnlock()      // 释放读锁
	return n            // 返回数量
}

// All returns a copy of all the observed logs.
// All返回所有观察日志的副本。
func (o *ObservedLogs) All() []LoggedEntry {
	o.mu.RLock()                             // 获取读锁
	ret := make([]LoggedEntry, len(o.logs)) // 创建副本切片
	copy(ret, o.logs)                       // 复制日志数据
	o.mu.RUnlock()                          // 释放读锁
	return ret                              // 返回副本
}

// TakeAll returns a copy of all the observed logs, and truncates the observed
// slice.
// TakeAll返回所有观察日志的副本，并截断观察切片。
func (o *ObservedLogs) TakeAll() []LoggedEntry {
	o.mu.Lock()     // 获取写锁
	ret := o.logs   // 获取现有日志
	o.logs = nil    // 清空日志列表
	o.mu.Unlock()   // 释放写锁
	return ret      // 返回旧日志
}

// AllUntimed returns a copy of all the observed logs, but overwrites the
// observed timestamps with time.Time's zero value. This is useful when making
// assertions in tests.
// AllUntimed返回所有观察日志的副本，但用time.Time的零值覆盖观察的时间戳。
// 这在测试中进行断言时很有用。
func (o *ObservedLogs) AllUntimed() []LoggedEntry {
	ret := o.All()                      // 获取所有日志的副本
	for i := range ret {                // 遍历所有日志
		ret[i].Time = time.Time{} // 将时间设置为零值
	}
	return ret // 返回修改后的副本
}

// FilterLevelExact filters entries to those logged at exactly the given level.
// FilterLevelExact过滤条目到在确切给定级别记录的那些。
func (o *ObservedLogs) FilterLevelExact(level zapcore.Level) *ObservedLogs {
	return o.Filter(func(e LoggedEntry) bool { // 使用通用过滤方法
		return e.Level == level // 检查级别是否完全匹配
	})
}

// FilterMessage filters entries to those that have the specified message.
// FilterMessage过滤条目到具有指定消息的那些。
func (o *ObservedLogs) FilterMessage(msg string) *ObservedLogs {
	return o.Filter(func(e LoggedEntry) bool { // 使用通用过滤方法
		return e.Message == msg // 检查消息是否匹配
	})
}

// FilterLoggerName filters entries to those logged through logger with the specified logger name.
// FilterLoggerName过滤条目到通过具有指定日志器名称的日志器记录的那些。
func (o *ObservedLogs) FilterLoggerName(name string) *ObservedLogs {
	return o.Filter(func(e LoggedEntry) bool { // 使用通用过滤方法
		return e.LoggerName == name // 检查日志器名称是否匹配
	})
}

// FilterMessageSnippet filters entries to those that have a message containing the specified snippet.
// FilterMessageSnippet过滤条目到具有包含指定片段的消息的那些。
func (o *ObservedLogs) FilterMessageSnippet(snippet string) *ObservedLogs {
	return o.Filter(func(e LoggedEntry) bool { // 使用通用过滤方法
		return strings.Contains(e.Message, snippet) // 检查消息是否包含片段
	})
}

// FilterField filters entries to those that have the specified field.
// FilterField过滤条目到具有指定字段的那些。
func (o *ObservedLogs) FilterField(field zapcore.Field) *ObservedLogs {
	return o.Filter(func(e LoggedEntry) bool { // 使用通用过滤方法
		for _, ctxField := range e.Context { // 遍历上下文字段
			if ctxField.Equals(field) { // 检查字段是否相等
				return true // 找到匹配的字段
			}
		}
		return false // 未找到匹配的字段
	})
}

// FilterFieldKey filters entries to those that have the specified key.
// FilterFieldKey过滤条目到具有指定键的那些。
func (o *ObservedLogs) FilterFieldKey(key string) *ObservedLogs {
	return o.Filter(func(e LoggedEntry) bool { // 使用通用过滤方法
		for _, ctxField := range e.Context { // 遍历上下文字段
			if ctxField.Key == key { // 检查字段键是否匹配
				return true // 找到匹配的键
			}
		}
		return false // 未找到匹配的键
	})
}

// Filter returns a copy of this ObservedLogs containing only those entries
// for which the provided function returns true.
// Filter返回此ObservedLogs的副本，只包含提供的函数返回true的那些条目。
func (o *ObservedLogs) Filter(keep func(LoggedEntry) bool) *ObservedLogs {
	o.mu.RLock()         // 获取读锁
	defer o.mu.RUnlock() // 函数返回时释放读锁

	var filtered []LoggedEntry       // 过滤后的条目列表
	for _, entry := range o.logs {   // 遍历所有日志条目
		if keep(entry) {             // 如果保留函数返回true
			filtered = append(filtered, entry) // 添加到过滤列表
		}
	}
	return &ObservedLogs{logs: filtered} // 返回新的ObservedLogs实例
}

// add appends a log entry to the observed logs.
// add将日志条目附加到观察日志中。
func (o *ObservedLogs) add(log LoggedEntry) {
	o.mu.Lock()                        // 获取写锁
	o.logs = append(o.logs, log)       // 添加日志条目
	o.mu.Unlock()                      // 释放写锁
}

// New creates a new Core that buffers logs in memory (without any encoding).
// It's particularly useful in tests.
// New创建一个在内存中缓冲日志的新Core（不进行任何编码）。
// 它在测试中特别有用。
func New(enab zapcore.LevelEnabler) (zapcore.Core, *ObservedLogs) {
	ol := &ObservedLogs{}                    // 创建观察日志实例
	return &contextObserver{                // 返图上下文观察者和观察日志
		LevelEnabler: enab, // 级别启用器
		logs:         ol,   // 观察日志实例
	}, ol
}

// contextObserver is a Core implementation that stores logs in memory for testing.
// contextObserver是一个Core实现，将日志存储在内存中以供测试使用。
type contextObserver struct {
	zapcore.LevelEnabler     // 嵌入级别启用器
	logs    *ObservedLogs    // 观察日志实例
	context []zapcore.Field // 上下文字段列表
}

var (
	_ zapcore.Core            = (*contextObserver)(nil) // 编译时检查Core接口实现
	_ internal.LeveledEnabler = (*contextObserver)(nil) // 编译时检查LeveledEnabler接口实现
)

// Level returns the minimum enabled log level.
// Level返回最低启用的日志级别。
func (co *contextObserver) Level() zapcore.Level {
	return zapcore.LevelOf(co.LevelEnabler) // 获取级别启用器的级别
}

// Check determines whether the supplied Entry should be logged.
// Check确定是否应该记录提供的Entry。
func (co *contextObserver) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if co.Enabled(ent.Level) { // 如果级别启用
		return ce.AddCore(ent, co) // 添加到检查条目
	}
	return ce // 返回原检查条目
}

// With adds fields to the logger's context.
// With将字段添加到日志器的上下文。
func (co *contextObserver) With(fields []zapcore.Field) zapcore.Core {
	return &contextObserver{                                                      // 返回新的上下文观察者
		LevelEnabler: co.LevelEnabler,                                           // 复用级别启用器
		logs:         co.logs,                                                   // 复用观察日志实例
		context:      append(co.context[:len(co.context):len(co.context)], fields...), // 添加新字段到上下文
	}
}

// Write writes the entry to the observed logs.
// Write将条目写入观察日志。
func (co *contextObserver) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	all := make([]zapcore.Field, 0, len(fields)+len(co.context)) // 分配合并字段的空间
	all = append(all, co.context...)                            // 添加上下文字段
	all = append(all, fields...)                                // 添加当前字段
	co.logs.add(LoggedEntry{ent, all})                          // 将条目添加到观察日志
	return nil                                                  // 返回nil表示成功
}

// Sync is a no-op for the in-memory observer.
// Sync对于内存观察者来说是无操作。
func (co *contextObserver) Sync() error {
	return nil // 返回nil表示成功，内存中无需同步
}
