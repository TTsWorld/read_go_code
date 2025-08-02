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

// Package exit provides stubs so that unit tests can exercise code that calls
// os.Exit(1).
// Package exit提供存根，以便单元测试可以执行调用os.Exit(1)的代码。
package exit

import "os" // os包：操作系统接口

var _exit = os.Exit // 可替换的退出函数

// With terminates the process by calling os.Exit(code). If the package is
// stubbed, it instead records a call in the testing spy.
// With通过调用os.Exit(code)终止进程。如果包被存根化，
// 它会在测试间谍中记录调用。
func With(code int) {
	_exit(code) // 调用退出函数
}

// A StubbedExit is a testing fake for os.Exit.
// StubbedExit是os.Exit的测试伪造。
type StubbedExit struct {
	Exited bool          // 是否已退出
	Code   int           // 退出代码
	prev   func(code int) // 先前的退出函数
}

// Stub substitutes a fake for the call to os.Exit(1).
// Stub为调用os.Exit(1)替换伪造。
func Stub() *StubbedExit {
	s := &StubbedExit{prev: _exit} // 保存先前的退出函数
	_exit = s.exit               // 设置为存根退出函数
	return s                     // 返回存根
}

// WithStub runs the supplied function with Exit stubbed. It returns the stub
// used, so that users can test whether the process would have crashed.
// WithStub在Exit被存根化的情况下运行提供的函数。它返回使用的存根，
// 以便用户可以测试进程是否会崩溃。
func WithStub(f func()) *StubbedExit {
	s := Stub()       // 创建存根
	defer s.Unstub()  // 延迟恢复
	f()               // 执行函数
	return s          // 返回存根
}

// Unstub restores the previous exit function.
// Unstub恢复先前的退出函数。
func (se *StubbedExit) Unstub() {
	_exit = se.prev // 恢复先前的退出函数
}

func (se *StubbedExit) exit(code int) { // 存根退出函数
	se.Exited = true  // 标记为已退出
	se.Code = code    // 记录退出代码
}
