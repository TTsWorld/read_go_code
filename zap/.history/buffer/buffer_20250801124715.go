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

// Package buffer provides a thin wrapper around a byte slice. Unlike the
// standard library's bytes.Buffer, it supports a portion of the strconv
// package's zero-allocation formatters.
// Package buffer提供字节切片的轻量级包装器。与标准库的bytes.Buffer不同，
// 它支持strconv包的部分零分配格式化器。
package buffer // import "go.uber.org/zap/buffer"

import (
	"strconv" // strconv包：字符串转换
	"time"    // time包：时间处理
)

const _size = 1024 // by default, create 1 KiB buffers
                  // 默认创建1KiB的缓冲区

// Buffer is a thin wrapper around a byte slice. It's intended to be pooled, so
// the only way to construct one is via a Pool.
// Buffer是字节切片的轻量级包装器。它设计用于对象池化，
// 因此构造它的唯一方式是通过Pool。
type Buffer struct {
	bs   []byte // 底层字节切片
	pool Pool   // 所属的对象池
}

// AppendByte writes a single byte to the Buffer.
// AppendByte向Buffer写入单个字节。
func (b *Buffer) AppendByte(v byte) {
	b.bs = append(b.bs, v) // 追加字节到切片
}

// AppendBytes writes the given slice of bytes to the Buffer.
// AppendBytes向Buffer写入给定的字节切片。
func (b *Buffer) AppendBytes(v []byte) {
	b.bs = append(b.bs, v...) // 追加字节切片到底层切片
}

// AppendString writes a string to the Buffer.
// AppendString向Buffer写入字符串。
func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...) // 追加字符串到切片（零拷贝转换）
}

// AppendInt appends an integer to the underlying buffer (assuming base 10).
// AppendInt向底层缓冲区追加整数（假设为十进制）。
func (b *Buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10) // 使用strconv的零分配追加
}

// AppendTime appends the time formatted using the specified layout.
// AppendTime使用指定的布局格式追加时间。
func (b *Buffer) AppendTime(t time.Time, layout string) {
	b.bs = t.AppendFormat(b.bs, layout) // 使用time的零分配格式化追加
}

// AppendUint appends an unsigned integer to the underlying buffer (assuming
// base 10).
// AppendUint向底层缓冲区追加无符号整数（假设为十进制）。
func (b *Buffer) AppendUint(i uint64) {
	b.bs = strconv.AppendUint(b.bs, i, 10) // 使用strconv的零分配追加
}

// AppendBool appends a bool to the underlying buffer.
// AppendBool向底层缓冲区追加布尔值。
func (b *Buffer) AppendBool(v bool) {
	b.bs = strconv.AppendBool(b.bs, v) // 使用strconv的零分配追加
}

// AppendFloat appends a float to the underlying buffer. It doesn't quote NaN
// or +/- Inf.
// AppendFloat向底层缓冲区追加浮点数。它不会对NaN或+/-Inf加引号。
func (b *Buffer) AppendFloat(f float64, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize) // 使用strconv的零分配追加
}

// Len returns the length of the underlying byte slice.
// Len返回底层字节切片的长度。
func (b *Buffer) Len() int {
	return len(b.bs) // 返回字节切片长度
}

// Cap returns the capacity of the underlying byte slice.
// Cap返回底层字节切片的容量。
func (b *Buffer) Cap() int {
	return cap(b.bs) // 返回字节切片容量
}

// Bytes returns a mutable reference to the underlying byte slice.
// Bytes返回底层字节切片的可变引用。
func (b *Buffer) Bytes() []byte {
	return b.bs // 返回字节切片（可修改）
}

// String returns a string copy of the underlying byte slice.
// String返回底层字节切片的字符串副本。
func (b *Buffer) String() string {
	return string(b.bs) // 转换为字符串（会产生拷贝）
}

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
// Reset重置底层字节切片。后续写入会重用切片的底层数组。
func (b *Buffer) Reset() {
	b.bs = b.bs[:0] // 将长度重置为0，但保留容量和底层数组
}

// Write implements io.Writer.
// Write实现io.Writer接口。
func (b *Buffer) Write(bs []byte) (int, error) {
	b.bs = append(b.bs, bs...) // 追加字节切片
	return len(bs), nil        // 返回写入字节数和nil错误
}

// WriteByte writes a single byte to the Buffer.
//
// Error returned is always nil, function signature is compatible
// with bytes.Buffer and bufio.Writer
// WriteByte向Buffer写入单个字节。
//
// 返回的错误总是nil，函数签名与bytes.Buffer和bufio.Writer兼容
func (b *Buffer) WriteByte(v byte) error {
	b.AppendByte(v) // 调用AppendByte
	return nil      // 总是返回nil错误
}

// WriteString writes a string to the Buffer.
//
// Error returned is always nil, function signature is compatible
// with bytes.Buffer and bufio.Writer
// WriteString向Buffer写入字符串。
//
// 返回的错误总是nil，函数签名与bytes.Buffer和bufio.Writer兼容
func (b *Buffer) WriteString(s string) (int, error) {
	b.AppendString(s) // 调用AppendString
	return len(s), nil // 返回字符串长度和nil错误
}

// TrimNewline trims any final "\n" byte from the end of the buffer.
// TrimNewline从缓冲区末尾修剪任何最后的"\n"字节。
func (b *Buffer) TrimNewline() {
	if i := len(b.bs) - 1; i >= 0 { // 获取最后一个字节的索引
		if b.bs[i] == '\n' {        // 如果最后一个字节是换行符
			b.bs = b.bs[:i]         // 移除最后一个字节
		}
	}
}

// Free returns the Buffer to its Pool.
//
// Callers must not retain references to the Buffer after calling Free.
// Free将Buffer返回到其Pool中。
//
// 调用者在调用Free后不得保留对Buffer的引用。
func (b *Buffer) Free() {
	b.pool.put(b) // 将buffer归还给对象池
}
