// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package io provides basic interfaces to I/O primitives.
// Its primary job is to wrap existing implementations of such primitives,
// such as those in package os, into shared public interfaces that
// abstract the functionality, plus some other related primitives.
// io 包提供了 I/O 原语的基本接口, 它的主要工作是包装这些原语的现有实现， 例如 os 包中的
// 那些，放到共享的公共接口中抽象功能，再加上一些其他相关的原语。
//
// Because these interfaces and primitives wrap lower-level operations with
// various implementations, unless otherwise informed clients should not
// assume they are safe for parallel execution.
// 因为这些接口和原语用各种实现包装的都是低层次的操作，除非经过确认，否则客户端不应该假定他们
// 是并发安全的
package io

import (
	"errors"
	"sync"
)

// Seek whence values.
// seek 位置值
const (
	SeekStart = 0 // seek relative to the origin of the file
	// 从文件起点
	SeekCurrent = 1 // seek relative to the current offset
	// 从当前 offset
	SeekEnd = 2 // seek relative to the end
	// 文件末尾
)

// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
// ErrShortWrite 含义是一个写操作接收到的字节数比申请的要小，但未能返回一个明确错误
var ErrShortWrite = errors.New("short write")

// errInvalidWrite means that a write returned an impossible count.
// errInvalidWrite 表示 写操作返回一个不可能的 计数
var errInvalidWrite = errors.New("invalid write result")

// ErrShortBuffer means that a read required a longer buffer than was provided.
// ErrShortBuffer 表示读操作需要比提供的 buffer 更大的缓冲区
var ErrShortBuffer = errors.New("short buffer")

// EOF is the error returned by Read when no more input is available.
// (Read must return EOF itself, not an error wrapping EOF,
// because callers will test for EOF using ==.)
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
// 当没有更多的输入可用时，读操作返回一个 EOF 错误。函数应该只返回 EOF 来表示输入的优雅结束。
// 如果 EOF 在一个结构化的数据流时意外出现，恰到的做法是提供ErrUnexpectedEOF 或其他更详细的信息
var EOF = errors.New("EOF")

// ErrUnexpectedEOF means that EOF was encountered in the
// middle of reading a fixed-size block or data structure.
// ErrUnexpectedEOF 的含义是在读取固定大小的块或数据结构的中间发生了 EOF 错误
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// ErrNoProgress is returned by some clients of a Reader when
// many calls to Read have failed to return any data or error,
// usually the sign of a broken Reader implementation.
// ErrNoProgress 在许多读操作无法返回任何数据或错误时由某些客户端返回，
// 该错误通常表示一个损坏的 Reader 实现
var ErrNoProgress = errors.New("multiple Read calls return no data or error")

// Reader is the interface that wraps the basic Read method.
// Reader 是一个封装了基本 Read 方法的接口
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
// Read 至多读取 len(p) 个自己到 p中。它返回读取的字节数和任何发生的 error。即使
// 在 Read 返回的 n < len(p)，它将在调用时使用所有p 作为暂存空间
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
// 当 Read 在成功读到 n > 0 bytes 后遇到一个错误或读到文件末尾，它将返回读到的字节
// 它将在一些调用时返回 (non-nil) error 或后续调用中返回一个 error ( 同时n=0)。
// 这个通用 case 的一个实例是 Reader 在在读取一个输入流末尾时返回一个非 0 的字节数，
// 在输入流的尾部将会返回 err == EOF 或 err == nil。 下一次 Read 将返回 0，EOF。
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
// 调用者在考虑err 错误之前应当总是先处理返回的 n>0 的字节。在读取一些字节的以后，
// 处理 I/O 错误 和
//
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
// Read的实现不鼓励返回一个 0 字节的计数和一个空的错误，除非 len(p) == 0。
// 调用者应当将返回值为 0 和 nil 的返回视为什么都没有发生的标志；特别是它不表示 EOF
//
// Implementations must not retain p.
// 实现一定不能保留 p
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
// Writer 是基本的 Write 方法的包装接口
//
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
// Write 从将 p 指向的 len(p) 个字节写入底层字节流。它返回从 p 写入的字节数和任何会导致
// 写入操作提前停止的 err。如果 n < len(p), Write 一定返回一个非 nil 的错误。
// Write 一定不要修改底层 slice 数据，甚至临时修改。
//
// Implementations must not retain p.
// 实现一定不能保留 p
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Closer is the interface that wraps the basic Close method.
// Closer 是 Close 基本方法的包装接口
//
// The behavior of Close after the first call is undefined.
// Specific implementations may document their own behavior.
// 在首次调用后再次调用 Close的行为是未定义的，具体实现应当在自己的文档中说明
//

