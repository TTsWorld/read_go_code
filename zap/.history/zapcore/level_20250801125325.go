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
	"bytes"  // bytes包：字节操作
	"errors" // errors包：错误处理
	"fmt"    // fmt包：格式化I/O
)

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level") // 无法反序列化nil *Level的错误

// A Level is a logging priority. Higher levels are more important.
// Level是日志优先级。更高的级别更重要。
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	// DebugLevel日志通常很详细，在生产环境中通常被禁用。
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	// InfoLevel是默认的日志优先级。
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	// WarnLevel日志比Info更重要，但不需要个别的人工审查。
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	// ErrorLevel日志是高优先级的。如果应用程序运行顺利，
	// 它不应该生成任何错误级别的日志。
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	// DPanicLevel日志是特别重要的错误。在开发模式下，
	// logger在写入消息后会panic。
	DPanicLevel
	// PanicLevel logs a message, then panics.
	// PanicLevel记录一条消息，然后panic。
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	// FatalLevel记录一条消息，然后调用os.Exit(1)。
	FatalLevel

	_minLevel = DebugLevel // 最小级别
	_maxLevel = FatalLevel // 最大级别

	// InvalidLevel is an invalid value for Level.
	//
	// Core implementations may panic if they see messages of this level.
	// InvalidLevel是Level的无效值。
	//
	// Core实现在看到此级别的消息时可能会panic。
	InvalidLevel = _maxLevel + 1
)

// ParseLevel parses a level based on the lower-case or all-caps ASCII
// representation of the log level. If the provided ASCII representation is
// invalid an error is returned.
//
// This is particularly useful when dealing with text input to configure log
// levels.
func ParseLevel(text string) (Level, error) {
	var level Level
	err := level.UnmarshalText([]byte(text))
	return level, err
}

type leveledEnabler interface {
	LevelEnabler

	Level() Level
}

// LevelOf reports the minimum enabled log level for the given LevelEnabler
// from Zap's supported log levels, or [InvalidLevel] if none of them are
// enabled.
//
// A LevelEnabler may implement a 'Level() Level' method to override the
// behavior of this function.
//
//	func (c *core) Level() Level {
//		return c.currentLevel
//	}
//
// It is recommended that [Core] implementations that wrap other cores use
// LevelOf to retrieve the level of the wrapped core. For example,
//
//	func (c *coreWrapper) Level() Level {
//		return zapcore.LevelOf(c.wrappedCore)
//	}
func LevelOf(enab LevelEnabler) Level {
	if lvler, ok := enab.(leveledEnabler); ok {
		return lvler.Level()
	}

	for lvl := _minLevel; lvl <= _maxLevel; lvl++ {
		if enab.Enabled(lvl) {
			return lvl
		}
	}

	return InvalidLevel
}

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DPanicLevel:
		return "dpanic"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// CapitalString returns an all-caps ASCII representation of the log level.
func (l Level) CapitalString() string {
	// Printing levels in all-caps is common enough that we should export this
	// functionality.
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DPanicLevel:
		return "DPANIC"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// MarshalText marshals the Level to text. Note that the text representation
// drops the -Level suffix (see example).
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText unmarshals text to a level. Like MarshalText, UnmarshalText
// expects the text representation of a Level to drop the -Level suffix (see
// example).
//
// In particular, this makes it easy to configure logging levels using YAML,
// TOML, or JSON files.
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug":
		*l = DebugLevel
	case "info", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "warning":
		*l = WarnLevel
	case "error":
		*l = ErrorLevel
	case "dpanic":
		*l = DPanicLevel
	case "panic":
		*l = PanicLevel
	case "fatal":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

// Set sets the level for the flag.Value interface.
func (l *Level) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

// Get gets the level for the flag.Getter interface.
func (l *Level) Get() interface{} {
	return *l
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// LevelEnabler decides whether a given logging level is enabled when logging a
// message.
//
// Enablers are intended to be used to implement deterministic filters;
// concerns like sampling are better implemented as a Core.
//
// Each concrete Level value implements a static LevelEnabler which returns
// true for itself and all higher logging levels. For example WarnLevel.Enabled()
// will return true for WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, and
// FatalLevel, but return false for InfoLevel and DebugLevel.
type LevelEnabler interface {
	Enabled(Level) bool
}
