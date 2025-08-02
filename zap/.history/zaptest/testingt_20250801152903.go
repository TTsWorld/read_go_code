// Copyright (c) 2017 Uber Technologies, Inc.
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

package zaptest // zaptest包：提供测试日志输出的各种助手函数

// TestingT is a subset of the API provided by all *testing.T and *testing.B
// objects.
// TestingT是所有*testing.T和*testing.B对象提供的API的子集。
type TestingT interface {
	// Logs the given message without failing the test.
	// 记录给定的消息而不使测试失败。
	Logf(string, ...interface{})

	// Logs the given message and marks the test as failed.
	// 记录给定的消息并标记测试为失败。
	Errorf(string, ...interface{})

	// Marks the test as failed.
	// 标记测试为失败。
	Fail()

	// Returns true if the test has been marked as failed.
	// 如果测试已被标记为失败，则返回true。
	Failed() bool

	// Returns the name of the test.
	// 返回测试的名称。
	Name() string

	// Marks the test as failed and stops execution of that test.
	// 标记测试为失败并停止该测试的执行。
	FailNow()
}

// Note: We currently only rely on Logf. We are including Errorf and FailNow
// in the interface in anticipation of future need since we can't extend the
// interface without a breaking change.
// 注意：我们目前只依赖Logf。我们在接口中包含Errorf和FailNow是为了预期未来的需要，
// 因为我们无法在不破坏更改的情况下扩展接口。