type Closer interface {
	Close() error
}

// Seeker is the interface that wraps the basic Seek method.
// Seeker 是基础 Seek 方法的包装接口
//
// Seek sets the offset for the next Read or Write to offset,
// interpreted according to whence:
// SeekStart means relative to the start of the file,
// SeekCurrent means relative to the current offset, and
// SeekEnd means relative to the end.
// Seek returns the new offset relative to the start of the
// file or an error, if any.
// Seek 设置下一次读或写的偏移，具体偏移根据 whence 参数解释:
// SeekStart 表示相对于文件的开头
// SeekCurrent 表示相对于当前偏移量
// SeekEnd 表示相对于文件末尾
// Seek 返回相对于文件起始位置的新偏移量或一个错误，如果有的话
//
// Seeking to an offset before the start of the file is an error.
// Seeking to any positive offset may be allowed, but if the new offset exceeds
// the size of the underlying object the behavior of subsequent I/O operations
// is implementation-dependent.
// 将偏移量设置在文件开始之前是一个 error。设置到任意正 offset 是被允许的，如果新的偏移量超过
// 底层对象的长度，那么后续 I/O操作的的行为依赖于具体实现
//
type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

// ReadWriter is the interface that groups the basic Read and Write methods.
// ReadWriter 是组合基本的 Read 和 Write 方法的 interface
type ReadWriter interface {
	Reader
	Writer
}

// ReadCloser is the interface that groups the basic Read and Close methods.
type ReadCloser interface {
	Reader
	Closer
}

// WriteCloser is the interface that groups the basic Write and Close methods.
type WriteCloser interface {
	Writer
	Closer
}

// ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// ReadSeeker is the interface that groups the basic Read and Seek methods.
type ReadSeeker interface {
	Reader
	Seeker
}

// ReadSeekCloser is the interface that groups the basic Read, Seek and Close
// methods.
type ReadSeekCloser interface {
	Reader
	Seeker
	Closer
}

// WriteSeeker is the interface that groups the basic Write and Seek methods.
type WriteSeeker interface {
	Writer
	Seeker
}

// ReadWriteSeeker is the interface that groups the basic Read, Write and Seek methods.
type ReadWriteSeeker interface {
	Reader
	Writer
	Seeker
}

// ReaderFrom is the interface that wraps the ReadFrom method.
// ReaderFrom 是一个保障了基本 ReadFrom 方法的接口
//
// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
// ReadFrom 从 r 读取数据知道遇到 EOF 或 error
// 返回值 n 是读到的字节数，咋读的过程中遇到除了 EOF 外的任何错误都会返回
//
// The Copy function uses ReaderFrom if available.
// 如果可用的话，copy 函数也使用 ReaderFrom
type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

// WriterTo is the interface that wraps the WriteTo method.
// WriteTo 是包装基本 WriteTo 方法的 interface
//
// WriteTo writes data to w until there's no more data to write or
// when an error occurs. The return value n is the number of bytes
// written. Any error encountered during the write is also returned.
// WriteTo 写数据到 w 直到没有数据科协或发生错误，返回值 n 是写入的字节数。写的过程中
// 遇到任何错误都将被返回
//
// The Copy function uses WriterTo if available.
// 如果可用的情况下，Copy 函数将使用 WriteTo
type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

