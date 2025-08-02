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

// Package bufferpool houses zap's shared internal buffer pool. Third-party
// packages can recreate the same functionality with buffers.NewPool.
// Package bufferpool包含zap的共享内部缓冲池。第三方包可以
// 使用buffers.NewPool重新创建相同的功能。
package bufferpool

import "go.uber.org/zap/buffer" // buffer包：高性能缓冲区

var (
	_pool = buffer.NewPool() // 创建新的缓冲池
	// Get retrieves a buffer from the pool, creating one if necessary.
	// Get从池中检索缓冲区，如果需要则创建一个。
	Get = _pool.Get // 导出池的Get方法
)
