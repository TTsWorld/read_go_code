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

// Package pool provides internal pool utilities.
// Package pool提供内部池工具。
package pool

import (
	"sync" // sync包：同步原语
)

// A Pool is a generic wrapper around [sync.Pool] to provide strongly-typed
// object pooling.
//
// Note that SA6002 (ref: https://staticcheck.io/docs/checks/#SA6002) will
// not be detected, so all internal pool use must take care to only store
// pointer types.
// Pool是sync.Pool的泛型包装器，用于提供强类型对象池。
//
// 注意SA6002（参考：https://staticcheck.io/docs/checks/#SA6002）
// 不会被检测到，因此所有内部池使用必须注意只存储指针类型。
type Pool[T any] struct {
	pool sync.Pool // 底层同步池
}

// New returns a new [Pool] for T, and will use fn to construct new Ts when
// the pool is empty.
// New返回T的新Pool，当池为空时将使用fn构造新的T。
func New[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{ // 创建同步池
			New: func() any { // 设置新对象构造函数
				return fn() // 调用用户提供的构造函数
			},
		},
	}
}

// Get gets a T from the pool, or creates a new one if the pool is empty.
// Get从池中获取T，如果池为空则创建新的。
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T) // 从池中获取并类型断言
}

// Put returns x into the pool.
// Put将x返回到池中。
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x) // 将对象放回池中
}
