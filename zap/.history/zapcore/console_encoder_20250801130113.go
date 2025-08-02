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
	"fmt" // fmt包：格式化I/O

	"go.uber.org/zap/buffer"             // buffer包：高性能缓冲区
	"go.uber.org/zap/internal/bufferpool" // 内部缓冲池
	"go.uber.org/zap/internal/pool"       // 内部对象池
)

var _sliceEncoderPool = pool.New(func() *sliceArrayEncoder { // 切片编码器对象池
	return &sliceArrayEncoder{
		elems: make([]interface{}, 0, 2), // 预分配2个元素的切片
	}
})

func getSliceEncoder() *sliceArrayEncoder { // 从对象池获取切片编码器
	return _sliceEncoderPool.Get()          // 从池中获取
}

func putSliceEncoder(e *sliceArrayEncoder) { // 将切片编码器归还到对象池
	e.elems = e.elems[:0]                    // 重置切片长度为0
	_sliceEncoderPool.Put(e)                 // 归还到池中
}

type consoleEncoder struct { // 控制台编码器结构体
	*jsonEncoder             // 嵌入JSON编码器
}

// NewConsoleEncoder creates an encoder whose output is designed for human -
// rather than machine - consumption. It serializes the core log entry data
// (message, level, timestamp, etc.) in a plain-text format and leaves the
// structured context as JSON.
//
// Note that although the console encoder doesn't use the keys specified in the
// encoder configuration, it will omit any element whose key is set to the empty
// string.
// NewConsoleEncoder创建一个为人类而非机器消费设计的编码器。
// 它以纯文本格式序列化核心日志条目数据（消息、级别、时间戳等），
// 并将结构化上下文保留为JSON。
//
// 注意，虽然控制台编码器不使用编码器配置中指定的键，
// 但它会忽略任何键设置为空字符串的元素。
func NewConsoleEncoder(cfg EncoderConfig) Encoder {
	if cfg.ConsoleSeparator == "" {           // 如果控制台分隔符为空
		// Use a default delimiter of '\t' for backwards compatibility
		// 为了向后兼容，使用默认的'\t'分隔符
		cfg.ConsoleSeparator = "\t"           // 设置为制表符
	}
	return consoleEncoder{newJSONEncoder(cfg, true)} // 返回控制台编码器，使用带空格的JSON编码器
}

func (c consoleEncoder) Clone() Encoder { // 克隆控制台编码器
	return consoleEncoder{c.jsonEncoder.Clone().(*jsonEncoder)} // 克隆底层JSON编码器
}

func (c consoleEncoder) EncodeEntry(ent Entry, fields []Field) (*buffer.Buffer, error) {
	line := bufferpool.Get()

	// We don't want the entry's metadata to be quoted and escaped (if it's
	// encoded as strings), which means that we can't use the JSON encoder. The
	// simplest option is to use the memory encoder and fmt.Fprint.
	//
	// If this ever becomes a performance bottleneck, we can implement
	// ArrayEncoder for our plain-text format.
	arr := getSliceEncoder()
	if c.TimeKey != "" && c.EncodeTime != nil && !ent.Time.IsZero() {
		c.EncodeTime(ent.Time, arr)
	}
	if c.LevelKey != "" && c.EncodeLevel != nil {
		c.EncodeLevel(ent.Level, arr)
	}
	if ent.LoggerName != "" && c.NameKey != "" {
		nameEncoder := c.EncodeName

		if nameEncoder == nil {
			// Fall back to FullNameEncoder for backward compatibility.
			nameEncoder = FullNameEncoder
		}

		nameEncoder(ent.LoggerName, arr)
	}
	if ent.Caller.Defined {
		if c.CallerKey != "" && c.EncodeCaller != nil {
			c.EncodeCaller(ent.Caller, arr)
		}
		if c.FunctionKey != "" {
			arr.AppendString(ent.Caller.Function)
		}
	}
	for i := range arr.elems {
		if i > 0 {
			line.AppendString(c.ConsoleSeparator)
		}
		_, _ = fmt.Fprint(line, arr.elems[i])
	}
	putSliceEncoder(arr)

	// Add the message itself.
	if c.MessageKey != "" {
		c.addSeparatorIfNecessary(line)
		line.AppendString(ent.Message)
	}

	// Add any structured context.
	c.writeContext(line, fields)

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	line.AppendString(c.LineEnding)
	return line, nil
}

func (c consoleEncoder) writeContext(line *buffer.Buffer, extra []Field) {
	context := c.jsonEncoder.Clone().(*jsonEncoder)
	defer func() {
		// putJSONEncoder assumes the buffer is still used, but we write out the buffer so
		// we can free it.
		context.buf.Free()
		putJSONEncoder(context)
	}()

	addFields(context, extra)
	context.closeOpenNamespaces()
	if context.buf.Len() == 0 {
		return
	}

	c.addSeparatorIfNecessary(line)
	line.AppendByte('{')
	line.Write(context.buf.Bytes())
	line.AppendByte('}')
}

func (c consoleEncoder) addSeparatorIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendString(c.ConsoleSeparator)
	}
}
