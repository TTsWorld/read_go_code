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

import "go.uber.org/zap/internal/color" // 内部颜色包

var (
	_levelToColor = map[Level]color.Color{ // 日志级别到颜色的映射
		DebugLevel:  color.Magenta, // Debug级别：紫红色
		InfoLevel:   color.Blue,    // Info级别：蓝色
		WarnLevel:   color.Yellow,  // Warn级别：黄色
		ErrorLevel:  color.Red,     // Error级别：红色
		DPanicLevel: color.Red,     // DPanic级别：红色
		PanicLevel:  color.Red,     // Panic级别：红色
		FatalLevel:  color.Red,     // Fatal级别：红色
	}
	_unknownLevelColor = color.Red // 未知级别颜色：红色

	_levelToLowercaseColorString = make(map[Level]string, len(_levelToColor)) // 级别到小写彩色字符串映射
	_levelToCapitalColorString   = make(map[Level]string, len(_levelToColor)) // 级别到大写彩色字符串映射
)

func init() { // 初始化函数
	for level, color := range _levelToColor { // 遍历级别颜色映射
		_levelToLowercaseColorString[level] = color.Add(level.String())         // 生成小写彩色字符串
		_levelToCapitalColorString[level] = color.Add(level.CapitalString())    // 生成大写彩色字符串
	}
}
