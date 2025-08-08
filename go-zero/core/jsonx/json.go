// 文件功能: 提供 JSON 编解码辅助函数，优化标准库默认行为（禁止 HTML 转义、去除 Encode 结尾换行），并增强错误上下文。
// 关键技术点:
// - 自定义 Marshal: 使用 `json.Encoder` 禁用 HTML 转义，避免 API 响应中把 <、>、& 转义。
// - 去除末尾换行: `Encoder.Encode` 会在末尾加换行，手动移除保持一致性。
// - UseNumber: `Decoder.UseNumber()` 避免大整数精度丢失，延迟解析为 `json.Number`。
// - 错误包装: 返回包含原始字符串内容的错误，便于问题定位。
// 适用场景: Web API 响应、配置加载、跨语言数据交换等。
package jsonx

// 导入依赖包列表
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Marshal marshals v into json bytes, without escaping HTML and removes the trailing newline.
func Marshal(v any) ([]byte, error) {
	// why not use json.Marshal?  https://github.com/golang/go/issues/28453
	// it changes the behavior of json.Marshal, like & -> \u0026, < -> \u003c, > -> \u003e
	// which is not what we want in API responses
	// 使用 Encoder 并禁用 HTML 转义，保证输出与期望一致
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	bs := buf.Bytes()
	// Remove trailing newline added by json.Encoder.Encode
	// `Encode` 默认在末尾添加换行符，这里主动移除
	if len(bs) > 0 && bs[len(bs)-1] == '\n' {
		bs = bs[:len(bs)-1]
	}

	return bs, nil
}

// MarshalToString marshals v into a string.
func MarshalToString(v any) (string, error) {
	// 复用上面的 Marshal，再转换为字符串
	data, err := Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Unmarshal unmarshals data bytes into v.
func Unmarshal(data []byte, v any) error {
	// 构造 Decoder 并启用 UseNumber 避免精度丢失
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		// 格式化错误信息时附带原始内容，便于定位
		return formatError(string(data), err)
	}

	return nil
}

// UnmarshalFromString unmarshals v from str.
func UnmarshalFromString(str string, v any) error {
	decoder := json.NewDecoder(strings.NewReader(str))
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return formatError(str, err)
	}

	return nil
}

// UnmarshalFromReader unmarshals v from reader.
func UnmarshalFromReader(reader io.Reader, v any) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := json.NewDecoder(teeReader)
	if err := unmarshalUseNumber(decoder, v); err != nil {
		return formatError(buf.String(), err)
	}

	return nil
}

func unmarshalUseNumber(decoder *json.Decoder, v any) error {
	// 使用 UseNumber 延迟解析数字类型，防止 float64 导致精度丢失
	decoder.UseNumber()
	return decoder.Decode(v)
}

func formatError(v string, err error) error {
	// 将原始字符串与底层错误一起返回，包含上下文以提升可观测性
	return fmt.Errorf("string: `%s`, error: `%w`", v, err)
}
