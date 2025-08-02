// Copyright (c) 2016-2022 Uber Technologies, Inc.
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
	"errors"       // errors包：错误处理
	"fmt"          // fmt包：格式化I/O
	"io"           // io包：基本I/O原语
	"net/url"      // url包：URL解析
	"os"           // os包：操作系统接口
	"path/filepath" // filepath包：文件路径操作
	"strings"      // strings包：字符串操作
	"sync"         // sync包：同步原语

	"go.uber.org/zap/zapcore" // zapcore包：核心接口和实现
)

const schemeFile = "file" // 文件URL方案常量

var _sinkRegistry = newSinkRegistry() // 全局sink注册表

// Sink defines the interface to write to and close logger destinations.
// Sink定义了写入和关闭日志目标的接口。
type Sink interface {
	zapcore.WriteSyncer // 嵌入WriteSyncer接口
	io.Closer           // 嵌入Closer接口
}

type errSinkNotFound struct { // sink未找到错误类型
	scheme string // URL方案
}

func (e *errSinkNotFound) Error() string { // 实现error接口
	return fmt.Sprintf("no sink found for scheme %q", e.scheme) // 格式化错误消息
}

type nopCloserSink struct{ zapcore.WriteSyncer } // 无操作关闭器sink

func (nopCloserSink) Close() error { return nil } // Close方法：无操作，直接返回nil

type sinkRegistry struct { // sink注册表结构体
	mu        sync.Mutex                                      // 互斥锁保护并发访问
	factories map[string]func(*url.URL) (Sink, error)          // 工厂函数映射，按URL方案索引
	openFile  func(string, int, os.FileMode) (*os.File, error) // 文件打开函数，类型匹配os.OpenFile
}

func newSinkRegistry() *sinkRegistry { // 创建新的sink注册表
	sr := &sinkRegistry{
		factories: make(map[string]func(*url.URL) (Sink, error)), // 初始化工厂函数映射
		openFile:  os.OpenFile,                                   // 设置默认文件打开函数
	}
	// Infallible operation: the registry is empty, so we can't have a conflict.
	// 无风险操作：注册表为空，不可能有冲突。
	_ = sr.RegisterSink(schemeFile, sr.newFileSinkFromURL) // 注册默认文件方案
	return sr                                             // 返回注册表
}

// RegisterSink registers the given factory for the specific scheme.
// RegisterSink为特定方案注册给定的工厂函数。
func (sr *sinkRegistry) RegisterSink(scheme string, factory func(*url.URL) (Sink, error)) error {
	sr.mu.Lock()         // 获取写锁
	defer sr.mu.Unlock() // 延迟释放锁

	if scheme == "" { // 检查方案是否为空
		return errors.New("can't register a sink factory for empty string") // 返回错误
	}
	normalized, err := normalizeScheme(scheme) // 规范化方案名称
	if err != nil {                            // 如果规范化失败
		return fmt.Errorf("%q is not a valid scheme: %v", scheme, err) // 返回格式化错误
	}
	if _, ok := sr.factories[normalized]; ok { // 检查是否已注册
		return fmt.Errorf("sink factory already registered for scheme %q", normalized) // 返回重复注册错误
	}
	sr.factories[normalized] = factory // 注册工厂函数
	return nil                        // 成功返回nil
}

func (sr *sinkRegistry) newSink(rawURL string) (Sink, error) { // 创建新的sink
	// URL parsing doesn't work well for Windows paths such as `c:\log.txt`, as scheme is set to
	// the drive, and path is unset unless `c:/log.txt` is used.
	// To avoid Windows-specific URL handling, we instead check IsAbs to open as a file.
	// filepath.IsAbs is OS-specific, so IsAbs('c:/log.txt') is false outside of Windows.
	// URL解析对Windows路径（如`c:\log.txt`）支持不好，方案会被设置为驱动器名，
	// 路径为空（除非使用`c:/log.txt`）。为避免Windows特定的URL处理，
	// 我们检查IsAbs来作为文件打开。filepath.IsAbs是OS特定的，
	// 所以在Windows外IsAbs('c:/log.txt')为false。
	if filepath.IsAbs(rawURL) { // 如果是绝对路径
		return sr.newFileSinkFromPath(rawURL) // 直接作为文件路径处理
	}

	u, err := url.Parse(rawURL) // 解析URL
	if err != nil {             // 如果解析失败
		return nil, fmt.Errorf("can't parse %q as a URL: %v", rawURL, err) // 返回错误
	}
	if u.Scheme == "" { // 如果没有方案
		u.Scheme = schemeFile // 默认为文件方案
	}

	sr.mu.Lock()                          // 获取读锁
	factory, ok := sr.factories[u.Scheme] // 查找工厂函数
	sr.mu.Unlock()                        // 释放锁
	if !ok {                              // 如果未找到工厂
		return nil, &errSinkNotFound{u.Scheme} // 返回未找到错误
	}
	return factory(u) // 调用工厂函数创建sink
}

