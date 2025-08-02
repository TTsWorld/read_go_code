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

package zap // zap包：提供快速、结构化、分级日志记录

import (
	"flag" // flag包：命令行标志解析

	"go.uber.org/zap/zapcore" // zapcore包：核心接口和实现
)

// LevelFlag uses the standard library's flag.Var to declare a global flag
// with the specified name, default, and usage guidance. The returned value is
// a pointer to the value of the flag.
//
// If you don't want to use the flag package's global state, you can use any
// non-nil *Level as a flag.Value with your own *flag.FlagSet.
// LevelFlag使用标准库的flag.Var声明一个具有指定名称、默认值和使用指导的全局标志。
// 返回值是指向标志值的指针。
//
// 如果不想使用flag包的全局状态，可以将任何非空*Level作为flag.Value
// 与自己的*flag.FlagSet一起使用。
func LevelFlag(name string, defaultLevel zapcore.Level, usage string) *zapcore.Level {
	lvl := defaultLevel         // 初始化级别变量
	flag.Var(&lvl, name, usage) // 注册全局标志
	return &lvl                 // 返回级别指针
}
