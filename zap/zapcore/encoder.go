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

import (
	"encoding/json" // json包：JSON编码和解码
	"io"            // io包：基本I/O原语
	"time"          // time包：时间处理

	"go.uber.org/zap/buffer" // buffer包：高性能缓冲区
)

// DefaultLineEnding defines the default line ending when writing logs.
// Alternate line endings specified in EncoderConfig can override this
// behavior.
// DefaultLineEnding定义写入日志时的默认行结尾。
// 在EncoderConfig中指定的替代行结尾可以覆盖此行为。
const DefaultLineEnding = "\n" // 默认的行结尾

// OmitKey defines the key to use when callers want to remove a key from log output.
// OmitKey定义当调用者想要从日志输出中删除键时使用的键。
const OmitKey = "" // 省略字段时使用的空键

// A LevelEncoder serializes a Level to a primitive type.
//
// This function must make exactly one call
// to a PrimitiveArrayEncoder's Append* method.
// LevelEncoder将Level序列化为原始类型。
//
// 此函数必须对PrimitiveArrayEncoder的Append*方法进行恰好一次调用。
type LevelEncoder func(Level, PrimitiveArrayEncoder) // 等级编码器函数类型

// LowercaseLevelEncoder serializes a Level to a lowercase string. For example,
// InfoLevel is serialized to "info".
// LowercaseLevelEncoder将Level序列化为小写字符串。
// 例如，InfoLevel被序列化为"info"。
func LowercaseLevelEncoder(l Level, enc PrimitiveArrayEncoder) { // 将等级编码为小写字符串
	enc.AppendString(l.String()) // 追加级别的字符串表示
}

// LowercaseColorLevelEncoder serializes a Level to a lowercase string and adds coloring.
// For example, InfoLevel is serialized to "info" and colored blue.
// LowercaseColorLevelEncoder将Level序列化为小写字符串并添加颜色。
// 例如，InfoLevel被序列化为"info"并着蓝色。
func LowercaseColorLevelEncoder(l Level, enc PrimitiveArrayEncoder) { // 小写并带颜色的等级编码
	s, ok := _levelToLowercaseColorString[l] // 查找级别对应的彩色字符串
	if !ok {                                 // 如果未找到
		s = _unknownLevelColor.Add(l.String()) // 使用未知级别颜色
	}
	enc.AppendString(s) // 追加彩色字符串
}

// CapitalLevelEncoder serializes a Level to an all-caps string. For example,
// InfoLevel is serialized to "INFO".
// CapitalLevelEncoder将Level序列化为全大写字符串。
// 例如，InfoLevel被序列化为"INFO"。
func CapitalLevelEncoder(l Level, enc PrimitiveArrayEncoder) { // 将等级编码为大写字符串
	enc.AppendString(l.CapitalString()) // 追加级别的大写字符串表示
}

// CapitalColorLevelEncoder serializes a Level to an all-caps string and adds color.
// For example, InfoLevel is serialized to "INFO" and colored blue.
// CapitalColorLevelEncoder将Level序列化为全大写字符串并添加颜色。
// 例如，InfoLevel被序列化为"INFO"并着蓝色。
func CapitalColorLevelEncoder(l Level, enc PrimitiveArrayEncoder) { // 大写并带颜色的等级编码
	s, ok := _levelToCapitalColorString[l] // 查找级别对应的彩色大写字符串
	if !ok {                               // 如果未找到
		s = _unknownLevelColor.Add(l.CapitalString()) // 使用未知级别颜色
	}
	enc.AppendString(s) // 追加彩色字符串
}

// UnmarshalText unmarshals text to a LevelEncoder. "capital" is unmarshaled to
// CapitalLevelEncoder, "coloredCapital" is unmarshaled to CapitalColorLevelEncoder,
// "colored" is unmarshaled to LowercaseColorLevelEncoder, and anything else
// is unmarshaled to LowercaseLevelEncoder.
func (e *LevelEncoder) UnmarshalText(text []byte) error { // 反序列化文本为等级编码器
	switch string(text) { // 根据文本选择对应编码器
	case "capital":
		*e = CapitalLevelEncoder // 全大写
	case "capitalColor":
		*e = CapitalColorLevelEncoder // 全大写并着色
	case "color":
		*e = LowercaseColorLevelEncoder // 小写并着色
	default:
		*e = LowercaseLevelEncoder // 默认小写
	}
	return nil // 无错误
}