// ReaderAt is the interface that wraps the basic ReadAt method.
// ReaderAt 是基本 ReadAt 函数的包装接口（WriterAt对称操作将不再翻译）
//
// ReadAt reads len(p) bytes into p starting at offset off in the
// underlying input source. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered.
// ReadAt 从一个底层输入源偏移开始读 len(p) 字节数据到 p，它返回读取的字节数
// 和发生的错误
//
// When ReadAt returns n < len(p), it returns a non-nil error
// explaining why more bytes were not returned. In this respect,
// ReadAt is stricter than Read.
// 当 ReadAt 返回 n < len(p), 它返回一个非 nil 的错误以说明为什么没有更多的字节返回
// 在这方面，ReadAt 比 Read 更加严格
//
// Even if ReadAt returns n < len(p), it may use all of p as scratch
// space during the call. If some data is available but not len(p) bytes,
// ReadAt blocks until either all the data is available or an error occurs.
// In this respect ReadAt is different from Read.
// 即使 ReadAt 返回的 n < len(p), 它在调用时，也将使用所有 p 作为暂存空间。如果一些数据可用
// 但是不够 len(p) 字节，ReadAt 将阻塞知道所有的数据可用或发生一个错误。
// 在这方面，ReadAt 和 Read 不同。
//
// If the n = len(p) bytes returned by ReadAt are at the end of the
// input source, ReadAt may return either err == EOF or err == nil.
// 在输入源结束时 如果 ReadAt 返回的 n = len(p)，ReadAt 将返回 err == EOF
// 或 err == nil 中任意一个
//
// If ReadAt is reading from an input source with a seek offset,
// ReadAt should not affect nor be affected by the underlying
// seek offset.
// 如果 ReadAt 从一个携带偏移量的输入源读数据，ReadAt 不应该影响也不会受到
// 底层 seek 偏移量的影响
//
// Clients of ReadAt can execute parallel ReadAt calls on the
// same input source.
// ReadAt 客户端可以在相同的输入源调用上执行并发操作
//
// Implementations must not retain p.
// 实现不得保留 p
type ReaderAt interface {
	ReadAt(p []byte, off int64) (n int, err error)
}

// WriterAt is the interface that wraps the basic WriteAt method.
//
// WriteAt writes len(p) bytes from p to the underlying data stream
// at offset off. It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// WriteAt must return a non-nil error if it returns n < len(p).
//
// If WriteAt is writing to a destination with a seek offset,
// WriteAt should not affect nor be affected by the underlying
// seek offset.
//
// Clients of WriteAt can execute parallel WriteAt calls on the same
// destination if the ranges do not overlap.
//
// Implementations must not retain p.
type WriterAt interface {
	WriteAt(p []byte, off int64) (n int, err error)
}

// ByteReader is the interface that wraps the ReadByte method.
// ByteReader 是包装了 ReadByte 方法的 interface
//
//
// ReadByte reads and returns the next byte from the input or
// any error encountered. If ReadByte returns an error, no input
// byte was consumed, and the returned byte value is undefined.
// ReadByte 从输入读取和返回下一个字节或遇到的错误。如果 ReadByte 返回一个错误，
// 没有输入字节被消费，并且返回字节值是未定义的
//
// ReadByte provides an efficient interface for byte-at-time
// processing. A Reader that does not implement  ByteReader
// can be wrapped using bufio.NewReader to add this method.
// ReadByte 提供为实时字节处理提供一个高效的接口。一个未实现 ByteReader接口
// 的 Reader 可以使用 bufio.NewReader 包装和添加这个方法
type ByteReader interface {
	ReadByte() (byte, error)
}

// ByteScanner is the interface that adds the UnreadByte method to the
// basic ReadByte method.
// ByteScanner  是一个在基础的 ReadByte 方法上添加了 UnreadByte 方法的接口
//
// UnreadByte causes the next call to ReadByte to return the last byte read.
// If the last operation was not a successful call to ReadByte, UnreadByte may
// return an error, unread the last byte read (or the byte prior to the
// last-unread byte), or (in implementations that support the Seeker interface)
// seek to one byte before the current offset.
// UnreadByte 导致下一次 ReadByte调用返回最后一个读到的字节。如果最后一次操作不是一个成功的 ReadByte 调用
// UnreadByte 会返回一个 error.unread 最后一个读取到的字节或 seek 1 个字节到当前 offset 之前
//
type ByteScanner interface {
	ByteReader
	UnreadByte() error
}

// ByteWriter is the interface that wraps the WriteByte method.
// ByteWriter 是一个包装了 WriteByte 方法的 interface
type ByteWriter interface {
	WriteByte(c byte) error
}

