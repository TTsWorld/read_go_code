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

// ObjectMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
//
// Note: ObjectMarshaler is only used when zap.Object is used or when
// passed directly to zap.Any. It is not used when reflection-based
// encoding is used.
// ObjectMarshaler允许用户定义的类型高效地将自己添加到日志上下文中，
// 并选择性地省略不应包含在日志中的信息（例如密码）。
//
// 注意：ObjectMarshaler仅在使用zap.Object或直接传递给zap.Any时使用。
// 在使用基于反射的编码时不使用。
type ObjectMarshaler interface {
	MarshalLogObject(ObjectEncoder) error // 序列化日志对象
}

// ObjectMarshalerFunc is a type adapter that turns a function into an
// ObjectMarshaler.
// ObjectMarshalerFunc是将函数转换为ObjectMarshaler的类型适配器。
type ObjectMarshalerFunc func(ObjectEncoder) error

// MarshalLogObject calls the underlying function.
// MarshalLogObject调用底层函数。
func (f ObjectMarshalerFunc) MarshalLogObject(enc ObjectEncoder) error {
	return f(enc) // 调用函数
}

// ArrayMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
//
// Note: ArrayMarshaler is only used when zap.Array is used or when
// passed directly to zap.Any. It is not used when reflection-based
// encoding is used.
// ArrayMarshaler允许用户定义的类型高效地将自己添加到日志上下文中，
// 并选择性地省略不应包含在日志中的信息（例如密码）。
//
// 注意：ArrayMarshaler仅在使用zap.Array或直接传递给zap.Any时使用。
// 在使用基于反射的编码时不使用。
type ArrayMarshaler interface {
	MarshalLogArray(ArrayEncoder) error // 序列化日志数组
}

// ArrayMarshalerFunc is a type adapter that turns a function into an
// ArrayMarshaler.
// ArrayMarshalerFunc是将函数转换为ArrayMarshaler的类型适配器。
type ArrayMarshalerFunc func(ArrayEncoder) error

// MarshalLogArray calls the underlying function.
// MarshalLogArray调用底层函数。
func (f ArrayMarshalerFunc) MarshalLogArray(enc ArrayEncoder) error {
	return f(enc) // 调用函数
}