// A TimeEncoder serializes a time.Time to a primitive type.
//
// This function must make exactly one call
// to a PrimitiveArrayEncoder's Append* method.
type TimeEncoder func(time.Time, PrimitiveArrayEncoder) // 时间编码器函数类型

// EpochTimeEncoder serializes a time.Time to a floating-point number of seconds
// since the Unix epoch.
func EpochTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) { // 将时间编码为自Unix纪元起的秒
	nanos := t.UnixNano()                        // 纳秒时间戳
	sec := float64(nanos) / float64(time.Second) // 转为秒（浮点）
	enc.AppendFloat64(sec)                       // 追加秒数
}

// EpochMillisTimeEncoder serializes a time.Time to a floating-point number of
// milliseconds since the Unix epoch.
func EpochMillisTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) { // 将时间编码为自Unix纪元起的毫秒
	nanos := t.UnixNano()                                // 纳秒时间戳
	millis := float64(nanos) / float64(time.Millisecond) // 转为毫秒（浮点）
	enc.AppendFloat64(millis)                            // 追加毫秒
}

// EpochNanosTimeEncoder serializes a time.Time to an integer number of
// nanoseconds since the Unix epoch.
func EpochNanosTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) { // 将时间编码为自Unix纪元起的纳秒
	enc.AppendInt64(t.UnixNano()) // 追加纳秒
}

func encodeTimeLayout(t time.Time, layout string, enc PrimitiveArrayEncoder) { // 按给定layout编码时间
	type appendTimeEncoder interface { // 可选接口：直接追加layout
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok { // 如果实现了AppendTimeLayout
		enc.AppendTimeLayout(t, layout) // 直接调用以避免分配
		return                          // 返回，已处理
	}

	enc.AppendString(t.Format(layout)) // 否则格式化为字符串再追加
}

// ISO8601TimeEncoder serializes a time.Time to an ISO8601-formatted string
// with millisecond precision.
//
// If enc supports AppendTimeLayout(t time.Time,layout string), it's used
// instead of appending a pre-formatted string value.
func ISO8601TimeEncoder(t time.Time, enc PrimitiveArrayEncoder) { // 按ISO8601（毫秒精度）编码
	encodeTimeLayout(t, "2006-01-02T15:04:05.000Z0700", enc) // 使用固定layout
}

// RFC3339TimeEncoder serializes a time.Time to an RFC3339-formatted string.
//
// If enc supports AppendTimeLayout(t time.Time,layout string), it's used
// instead of appending a pre-formatted string value.
func RFC3339TimeEncoder(t time.Time, enc PrimitiveArrayEncoder) { // 按RFC3339编码
	encodeTimeLayout(t, time.RFC3339, enc) // 标准layout
}

// RFC3339NanoTimeEncoder serializes a time.Time to an RFC3339-formatted string
// with nanosecond precision.
//
// If enc supports AppendTimeLayout(t time.Time,layout string), it's used
// instead of appending a pre-formatted string value.
func RFC3339NanoTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) { // 按RFC3339纳秒精度编码
	encodeTimeLayout(t, time.RFC3339Nano, enc) // 纳秒layout
}

// TimeEncoderOfLayout returns TimeEncoder which serializes a time.Time using
// given layout.
func TimeEncoderOfLayout(layout string) TimeEncoder { // 返回使用指定layout的时间编码器
	return func(t time.Time, enc PrimitiveArrayEncoder) { // 闭包捕获layout
		encodeTimeLayout(t, layout, enc) // 按layout编码
	}
}

// UnmarshalText unmarshals text to a TimeEncoder.
// "rfc3339nano" and "RFC3339Nano" are unmarshaled to RFC3339NanoTimeEncoder.
// "rfc3339" and "RFC3339" are unmarshaled to RFC3339TimeEncoder.
// "iso8601" and "ISO8601" are unmarshaled to ISO8601TimeEncoder.
// "millis" is unmarshaled to EpochMillisTimeEncoder.
// "nanos" is unmarshaled to EpochNanosEncoder.
// Anything else is unmarshaled to EpochTimeEncoder.
func (e *TimeEncoder) UnmarshalText(text []byte) error { // 文本反序列化为时间编码器
	switch string(text) { // 区分不同格式关键字
	case "rfc3339nano", "RFC3339Nano":
		*e = RFC3339NanoTimeEncoder // 纳秒精度
	case "rfc3339", "RFC3339":
		*e = RFC3339TimeEncoder // 秒级精度
	case "iso8601", "ISO8601":
		*e = ISO8601TimeEncoder // ISO8601
	case "millis":
		*e = EpochMillisTimeEncoder // 毫秒
	case "nanos":
		*e = EpochNanosTimeEncoder // 纳秒
	default:
		*e = EpochTimeEncoder // 秒（默认）
	}
	return nil // 无错误
}

