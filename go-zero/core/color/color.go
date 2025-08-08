// 文件功能: 定义颜色枚举与彩色字符串渲染工具，基于第三方 `github.com/fatih/color` 实现。
// 关键技术点:
// - 颜色枚举: 使用自定义 `Color` 类型与常量枚举前景/背景颜色。
// - 第三方库: 借助 `fatih/color` 完成跨平台终端彩色输出。
// - 文本包装: 提供带填充的彩色输出，便于在日志中突出显示。
// 适用场景: 终端日志/提示信息的彩色渲染，增强可读性。
// 注意: 若终端或输出环境不支持 ANSI 颜色，渲染效果可能受限。
package color

// 导入第三方颜色库，用于终端彩色输出
import "github.com/fatih/color"

// 常量块: 定义前景/背景颜色的枚举值
const (
	// NoColor is no color for both foreground and background.
	NoColor Color = iota
	// FgBlack is the foreground color black.
	FgBlack
	// FgRed is the foreground color red.
	FgRed
	// FgGreen is the foreground color green.
	FgGreen
	// FgYellow is the foreground color yellow.
	FgYellow
	// FgBlue is the foreground color blue.
	FgBlue
	// FgMagenta is the foreground color magenta.
	FgMagenta
	// FgCyan is the foreground color cyan.
	FgCyan
	// FgWhite is the foreground color white.
	FgWhite

	// BgBlack is the background color black.
	BgBlack
	// BgRed is the background color red.
	BgRed
	// BgGreen is the background color green.
	BgGreen
	// BgYellow is the background color yellow.
	BgYellow
	// BgBlue is the background color blue.
	BgBlue
	// BgMagenta is the background color magenta.
	BgMagenta
	// BgCyan is the background color cyan.
	BgCyan
	// BgWhite is the background color white.
	BgWhite
)

// 颜色到 fatih/color 属性的映射表
var colors = map[Color][]color.Attribute{
	FgBlack:   {color.FgBlack, color.Bold},
	FgRed:     {color.FgRed, color.Bold},
	FgGreen:   {color.FgGreen, color.Bold},
	FgYellow:  {color.FgYellow, color.Bold},
	FgBlue:    {color.FgBlue, color.Bold},
	FgMagenta: {color.FgMagenta, color.Bold},
	FgCyan:    {color.FgCyan, color.Bold},
	FgWhite:   {color.FgWhite, color.Bold},
	BgBlack:   {color.BgBlack, color.FgHiWhite, color.Bold},
	BgRed:     {color.BgRed, color.FgHiWhite, color.Bold},
	BgGreen:   {color.BgGreen, color.FgHiWhite, color.Bold},
	BgYellow:  {color.BgHiYellow, color.FgHiBlack, color.Bold},
	BgBlue:    {color.BgBlue, color.FgHiWhite, color.Bold},
	BgMagenta: {color.BgMagenta, color.FgHiWhite, color.Bold},
	BgCyan:    {color.BgCyan, color.FgHiWhite, color.Bold},
	BgWhite:   {color.BgHiWhite, color.FgHiBlack, color.Bold},
}

// 自定义颜色类型，使用无符号整型承载枚举
type Color uint32

// WithColor returns a string with the given color applied.
// 功能: 将文本渲染为指定颜色的字符串
func WithColor(text string, colour Color) string {
	// 创建一个带有指定颜色属性的渲染器
	c := color.New(colors[colour]...)
	// 返回渲染后的字符串
	return c.Sprint(text)
}

// WithColorPadding returns a string with the given color applied with leading and trailing spaces.
// 功能: 在渲染文本的同时为其添加首尾空格，提升在日志中的可辨识度
func WithColorPadding(text string, colour Color) string {
	// 先在文本两侧添加空格，再调用 WithColor 进行渲染
	return WithColor(" "+text+" ", colour)
}
