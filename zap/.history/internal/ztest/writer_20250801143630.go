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

package ztest // ztest包：测试工具

import (
	"bytes"   // bytes包：字节操作
	"errors"  // errors包：错误处理
	"io"      // io包：基本I/O原语
	"strings" // strings包：字符串操作
)

// A Syncer is a spy for the Sync portion of zapcore.WriteSyncer.
// Syncer是zapcore.WriteSyncer的Sync部分的间谍。
type Syncer struct {
	err    error // 返回的错误
	called bool  // 是否被调用过
}

// SetError sets the error that the Sync method will return.
// SetError设置Sync方法将返回的错误。
func (s *Syncer) SetError(err error) {
	s.err = err // 设置错误
}

// Sync records that it was called, then returns the user-supplied error (if
// any).
// Sync记录它被调用，然后返回用户提供的错误（如果有）。
func (s *Syncer) Sync() error {
	s.called = true // 标记为已调用
	return s.err    // 返回设置的错误
}

// Called reports whether the Sync method was called.
// Called报告Sync方法是否被调用。
func (s *Syncer) Called() bool {
	return s.called // 返回调用状态
}

// A Discarder sends all writes to io.Discard.
// Discarder将所有写入发送到io.Discard。
type Discarder struct{ Syncer } // 嵌入Syncer

// Write implements io.Writer.
// Write实现io.Writer接口。
func (d *Discarder) Write(b []byte) (int, error) {
	return io.Discard.Write(b) // 将数据写入到丢弃器
}

// FailWriter is a WriteSyncer that always returns an error on writes.
// FailWriter是一个在写入时总是返回错误的WriteSyncer。
type FailWriter struct{ Syncer } // 嵌入Syncer

// Write implements io.Writer.
// Write实现io.Writer接口。
func (w FailWriter) Write(b []byte) (int, error) {
	return len(b), errors.New("failed") // 返回长度但报告失败错误
}

// ShortWriter is a WriteSyncer whose write method never fails, but
// nevertheless fails to the last byte of the input.
// ShortWriter是一个WriteSyncer，其写入方法从不失败，但却无法写入输入的最后一个字节。
type ShortWriter struct{ Syncer } // 嵌入Syncer

// Write implements io.Writer.
// Write实现io.Writer接口。
func (w ShortWriter) Write(b []byte) (int, error) {
	return len(b) - 1, nil // 返回比实际少一个字节的长度
}

// Buffer is an implementation of zapcore.WriteSyncer that sends all writes to
// a bytes.Buffer. It has convenience methods to split the accumulated buffer
// on newlines.
// Buffer是zapcore.WriteSyncer的实现，将所有写入发送到bytes.Buffer。
// 它有便利方法来在换行符上分割累积的缓冲区。
type Buffer struct {
	bytes.Buffer // 嵌入字节缓冲区
	Syncer       // 嵌入同步器
}

// Lines returns the current buffer contents, split on newlines.
// Lines返回当前缓冲区内容，在换行符上分割。
func (b *Buffer) Lines() []string {
	output := strings.Split(b.String(), "\n") // 按换行符分割
	return output[:len(output)-1]            // 返回除最后一个空元素外的所有行
}

// Stripped returns the current buffer contents with the last trailing newline
// stripped.
// Stripped返回当前缓冲区内容，去除最后的尾随换行符。
func (b *Buffer) Stripped() string {
	return strings.TrimRight(b.String(), "\n") // 去除右侧换行符
}