// UnmarshalYAML unmarshals YAML to a TimeEncoder.
// If value is an object with a "layout" field, it will be unmarshaled to  TimeEncoder with given layout.
//
//	timeEncoder:
//	  layout: 06/01/02 03:04pm
//
// If value is string, it uses UnmarshalText.
//
//	timeEncoder: iso8601
func (e *TimeEncoder) UnmarshalYAML(unmarshal func(interface{}) error) error { // 从YAML反序列化时间编码器
	var o struct { // 支持对象形式指定layout
		Layout string `json:"layout" yaml:"layout"`
	}
	if err := unmarshal(&o); err == nil { // 若可解析为对象
		*e = TimeEncoderOfLayout(o.Layout) // 使用自定义layout
		return nil                         // 完成
	}

	var s string
	if err := unmarshal(&s); err != nil { // 尝试解析为字符串
		return err // 返回解析错误
	}
	return e.UnmarshalText([]byte(s)) // 复用文本解析
}

// UnmarshalJSON unmarshals JSON to a TimeEncoder as same way UnmarshalYAML does.
func (e *TimeEncoder) UnmarshalJSON(data []byte) error { // 以与YAML一致的逻辑从JSON反序列化
	return e.UnmarshalYAML(func(v interface{}) error { // 复用YAML分支
		return json.Unmarshal(data, v) // 解JSON到临时结构
	})
}

// A DurationEncoder serializes a time.Duration to a primitive type.
//
// This function must make exactly one call
// to a PrimitiveArrayEncoder's Append* method.
type DurationEncoder func(time.Duration, PrimitiveArrayEncoder) // 时长编码器函数类型

// SecondsDurationEncoder serializes a time.Duration to a floating-point number of seconds elapsed.
func SecondsDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) { // 按秒（浮点）编码时长
	enc.AppendFloat64(float64(d) / float64(time.Second)) // 纳秒转秒
}

// NanosDurationEncoder serializes a time.Duration to an integer number of
// nanoseconds elapsed.
func NanosDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) { // 按纳秒（整数）编码时长
	enc.AppendInt64(int64(d)) // 直接纳秒数
}

// MillisDurationEncoder serializes a time.Duration to an integer number of
// milliseconds elapsed.
func MillisDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) { // 按毫秒（整数）编码时长
	enc.AppendInt64(d.Nanoseconds() / 1e6) // 纳秒转毫秒
}

// StringDurationEncoder serializes a time.Duration using its built-in String
// method.
func StringDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) { // 使用Go内置字符串格式
	enc.AppendString(d.String()) // 例如 "1s"、"200ms"
}

// UnmarshalText unmarshals text to a DurationEncoder. "string" is unmarshaled
// to StringDurationEncoder, and anything else is unmarshaled to
// NanosDurationEncoder.
func (e *DurationEncoder) UnmarshalText(text []byte) error { // 文本反序列化为时长编码器
	switch string(text) { // 依据关键字选择
	case "string":
		*e = StringDurationEncoder // 使用字符串格式
	case "nanos":
		*e = NanosDurationEncoder // 纳秒
	case "ms":
		*e = MillisDurationEncoder // 毫秒
	default:
		*e = SecondsDurationEncoder // 秒（默认）
	}
	return nil // 无错误
}

// A CallerEncoder serializes an EntryCaller to a primitive type.
//
// This function must make exactly one call
// to a PrimitiveArrayEncoder's Append* method.
type CallerEncoder func(EntryCaller, PrimitiveArrayEncoder) // 调用方编码器函数类型

