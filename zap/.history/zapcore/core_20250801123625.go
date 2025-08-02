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

// Core is a minimal, fast logger interface. It's designed for library authors
// to wrap in a more user-friendly API.
// Core是一个最小化、高速的日志接口。它设计给库作者使用，
// 用于封装成更用户友好的API
type Core interface {
	LevelEnabler // 嵌入LevelEnabler接口，提供级别启用功能

	// With adds structured context to the Core.
	// With方法向Core添加结构化上下文字段
	With([]Field) Core
	// Check determines whether the supplied Entry should be logged (using the
	// embedded LevelEnabler and possibly some extra logic). If the entry
	// should be logged, the Core adds itself to the CheckedEntry and returns
	// the result.
	//
	// Callers must use Check before calling Write.
	// Check方法确定提供的Entry是否应该被记录（使用嵌入的LevelEnabler
	// 以及可能的额外逻辑）。如果条目应该被记录，Core将自己添加到
	// CheckedEntry中并返回结果。
	//
	// 调用者必须在调用Write之前使用Check
	Check(Entry, *CheckedEntry) *CheckedEntry
	// Write serializes the Entry and any Fields supplied at the log site and
	// writes them to their destination.
	//
	// If called, Write should always log the Entry and Fields; it should not
	// replicate the logic of Check.
	// Write方法序列化Entry和在日志记录点提供的任何Fields，
	// 并将它们写入目标位置。
	//
	// 如果被调用，Write应该总是记录Entry和Fields；
	// 它不应该重复Check的逻辑
	Write(Entry, []Field) error
	// Sync flushes buffered logs (if any).
	// Sync方法刷新缓冲的日志（如果有的话）
	Sync() error
}

type nopCore struct{} // nopCore结构体：无操作的Core实现

// NewNopCore returns a no-op Core.
// NewNopCore返回一个无操作的Core实现
func NewNopCore() Core                                        { return nopCore{} }
func (nopCore) Enabled(Level) bool                            { return false } // Enabled始终返回false，表示不启用任何级别
func (n nopCore) With([]Field) Core                           { return n }     // With返回自身，不添加任何字段
func (nopCore) Check(_ Entry, ce *CheckedEntry) *CheckedEntry { return ce }    // Check直接返回传入的CheckedEntry，不做任何处理
func (nopCore) Write(Entry, []Field) error                    { return nil }   // Write不执行任何写入操作，直接返回nil
func (nopCore) Sync() error                                   { return nil }   // Sync不执行任何同步操作，直接返回nil

// NewCore creates a Core that writes logs to a WriteSyncer.
// NewCore创建一个将日志写入WriteSyncer的Core
func NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core {
	return &ioCore{
		LevelEnabler: enab, // 级别启用器，控制哪些级别的日志被记录
		enc:          enc,  // 编码器，负责序列化日志条目
		out:          ws,   // 写入同步器，负责将编码后的日志写入目标
	}
}

type ioCore struct { // ioCore结构体：基于IO的Core实现
	LevelEnabler        // 嵌入级别启用器
	enc Encoder         // 编码器实例
	out WriteSyncer     // 输出写入器
}

var (
	_ Core           = (*ioCore)(nil) // 编译时检查：确保ioCore实现了Core接口
	_ leveledEnabler = (*ioCore)(nil) // 编译时检查：确保ioCore实现了leveledEnabler接口
)

func (c *ioCore) Level() Level { // Level返回当前Core的日志级别
	return LevelOf(c.LevelEnabler) // 从LevelEnabler中获取级别
}

func (c *ioCore) With(fields []Field) Core { // With方法创建带有额外字段的新Core
	clone := c.clone()              // 克隆当前Core
	addFields(clone.enc, fields)    // 向克隆的编码器添加字段
	return clone                    // 返回新的Core实例
}

func (c *ioCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry { // Check检查是否应该记录指定级别的日志
	if c.Enabled(ent.Level) {       // 如果当前级别启用了该日志级别
		return ce.AddCore(ent, c)   // 将当前Core添加到CheckedEntry中
	}
	return ce                       // 否则直接返回原CheckedEntry
}

func (c *ioCore) Write(ent Entry, fields []Field) error { // Write方法执行实际的日志写入操作
	buf, err := c.enc.EncodeEntry(ent, fields) // 使用编码器编码日志条目和字段
	if err != nil {                             // 如果编码失败
		return err                              // 返回错误
	}
	_, err = c.out.Write(buf.Bytes())          // 将编码后的字节写入输出器
	buf.Free()                                 // 释放缓冲区资源
	if err != nil {                            // 如果写入失败
		return err                             // 返回错误
	}
	if ent.Level > ErrorLevel {                // 如果是严重错误级别（Fatal/Panic）
		// Since we may be crashing the program, sync the output.
		// Ignore Sync errors, pending a clean solution to issue #370.
		// 由于程序可能即将崩溃，同步输出
		// 忽略Sync错误，等待issue #370的完整解决方案
		_ = c.Sync()                           // 强制同步输出（忽略错误）
	}
	return nil                                 // 成功返回nil
}

func (c *ioCore) Sync() error { // Sync方法同步输出器的缓冲数据
	return c.out.Sync()         // 调用输出器的Sync方法
}

func (c *ioCore) clone() *ioCore { // clone方法创建ioCore的副本
	return &ioCore{
		LevelEnabler: c.LevelEnabler, // 复用相同的级别启用器
		enc:          c.enc.Clone(),  // 克隆编码器（编码器可能有状态）
		out:          c.out,          // 复用相同的输出器（输出器通常是无状态的）
	}
}
