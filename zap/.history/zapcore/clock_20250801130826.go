// Copyright (c) 2021 Uber Technologies, Inc.
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

import "time" // time包：时间处理

// DefaultClock is the default clock used by Zap in operations that require
// time. This clock uses the system clock for all operations.
// DefaultClock是Zap在需要时间的操作中使用的默认时钟。
// 此时钟在所有操作中都使用系统时钟。
var DefaultClock = systemClock{}

// Clock is a source of time for logged entries.
// Clock是日志条目的时间源。
type Clock interface {
	// Now returns the current local time.
	// Now返回当前本地时间。
	Now() time.Time

	// NewTicker returns *time.Ticker that holds a channel
	// that delivers "ticks" of a clock.
	// NewTicker返回*time.Ticker，它持有一个传递时钟"滴答"的通道。
	NewTicker(time.Duration) *time.Ticker
}

// systemClock implements default Clock that uses system time.
// systemClock实现使用系统时间的默认Clock。
type systemClock struct{}

func (systemClock) Now() time.Time { // 获取当前时间
	return time.Now() // 返回系统当前时间
}

func (systemClock) NewTicker(duration time.Duration) *time.Ticker { // 创建定时器
	return time.NewTicker(duration) // 返回系统定时器
}