// RuneReader is the interface that wraps the ReadRune method.
//
// ReadRune reads a single encoded Unicode character
// and returns the rune and its size in bytes. If no character is
// available, err will be set.
type RuneReader interface {
	ReadRune() (r rune, size int, err error)
}

// RuneScanner is the interface that adds the UnreadRune method to the
// basic ReadRune method.
//
// UnreadRune causes the next call to ReadRune to return the last rune read.
// If the last operation was not a successful call to ReadRune, UnreadRune may
// return an error, unread the last rune read (or the rune prior to the
// last-unread rune), or (in implementations that support the Seeker interface)
// seek to the start of the rune before the current offset.
type RuneScanner interface {
	RuneReader
	UnreadRune() error
}

// StringWriter is the interface that wraps the WriteString method.
// StringWriter 是一个包装了 WriteString 的接口
type StringWriter interface {
	WriteString(s string) (n int, err error)
}

// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
// If w implements StringWriter, its WriteString method is invoked directly.
// Otherwise, w.Write is called exactly once.
// WriteString 接收一个 bytes 数组，将 string 的内容写入到 w
// 如果 w 实现了 StringWriter, 会直接调用它的 WriteString 方法。否则会调用 w.Write。
func WriteString(w Writer, s string) (n int, err error) {
	if sw, ok := w.(StringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}

// ReadAtLeast reads from r into buf until it has read at least min bytes.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading fewer than min bytes,
// ReadAtLeast returns ErrUnexpectedEOF.
// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
// On return, n >= min if and only if err == nil.
// If r returns an error having read at least min bytes, the error is dropped.
// ReadAtLeast 从 r 读到 buf, 直到读到 min 个字节。EOF error 尽在没有字节可读的情况下返回。
// 如果在读到 的字节数少于 min bytes 遇到了 EOF ，ReadAtLeast 会返回 ErrUnexpectedEOF
// 如果 min 比 buf 的长度大，ReadAtLeast 会返回 ErrShortBuffer。
// 在返回上，当且仅当 n >= min 时 err == nil。
// 如果 r 在读取至少 min 字节时返回了一个错误，则丢弃该错误。
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == EOF {
		err = ErrUnexpectedEOF
	}
	return
}

// ReadFull reads exactly len(buf) bytes from r into buf.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading some but not all the bytes,
// ReadFull returns ErrUnexpectedEOF.
// On return, n == len(buf) if and only if err == nil.
// If r returns an error having read at least len(buf) bytes, the error is dropped.
// ReadFull 将 r 中的 len(buf) 个字节准确地读入 buf. 它将返回复制的字节数和如果读到的字节较少会返回
// 错误，仅当未读取任何字节时错误为 EOF。如果在读取一些但不是所有字节时遇到 EOF，则会返回ErrUnexpectedEOF
// 在返回时，当且仅当 n == len(buf)时，err == nil
// 如果在读取至少 len(buf)个字节时 r 返回一个错误，error 会被丢弃。
func ReadFull(r Reader, buf []byte) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf))
}

// CopyN copies n bytes (or until an error) from src to dst.
// It returns the number of bytes copied and the earliest
// error encountered while copying.
// On return, written == n if and only if err == nil.
//
// If dst implements the ReaderFrom interface,
// the copy is implemented using it.
// CopyN 从 src 到 dst 拷贝 n 个字节（或直到遇到一个 error）
// 它将返回复制的字节数 和 copy 时最早遇到的错误
// 作为返回值， 当且仅当 err == nil 时候 written == n。
//
// 如果 dst 实现了 ReaderFrom interface，copy 的视线会使用它。
//

