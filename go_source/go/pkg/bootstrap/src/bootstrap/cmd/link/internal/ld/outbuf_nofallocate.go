// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/link/internal/ld/outbuf_nofallocate.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/link/internal/ld/outbuf_nofallocate.go:1
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!darwin && !linux) || (darwin && !go1.12)
// +build !darwin,!linux darwin,!go1.12

package ld

func (out *OutBuf) fallocate(size uint64) error {
	return errNoFallocate
}
