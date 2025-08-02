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

import "go.uber.org/multierr" // multierr包：多重错误处理

type hooked struct { // 带钩子的Core
	Core                      // 嵌入Core
	funcs []func(Entry) error // 钩子函数切片
}

var (
	_ Core           = (*hooked)(nil) // 确保hooked实现Core接口
	_ leveledEnabler = (*hooked)(nil) // 确保hooked实现leveledEnabler接口
)

// RegisterHooks wraps a Core and runs a collection of user-defined callback
// hooks each time a message is logged. Execution of the callbacks is blocking.
//
// This offers users an easy way to register simple callbacks (e.g., metrics
// collection) without implementing the full Core interface.
// RegisterHooks包装一个Core，并在每次记录消息时运行用户定义的回调钩子集合。
// 回调的执行是阻塞的。
//
// 这为用户提供了一种简单的方式来注册简单的回调（例如指标收集），
// 而无需实现完整的Core接口。
func RegisterHooks(core Core, hooks ...func(Entry) error) Core {
	funcs := append([]func(Entry) error{}, hooks...) // 复制钩子函数
	return &hooked{ // 返回带钩子的Core
		Core:  core,  // 设置底层Core
		funcs: funcs, // 设置钩子函数
	}
}

func (h *hooked) Level() Level {
	return LevelOf(h.Core)
}

func (h *hooked) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	// Let the wrapped Core decide whether to log this message or not. This
	// also gives the downstream a chance to register itself directly with the
	// CheckedEntry.
	if downstream := h.Core.Check(ent, ce); downstream != nil {
		return downstream.AddCore(ent, h)
	}
	return ce
}

func (h *hooked) With(fields []Field) Core {
	return &hooked{
		Core:  h.Core.With(fields),
		funcs: h.funcs,
	}
}

func (h *hooked) Write(ent Entry, _ []Field) error {
	// Since our downstream had a chance to register itself directly with the
	// CheckedMessage, we don't need to call it here.
	var err error
	for i := range h.funcs {
		err = multierr.Append(err, h.funcs[i](ent))
	}
	return err
}