func CopyN(dst Writer, src Reader, n int64) (written int64, err error) {
	written, err = Copy(dst, LimitReader(src, n))
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early; must have been EOF.
		err = EOF
	}
	return
}

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
// Copy 从 src 拷贝数据到 dst，直到遇到 EOF 或一个错误。它返回已拷贝
// 的字节数和拷贝时最先遇到任何错误
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
// 一个成功的 Copy 返回 err == nil， 而不是 err == EOF
// 因为 Copy 被定义为从 src 读书数据直到遇到 EOF，它并不认为从 Read 读到 EOF
// 是一个错误 而上报它
//
// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).
// 如果 src 实现了 WriteTo 接口，则 copy 将被时限为调用 src.WriteTo(dst)
// 另外，如果 dst 实现了 ReaderFrom 接口，copy 将被实现为调用
// dst.ReaderFrom(src)
//
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer (if one is required) rather than allocating a
// temporary one. If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
// CopyBuffer 是和 Copy 一样的方法，除了它通过将数据暂存在提供的缓冲区外而不是申请
// 一块临时的空间。如果 buf 为 nil，将申请一个；另外如果 buf 长度为 0，则 panic
//
// If either src implements WriterTo or dst implements ReaderFrom,
// buf will not be used to perform the copy.
// 如果 src 实现了 WriteTo 或 dst 实现了 ReaderFrom，但是buf 不会用于执行 copy。
//
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}
	return copyBuffer(dst, src, buf)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// LimitReader returns a Reader that reads from r
// but stops with EOF after n bytes.
// The underlying implementation is a *LimitedReader.
// LimitReader 返回一个 从 r 读但在读 n 个字节后以 Eof 结束的 Reader
// 底层实现是一个 *LimitedReader
func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
// 一个 LimitedReader 从 R 读数据但限制读取返回的字节数为 n。
// 每次调用 Read 更新 N 以反映新的剩余数量,当 N <= 0 或底层 R 返回 EOF 时，Read 返回 EOF。
type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

// NewSectionReader returns a SectionReader that reads from r
// starting at offset off and stops with EOF after n bytes.
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader {
	var remaining int64
	const maxint64 = 1<<63 - 1
	if off <= maxint64-n {
		remaining = n + off
	} else {
		// Overflow, with no way to return error.
		// Assume we can read up to an offset of 1<<63 - 1.
		remaining = maxint64
	}
	return &SectionReader{r, off, off, remaining}
}

// SectionReader implements Read, Seek, and ReadAt on a section
// of an underlying ReaderAt.
// SectionReader 在一个底层 ReaderAt section 上 实现了 Read, Seek, and ReadAt
type SectionReader struct {
	r     ReaderAt
	base  int64
	off   int64
	limit int64
}

func (s *SectionReader) Read(p []byte) (n int, err error) {
	if s.off >= s.limit {
		return 0, EOF
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
	}
	n, err = s.r.ReadAt(p, s.off)
	s.off += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case SeekStart:
		offset += s.base
	case SeekCurrent:
		offset += s.off
	case SeekEnd:
		offset += s.limit
	}
	if offset < s.base {
		return 0, errOffset
	}
	s.off = offset
	return offset - s.base, nil
}

func (s *SectionReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= s.limit-s.base {
		return 0, EOF
	}
	off += s.base
	if max := s.limit - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = s.r.ReadAt(p, off)
		if err == nil {
			err = EOF
		}
		return n, err
	}
	return s.r.ReadAt(p, off)
}

// Size returns the size of the section in bytes.
func (s *SectionReader) Size() int64 { return s.limit - s.base }

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r Reader, w Writer) Reader {
	return &teeReader{r, w}
}

type teeReader struct {
	r Reader
	w Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

// Discard is a Writer on which all Write calls succeed
// without doing anything.
var Discard Writer = discard{}

type discard struct{}

// discard implements ReaderFrom as an optimization so Copy to
// io.Discard can avoid doing unnecessary work.
var _ ReaderFrom = discard{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discard) WriteString(s string) (int, error) {
	return len(s), nil
}

var blackHolePool = sync.Pool{
	New: func() any {
		b := make([]byte, 8192)
		return &b
	},
}

func (discard) ReadFrom(r Reader) (n int64, err error) {
	bufp := blackHolePool.Get().(*[]byte)
	readSize := 0
	for {
		readSize, err = r.Read(*bufp)
		n += int64(readSize)
		if err != nil {
			blackHolePool.Put(bufp)
			if err == EOF {
				return n, nil
			}
			return
		}
	}
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
//  NopCloser 返回一个 ReadCloser， 包装了提供的 Reader r 的无操作方法Close 方法
func NopCloser(r Reader) ReadCloser {
	return nopCloser{r}
}

type nopCloser struct {
	Reader
}

func (nopCloser) Close() error { return nil }

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == EOF {
				err = nil
			}
			return b, err
		}
	}
}
