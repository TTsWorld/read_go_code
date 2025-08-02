// Copyright (c) 2022 Uber Technologies, Inc.
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

// Package internal and its subpackages hold types and functionality
// that are not part of Zap's public API.
// Package internal及其子包包含不属于Zap公共API的类型和功能。
package internal

import "go.uber.org/zap/zapcore" // zapcore包：核心接口和实现

// LeveledEnabler is an interface satisfied by LevelEnablers that are able to
// report their own level.
//
// This interface is defined to use more conveniently in tests and non-zapcore
// packages.
// This cannot be imported from zapcore because of the cyclic dependency.
// LeveledEnabler是一个由能够报告自己级别的LevelEnabler满足的接口。
//
// 此接口的定义是为了在测试和非zapcore包中更方便地使用。
// 由于循环依赖，这不能从zapcore导入。
type LeveledEnabler interface {
	zapcore.LevelEnabler // 嵌入LevelEnabler接口

	Level() zapcore.Level // 返回当前级别
}
