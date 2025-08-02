// Copyright (c) 2020 Uber Technologies, Inc.
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

import "fmt" // fmt包：格式化I/O

type levelFilterCore struct { // 级别过滤Core
	core  Core         // 底层Core
	level LevelEnabler // 级别启用器
}

var (
	_ Core           = (*levelFilterCore)(nil) // 确保levelFilterCore实现Core接口
	_ leveledEnabler = (*levelFilterCore)(nil) // 确保levelFilterCore实现leveledEnabler接口
)

// NewIncreaseLevelCore creates a core that can be used to increase the level of
// an existing Core. It cannot be used to decrease the logging level, as it acts
// as a filter before calling the underlying core. If level decreases the log level,
// an error is returned.
// NewIncreaseLevelCore创建一个可用于提高现有Core级别的core。
// 它不能用于降低日志级别，因为它在调用底层core之前充当过滤器。
// 如果级别降低了日志级别，将返回错误。
func NewIncreaseLevelCore(core Core, level LevelEnabler) (Core, error) {
	for l := _maxLevel; l >= _minLevel; l-- { // 从最高级别到最低级别遍历
		if !core.Enabled(l) && level.Enabled(l) { // 如果原core不启用但新级别启用
			return nil, fmt.Errorf("invalid increase level, as level %q is allowed by increased level, but not by existing core", l) // 返回错误
		}
	}

	return &levelFilterCore{core, level}, nil // 返回级别过滤Core
}

func (c *levelFilterCore) Enabled(lvl Level) bool {
	return c.level.Enabled(lvl)
}

func (c *levelFilterCore) Level() Level {
	return LevelOf(c.level)
}

func (c *levelFilterCore) With(fields []Field) Core {
	return &levelFilterCore{c.core.With(fields), c.level}
}

func (c *levelFilterCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	if !c.Enabled(ent.Level) {
		return ce
	}

	return c.core.Check(ent, ce)
}

func (c *levelFilterCore) Write(ent Entry, fields []Field) error {
	return c.core.Write(ent, fields)
}

func (c *levelFilterCore) Sync() error {
	return c.core.Sync()
}
