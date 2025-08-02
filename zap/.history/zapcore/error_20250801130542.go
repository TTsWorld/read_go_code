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

package zapcore // zapcore包：定义zap日志库的核心低层接口

import (
	"fmt"     // fmt包：格式化I/O
	"reflect" // reflect包：运行时反射

	"go.uber.org/zap/internal/pool" // 内部对象池
)

// Encodes the given error into fields of an object. A field with the given
// name is added for the error message.
//
// If the error implements fmt.Formatter, a field with the name ${key}Verbose
// is also added with the full verbose error message.
//
// Finally, if the error implements errorGroup (from go.uber.org/multierr) or
// causer (from github.com/pkg/errors), a ${key}Causes field is added with an
// array of objects containing the errors this error was comprised of.
//
//	{
//	  "error": err.Error(),
//	  "errorVerbose": fmt.Sprintf("%+v", err),
//	  "errorCauses": [
//	    ...
//	  ],
//	}
// encodeError将给定的错误编码到对象的字段中。为错误消息添加具有给定名称的字段。
//
// 如果错误实现了fmt.Formatter，还会添加名为${key}Verbose的字段，
// 包含完整的详细错误消息。
//
// 最后，如果错误实现了errorGroup（来自go.uber.org/multierr）或
// causer（来自github.com/pkg/errors），会添加${key}Causes字段，
// 包含组成此错误的错误对象数组。
func encodeError(key string, err error, enc ObjectEncoder) (retErr error) {
	// Try to capture panics (from nil references or otherwise) when calling
	// the Error() method
	// 尝试捕获调用Error()方法时的panic（来自nil引用或其他原因）
	defer func() {
		if rerr := recover(); rerr != nil { // 如果发生panic
			// If it's a nil pointer, just say "<nil>". The likeliest causes are a
			// error that fails to guard against nil or a nil pointer for a
			// value receiver, and in either case, "<nil>" is a nice result.
			// 如果是nil指针，只需说"<nil>"。最可能的原因是
			// 错误未能防范nil或值接收器的nil指针，
			// 无论哪种情况，"<nil>"都是很好的结果。
			if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
				enc.AddString(key, "<nil>") // 添加"<nil>"字符串
				return
			}

			retErr = fmt.Errorf("PANIC=%v", rerr) // 格式化panic错误
		}
	}()

	basic := err.Error()          // 获取基本错误消息
	enc.AddString(key, basic)     // 添加基本错误字段

	switch e := err.(type) {      // 类型断言检查错误类型
	case errorGroup:              // 如果是错误组
		return enc.AddArray(key+"Causes", errArray(e.Errors())) // 添加原因数组
	case fmt.Formatter:           // 如果实现了格式化器
		verbose := fmt.Sprintf("%+v", e) // 获取详细错误信息
		if verbose != basic {     // 如果详细信息与基本信息不同
			// This is a rich error type, like those produced by
			// github.com/pkg/errors.
			// 这是丰富的错误类型，如github.com/pkg/errors产生的错误。
			enc.AddString(key+"Verbose", verbose) // 添加详细错误字段
		}
	}
	return nil // 返回nil表示成功
}

type errorGroup interface {
	// Provides read-only access to the underlying list of errors, preferably
	// without causing any allocs.
	Errors() []error
}

// Note that errArray and errArrayElem are very similar to the version
// implemented in the top-level error.go file. We can't re-use this because
// that would require exporting errArray as part of the zapcore API.

// Encodes a list of errors using the standard error encoding logic.
type errArray []error

func (errs errArray) MarshalLogArray(arr ArrayEncoder) error {
	for i := range errs {
		if errs[i] == nil {
			continue
		}

		el := newErrArrayElem(errs[i])
		err := arr.AppendObject(el)
		el.Free()
		if err != nil {
			return err
		}
	}
	return nil
}

var _errArrayElemPool = pool.New(func() *errArrayElem {
	return &errArrayElem{}
})

// Encodes any error into a {"error": ...} re-using the same errors logic.
//
// May be passed in place of an array to build a single-element array.
type errArrayElem struct{ err error }

func newErrArrayElem(err error) *errArrayElem {
	e := _errArrayElemPool.Get()
	e.err = err
	return e
}

func (e *errArrayElem) MarshalLogArray(arr ArrayEncoder) error {
	return arr.AppendObject(e)
}

func (e *errArrayElem) MarshalLogObject(enc ObjectEncoder) error {
	return encodeError("error", e.err, enc)
}

func (e *errArrayElem) Free() {
	e.err = nil
	_errArrayElemPool.Put(e)
}