// RegisterSink registers a user-supplied factory for all sinks with a
// particular scheme.
//
// All schemes must be ASCII, valid under section 0.1 of RFC 3986
// (https://tools.ietf.org/html/rfc3983#section-3.1), and must not already
// have a factory registered. Zap automatically registers a factory for the
// "file" scheme.
// RegisterSink为特定方案的所有sink注册用户提供的工厂函数。
//
// 所有方案必须是ASCII，符合RFC 3986第0.1节
// (https://tools.ietf.org/html/rfc3983#section-3.1)的规范，
// 且不能已经注册了工厂函数。Zap自动为"file"方案注册工厂。
func RegisterSink(scheme string, factory func(*url.URL) (Sink, error)) error {
	return _sinkRegistry.RegisterSink(scheme, factory) // 委托给全局注册表
}

func (sr *sinkRegistry) newFileSinkFromURL(u *url.URL) (Sink, error) { // 从URL创建文件sink
	if u.User != nil { // 如果包含用户信息
		return nil, fmt.Errorf("user and password not allowed with file URLs: got %v", u) // 返回错误
	}
	if u.Fragment != "" { // 如果包含片段
		return nil, fmt.Errorf("fragments not allowed with file URLs: got %v", u) // 返回错误
	}
	if u.RawQuery != "" { // 如果包含查询参数
		return nil, fmt.Errorf("query parameters not allowed with file URLs: got %v", u) // 返回错误
	}
	// Error messages are better if we check hostname and port separately.
	// 分别检查主机名和端口，错误消息会更好。
	if u.Port() != "" { // 如果指定了端口
		return nil, fmt.Errorf("ports not allowed with file URLs: got %v", u) // 返回错误
	}
	if hn := u.Hostname(); hn != "" && hn != "localhost" { // 如果主机名非空且非localhost
		return nil, fmt.Errorf("file URLs must leave host empty or use localhost: got %v", u) // 返回错误
	}

	return sr.newFileSinkFromPath(u.Path) // 使用路径创建文件sink
}

func (sr *sinkRegistry) newFileSinkFromPath(path string) (Sink, error) { // 从路径创建文件sink
	switch path { // 根据路径判断
	case "stdout": // 标准输出
		return nopCloserSink{os.Stdout}, nil // 返回标准输出sink
	case "stderr": // 标准错误
		return nopCloserSink{os.Stderr}, nil // 返回标准错误sink
	}
	return sr.openFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666) // 打开文件，写入、追加、创建模式，权限666
}

func normalizeScheme(s string) (string, error) { // 规范化URL方案
	// https://tools.ietf.org/html/rfc3986#section-3.1
	// 参考RFC 3986第3.1节
	s = strings.ToLower(s)                        // 转换为小写
	if first := s[0]; 'a' > first || 'z' < first { // 检查首字符是否为字母
		return "", errors.New("must start with a letter") // 必须以字母开头
	}
	for i := 1; i < len(s); i++ { // 遍历字节（不是rune）
		c := s[i]         // 获取当前字符
		switch {          // 检查字符类型
		case 'a' <= c && c <= 'z': // 小写字母
			continue      // 继续
		case '0' <= c && c <= '9': // 数字
			continue      // 继续
		case c == '.' || c == '+' || c == '-': // 特殊字符
			continue      // 继续
		}
		return "", fmt.Errorf("may not contain %q", c) // 不能包含其他字符
	}
	return s, nil // 返回规范化后的方案
}
