// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/math/big/arith_decl.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/math/big/arith_decl.go:1
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !math_big_pure_go
// +build !math_big_pure_go

package big

// implemented in arith_$GOARCH.s
func mulWW(x, y Word) (z1, z0 Word)
func addVV(z, x, y []Word) (c Word)
func subVV(z, x, y []Word) (c Word)
func addVW(z, x []Word, y Word) (c Word)
func subVW(z, x []Word, y Word) (c Word)
func shlVU(z, x []Word, s uint) (c Word)
func shrVU(z, x []Word, s uint) (c Word)
func mulAddVWW(z, x []Word, y, r Word) (c Word)
func addMulVVW(z, x []Word, y Word) (c Word)
