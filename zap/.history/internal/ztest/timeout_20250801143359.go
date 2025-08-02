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

package ztest // ztest包：测试工具

import (
	"log"    // log包：日志记录
	"os"     // os包：操作系统接口
	"strconv" // strconv包：字符串转换
	"time"   // time包：时间处理
)

var _timeoutScale = 1.0 // 超时缩放系数

// Timeout scales the provided duration by $TEST_TIMEOUT_SCALE.
// Timeout根据$TEST_TIMEOUT_SCALE缩放提供的持续时间。
func Timeout(base time.Duration) time.Duration {
	return time.Duration(float64(base) * _timeoutScale) // 应用缩放系数
}

// Sleep scales the sleep duration by $TEST_TIMEOUT_SCALE.
// Sleep根据$TEST_TIMEOUT_SCALE缩放睡眠持续时间。
func Sleep(base time.Duration) {
	time.Sleep(Timeout(base)) // 睡眠缩放后的时间
}

// Initialize checks the environment and alters the timeout scale accordingly.
// It returns a function to undo the scaling.
// Initialize检查环境并相应地更改超时缩放。它返回一个撤销缩放的函数。
func Initialize(factor string) func() {
	fv, err := strconv.ParseFloat(factor, 64) // 解析浮点数
	if err != nil {                           // 如果解析失败
		panic(err) // 触发panic
	}
	original := _timeoutScale // 保存原始值
	_timeoutScale = fv        // 设置新的缩放系数
	return func() { _timeoutScale = original } // 返回恢复函数
}

func init() { // 包初始化函数
	if v := os.Getenv("TEST_TIMEOUT_SCALE"); v != "" { // 检查环境变量
		Initialize(v)                                  // 初始化缩放系数
		log.Printf("Scaling timeouts by %vx.\n", _timeoutScale) // 记录缩放信息
	}
}
