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

package observer // observer包：提供测试观察者核心实现

import "go.uber.org/zap/zapcore" // zapcore包：核心接口和实现

// A LoggedEntry is an encoding-agnostic representation of a log message.
// Field availability is context dependent.
// LoggedEntry是日志消息的编码无关表示。
// 字段可用性取决于上下文。
type LoggedEntry struct {
	zapcore.Entry                 // 嵌入zapcore.Entry，包含基本日志信息
	Context       []zapcore.Field // 上下文字段列表
}

// ContextMap returns a map for all fields in Context.
// ContextMap返回Context中所有字段的映射。
func (e LoggedEntry) ContextMap() map[string]interface{} {
	encoder := zapcore.NewMapObjectEncoder() // 创建映射对象编码器
	for _, f := range e.Context {            // 遍历上下文字段
		f.AddTo(encoder) // 将字段添加到编码器
	}
	return encoder.Fields // 返回编码器的字段映射
}
