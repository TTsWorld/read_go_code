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
func (o *ObservedLogs) Filter(keep func(LoggedEntry) bool) *ObservedLogs {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var filtered []LoggedEntry
	for _, entry := range o.logs {
		if keep(entry) {
			filtered = append(filtered, entry)
		}
	}
	return &ObservedLogs{logs: filtered}
}

func (o *ObservedLogs) add(log LoggedEntry) {
	o.mu.Lock()
	o.logs = append(o.logs, log)
	o.mu.Unlock()
}

// New creates a new Core that buffers logs in memory (without any encoding).
// It's particularly useful in tests.
func New(enab zapcore.LevelEnabler) (zapcore.Core, *ObservedLogs) {
	ol := &ObservedLogs{}
	return &contextObserver{
		LevelEnabler: enab,
		logs:         ol,
	}, ol
}

type contextObserver struct {
	zapcore.LevelEnabler
	logs    *ObservedLogs
	context []zapcore.Field
}

var (
	_ zapcore.Core            = (*contextObserver)(nil)
	_ internal.LeveledEnabler = (*contextObserver)(nil)
)

func (co *contextObserver) Level() zapcore.Level {
	return zapcore.LevelOf(co.LevelEnabler)
}

func (co *contextObserver) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if co.Enabled(ent.Level) {
		return ce.AddCore(ent, co)
	}
	return ce
}

func (co *contextObserver) With(fields []zapcore.Field) zapcore.Core {
	return &contextObserver{
		LevelEnabler: co.LevelEnabler,
		logs:         co.logs,
		context:      append(co.context[:len(co.context):len(co.context)], fields...),
	}
}

func (co *contextObserver) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	all := make([]zapcore.Field, 0, len(fields)+len(co.context))
	all = append(all, co.context...)
	all = append(all, fields...)
	co.logs.add(LoggedEntry{ent, all})
	return nil
}

func (co *contextObserver) Sync() error {
	return nil
}
