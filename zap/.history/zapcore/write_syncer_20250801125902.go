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
	"io"   // io包：基本I/O原语
	"sync" // sync包：同步原语

	"go.uber.org/multierr" // multierr包：多重错误处理
)

// A WriteSyncer is an io.Writer that can also flush any buffered data. Note
// that *os.File (and thus, os.Stderr and os.Stdout) implement WriteSyncer.
// WriteSyncer是一个io.Writer，也可以刷新任何缓冲数据。
// 注意*os.File（因此os.Stderr和os.Stdout）实现了WriteSyncer。
type WriteSyncer interface {
	io.Writer        // 嵌入Writer接口
	Sync() error     // 同步方法，刷新缓冲数据
}

// AddSync converts an io.Writer to a WriteSyncer. It attempts to be
// intelligent: if the concrete type of the io.Writer implements WriteSyncer,
// we'll use the existing Sync method. If it doesn't, we'll add a no-op Sync.
// AddSync将io.Writer转换为WriteSyncer。它尝试智能处理：
// 如果io.Writer的具体类型实现了WriteSyncer，我们将使用现有的Sync方法。
// 如果没有，我们将添加一个无操作的Sync。
func AddSync(w io.Writer) WriteSyncer {
	switch w := w.(type) {    // 类型断言检查
	case WriteSyncer:         // 如果已经是WriteSyncer
		return w              // 直接返回
	default:                  // 否则
		return writerWrapper{w} // 包装为writerWrapper
	}
}

type lockedWriteSyncer struct { // 带锁的WriteSyncer
	sync.Mutex         // 嵌入互斥锁
	ws WriteSyncer     // 底层WriteSyncer
}

// Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use. In
// particular, *os.Files must be locked before use.
// Lock用互斥锁包装WriteSyncer以使其并发安全。
// 特别是，*os.Files在使用前必须加锁。
func Lock(ws WriteSyncer) WriteSyncer {
	if _, ok := ws.(*lockedWriteSyncer); ok { // 如果已经是带锁的
		// no need to layer on another lock
		// 无需再加一层锁
		return ws                             // 直接返回
	}
	return &lockedWriteSyncer{ws: ws}         // 包装为带锁版本
}

func (s *lockedWriteSyncer) Write(bs []byte) (int, error) { // Write方法（带锁）
	s.Lock()                      // 加锁
	n, err := s.ws.Write(bs)      // 调用底层写入
	s.Unlock()                    // 解锁
	return n, err                 // 返回结果
}

func (s *lockedWriteSyncer) Sync() error { // Sync方法（带锁）
	s.Lock()                // 加锁
	err := s.ws.Sync()      // 调用底层同步
	s.Unlock()              // 解锁
	return err              // 返回错误
}

type writerWrapper struct { // Writer包装器
	io.Writer               // 嵌入Writer
}

func (w writerWrapper) Sync() error { // 无操作的Sync实现
	return nil                        // 直接返回nil
}

type multiWriteSyncer []WriteSyncer

// NewMultiWriteSyncer creates a WriteSyncer that duplicates its writes
// and sync calls, much like io.MultiWriter.
func NewMultiWriteSyncer(ws ...WriteSyncer) WriteSyncer {
	if len(ws) == 1 {
		return ws[0]
	}
	return multiWriteSyncer(ws)
}

// See https://golang.org/src/io/multi.go
// When not all underlying syncers write the same number of bytes,
// the smallest number is returned even though Write() is called on
// all of them.
func (ws multiWriteSyncer) Write(p []byte) (int, error) {
	var writeErr error
	nWritten := 0
	for _, w := range ws {
		n, err := w.Write(p)
		writeErr = multierr.Append(writeErr, err)
		if nWritten == 0 && n != 0 {
			nWritten = n
		} else if n < nWritten {
			nWritten = n
		}
	}
	return nWritten, writeErr
}

func (ws multiWriteSyncer) Sync() error {
	var err error
	for _, w := range ws {
		err = multierr.Append(err, w.Sync())
	}
	return err
}
