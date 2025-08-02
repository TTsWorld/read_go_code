// Copyright (c) 2019 Uber Technologies, Inc.
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

//go:build tools
// +build tools
// 构建约束：仅在tools构建标签下编译，正常构建时不会包含此文件

// Package tools contains import references to tools used in the build process.
// This package will not be included in the final binary since it uses the "tools" build tag.
// However, it ensures that the imported tools are tracked in go.mod and available via `go install`.
//
// This is a common Go pattern for managing development tools as dependencies.
// Package tools包含构建过程中使用的工具的导入引用。
// 由于使用了"tools"构建标签，此包不会包含在最终的二进制文件中。
// 但是，它确保导入的工具在go.mod中被跟踪，并可通过`go install`获得。
//
// 这是一个常见的Go模式，用于将开发工具作为依赖项管理。
package tools

import (
	// Tools we use during development.
	// 开发过程中使用的工具。
	
	// govulncheck: Go vulnerability checker - scans Go code for known security vulnerabilities
	// govulncheck: Go漏洞检查器 - 扫描Go代码中的已知安全漏洞
	_ "golang.org/x/vuln/cmd/govulncheck"
)