// FullCallerEncoder serializes a caller in /full/path/to/package/file:line
// format.
func FullCallerEncoder(caller EntryCaller, enc PrimitiveArrayEncoder) { // 全路径 caller 编码
	// TODO: consider using a byte-oriented API to save an allocation.
	enc.AppendString(caller.String()) // 例如 /path/to/pkg/file.go:123
}

// ShortCallerEncoder serializes a caller in package/file:line format, trimming
// all but the final directory from the full path.
func ShortCallerEncoder(caller EntryCaller, enc PrimitiveArrayEncoder) { // 短路径 caller 编码
	// TODO: consider using a byte-oriented API to save an allocation.
	enc.AppendString(caller.TrimmedPath()) // 例如 pkg/file.go:123
}

// UnmarshalText unmarshals text to a CallerEncoder. "full" is unmarshaled to
// FullCallerEncoder and anything else is unmarshaled to ShortCallerEncoder.
func (e *CallerEncoder) UnmarshalText(text []byte) error { // 文本反序列化为调用方编码器
	switch string(text) { // full 或默认
	case "full":
		*e = FullCallerEncoder // 使用全路径
	default:
		*e = ShortCallerEncoder // 使用短路径
	}
	return nil // 无错误
}

// A NameEncoder serializes a period-separated logger name to a primitive
// type.
//
// This function must make exactly one call
// to a PrimitiveArrayEncoder's Append* method.
type NameEncoder func(string, PrimitiveArrayEncoder) // 日志器名称编码器类型

// FullNameEncoder serializes the logger name as-is.
func FullNameEncoder(loggerName string, enc PrimitiveArrayEncoder) { // 原样输出名称
	enc.AppendString(loggerName) // 追加名称
}

// UnmarshalText unmarshals text to a NameEncoder. Currently, everything is
// unmarshaled to FullNameEncoder.
func (e *NameEncoder) UnmarshalText(text []byte) error { // 文本反序列化为名称编码器
	switch string(text) { // 当前仅支持full
	case "full":
		*e = FullNameEncoder // 全名称
	default:
		*e = FullNameEncoder // 默认同full
	}
	return nil // 无错误
}

// An EncoderConfig allows users to configure the concrete encoders supplied by
// zapcore.
type EncoderConfig struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey     string `json:"messageKey" yaml:"messageKey"`         // 消息字段键
	LevelKey       string `json:"levelKey" yaml:"levelKey"`             // 级别字段键
	TimeKey        string `json:"timeKey" yaml:"timeKey"`               // 时间字段键
	NameKey        string `json:"nameKey" yaml:"nameKey"`               // 日志器名称键
	CallerKey      string `json:"callerKey" yaml:"callerKey"`           // 调用方字段键
	FunctionKey    string `json:"functionKey" yaml:"functionKey"`       // 函数名字段键
	StacktraceKey  string `json:"stacktraceKey" yaml:"stacktraceKey"`   // 堆栈字段键
	SkipLineEnding bool   `json:"skipLineEnding" yaml:"skipLineEnding"` // 是否跳过行结尾
	LineEnding     string `json:"lineEnding" yaml:"lineEnding"`         // 行结尾字符串
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`       // 级别编码器
	EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`         // 时间编码器
	EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"` // 时长编码器
	EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`     // 调用方编码器
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"` // 名称编码器
	// Configure the encoder for interface{} type objects.
	// If not provided, objects are encoded using json.Encoder
	NewReflectedEncoder func(io.Writer) ReflectedEncoder `json:"-" yaml:"-"` // 自定义反射编码器
	// Configures the field separator used by the console encoder. Defaults
	// to tab.
	ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"` // 控制台分隔符
}

