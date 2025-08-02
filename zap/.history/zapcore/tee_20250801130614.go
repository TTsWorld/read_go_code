// Copyright (c) 2016-2022 Uber Technologies, Inc.
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

type multiCore []Core // 多路Core，实现日志条目分发到多个Core

var (
	_ leveledEnabler = multiCore(nil) // 确保multiCore实现leveledEnabler接口
	_ Core           = multiCore(nil) // 确保multiCore实现Core接口
)

// NewTee creates a Core that duplicates log entries into two or more
// underlying Cores.
//
// Calling it with a single Core returns the input unchanged, and calling
// it with no input returns a no-op Core.
// NewTee创建一个将日志条目复制到两个或多个底层Core的Core。
//
// 使用单个Core调用时返回输入不变，不使用输入调用时返回无操作Core。
func NewTee(cores ...Core) Core {
	switch len(cores) {   // 根据Core数量选择策略
	case 0:               // 没有Core
		return NewNopCore() // 返回无操作Core
	case 1:               // 只有一个Core
		return cores[0]     // 直接返回该Core
	default:              // 多个Core
		return multiCore(cores) // 返回多路Core
	}
}

func (mc multiCore) With(fields []Field) Core {
	clone := make(multiCore, len(mc))
	for i := range mc {
		clone[i] = mc[i].With(fields)
	}
	return clone
}

func (mc multiCore) Level() Level {
	minLvl := _maxLevel // mc is never empty
	for i := range mc {
		if lvl := LevelOf(mc[i]); lvl < minLvl {
			minLvl = lvl
		}
	}
	return minLvl
}

func (mc multiCore) Enabled(lvl Level) bool {
	for i := range mc {
		if mc[i].Enabled(lvl) {
			return true
		}
	}
	return false
}

func (mc multiCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	for i := range mc {
		ce = mc[i].Check(ent, ce)
	}
	return ce
}

func (mc multiCore) Write(ent Entry, fields []Field) error {
	var err error
	for i := range mc {
		err = multierr.Append(err, mc[i].Write(ent, fields))
	}
	return err
}

func (mc multiCore) Sync() error {
	var err error
	for i := range mc {
		err = multierr.Append(err, mc[i].Sync())
	}
	return err
}
