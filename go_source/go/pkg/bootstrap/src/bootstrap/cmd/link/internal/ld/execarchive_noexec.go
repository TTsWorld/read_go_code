// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/link/internal/ld/execarchive_noexec.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/link/internal/ld/execarchive_noexec.go:1
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm || windows
// +build wasm windows

package ld

const syscallExecSupported = false

func (ctxt *Link) execArchive(argv []string) {
	panic("should never arrive here")
}