// ObjectEncoder is a strongly-typed, encoding-agnostic interface for adding a
// map- or struct-like object to the logging context. Like maps, ObjectEncoders
// aren't safe for concurrent use (though typical use shouldn't require locks).
type ObjectEncoder interface {
	// Logging-specific marshalers.
	AddArray(key string, marshaler ArrayMarshaler) error   // 添加数组字段
	AddObject(key string, marshaler ObjectMarshaler) error // 添加对象字段

	// Built-in types.
	AddBinary(key string, value []byte)          // 任意字节
	AddByteString(key string, value []byte)      // UTF-8字节
	AddBool(key string, value bool)              // 布尔
	AddComplex128(key string, value complex128)  // 复数128
	AddComplex64(key string, value complex64)    // 复数64
	AddDuration(key string, value time.Duration) // 时长
	AddFloat64(key string, value float64)        // 浮点64
	AddFloat32(key string, value float32)        // 浮点32
	AddInt(key string, value int)                // 整型
	AddInt64(key string, value int64)            // 整型64
	AddInt32(key string, value int32)            // 整型32
	AddInt16(key string, value int16)            // 整型16
	AddInt8(key string, value int8)              // 整型8
	AddString(key, value string)                 // 字符串
	AddTime(key string, value time.Time)         // 时间
	AddUint(key string, value uint)              // 无符号
	AddUint64(key string, value uint64)          // 无符号64
	AddUint32(key string, value uint32)          // 无符号32
	AddUint16(key string, value uint16)          // 无符号16
	AddUint8(key string, value uint8)            // 无符号8
	AddUintptr(key string, value uintptr)        // 指针大小无符号

	// AddReflected uses reflection to serialize arbitrary objects, so it can be
	// slow and allocation-heavy.
	AddReflected(key string, value interface{}) error // 通过反射添加任意对象
	// OpenNamespace opens an isolated namespace where all subsequent fields will
	// be added. Applications can use namespaces to prevent key collisions when
	// injecting loggers into sub-components or third-party libraries.
	OpenNamespace(key string) // 打开子命名空间
}

// ArrayEncoder is a strongly-typed, encoding-agnostic interface for adding
// array-like objects to the logging context. Of note, it supports mixed-type
// arrays even though they aren't typical in Go. Like slices, ArrayEncoders
// aren't safe for concurrent use (though typical use shouldn't require locks).
type ArrayEncoder interface {
	// Built-in types.
	PrimitiveArrayEncoder // 内置类型追加器子集

	// Time-related types.
	AppendDuration(time.Duration) // 追加时长
	AppendTime(time.Time)         // 追加时间

	// Logging-specific marshalers.
	AppendArray(ArrayMarshaler) error   // 追加数组类型
	AppendObject(ObjectMarshaler) error // 追加对象类型

	// AppendReflected uses reflection to serialize arbitrary objects, so it's
	// slow and allocation-heavy.
	AppendReflected(value interface{}) error // 通过反射追加任意对象
}

// PrimitiveArrayEncoder is the subset of the ArrayEncoder interface that deals
// only in Go's built-in types. It's included only so that Duration- and
// TimeEncoders cannot trigger infinite recursion.
type PrimitiveArrayEncoder interface {
	// Built-in types.
	AppendBool(bool)             // 布尔
	AppendByteString([]byte)     // for UTF-8 encoded bytes // UTF-8字节
	AppendComplex128(complex128) // 复数128
	AppendComplex64(complex64)   // 复数64
	AppendFloat64(float64)       // 浮点64
	AppendFloat32(float32)       // 浮点32
	AppendInt(int)               // 整型
	AppendInt64(int64)           // 整型64
	AppendInt32(int32)           // 整型32
	AppendInt16(int16)           // 整型16
	AppendInt8(int8)             // 整型8
	AppendString(string)         // 字符串
	AppendUint(uint)             // 无符号
	AppendUint64(uint64)         // 无符号64
	AppendUint32(uint32)         // 无符号32
	AppendUint16(uint16)         // 无符号16
	AppendUint8(uint8)           // 无符号8
	AppendUintptr(uintptr)       // 指针大小无符号
}

// Encoder is a format-agnostic interface for all log entry marshalers. Since
// log encoders don't need to support the same wide range of use cases as
// general-purpose marshalers, it's possible to make them faster and
// lower-allocation.
//
// Implementations of the ObjectEncoder interface's methods can, of course,
// freely modify the receiver. However, the Clone and EncodeEntry methods will
// be called concurrently and shouldn't modify the receiver.
type Encoder interface {
	ObjectEncoder

	// Clone copies the encoder, ensuring that adding fields to the copy doesn't
	// affect the original.
	Clone() Encoder // 拷贝编码器（追加字段互不影响）

	// EncodeEntry encodes an entry and fields, along with any accumulated
	// context, into a byte buffer and returns it. Any fields that are empty,
	// including fields on the `Entry` type, should be omitted.
	EncodeEntry(Entry, []Field) (*buffer.Buffer, error) // 编码日志条目并返回字节缓冲
}
