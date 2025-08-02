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

package zap // zap包：提供快速、结构化、分级日志记录

import (
	"go.uber.org/zap/internal/pool" // 内部对象池
	"go.uber.org/zap/zapcore"       // zapcore包：核心接口和实现
)

var _errArrayElemPool = pool.New(func() *errArrayElem { // 错误数组元素对象池
	return &errArrayElem{} // 返回新的错误数组元素
})

// Error is shorthand for the common idiom NamedError("error", err).
// Error是常用惯用法NamedError("error", err)的简写。
func Error(err error) Field {
	return NamedError("error", err) // 调用NamedError使用默认键名"error"
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
// NamedError构造一个在提供的键下延迟存储err.Error()的字段。
// 同时实现fmt.Formatter的错误（如github.com/pkg/errors生成的错误）
// 也将在key+"Verbose"下存储其详细表示。如果传入nil错误，字段为无操作。
//
// 对于键简单为"error"的常见情况，Error函数更简短且重复性较少。
func NamedError(key string, err error) Field {
	if err == nil { // 如果错误为nil
		return Skip() // 返回跳过字段
	}
	return Field{Key: key, Type: zapcore.ErrorType, Interface: err} // 返回错误类型字段
}

type errArray []error // 错误数组类型

func (errs errArray) MarshalLogArray(arr zapcore.ArrayEncoder) error { // 实现ArrayMarshaler接口
	for i := range errs { // 遍历所有错误
		if errs[i] == nil { // 如果错误为nil
			continue // 跳过
		}
		// To represent each error as an object with an "error" attribute and
		// potentially an "errorVerbose" attribute, we need to wrap it in a
		// type that implements LogObjectMarshaler. To prevent this from
		// allocating, pool the wrapper type.
		// 为了将每个错误表示为具有"error"属性和潜在"errorVerbose"属性的对象，
		// 我们需要将其包装在实现LogObjectMarshaler的类型中。
		// 为了防止分配，池化包装器类型。
		elem := _errArrayElemPool.Get() // 从对象池获取元素
		elem.error = errs[i]            // 设置错误
		err := arr.AppendObject(elem)   // 将对象添加到数组编码器
		elem.error = nil                // 清空错误（避免内存泄漏）
		_errArrayElemPool.Put(elem)     // 将元素返回对象池
		if err != nil {                 // 如果添加失败
			return err // 返回错误
		}
	}
	return nil // 成功返回nil
}

type errArrayElem struct { // 错误数组元素结构体
	error // 嵌入错误接口
}

func (e *errArrayElem) MarshalLogObject(enc zapcore.ObjectEncoder) error { // 实现ObjectMarshaler接口
	// Re-use the error field's logic, which supports non-standard error types.
	// 重用错误字段的逻辑，支持非标准错误类型。
	Error(e.error).AddTo(enc) // 将错误字段添加到编码器
	return nil                // 返回nil表示成功
}
