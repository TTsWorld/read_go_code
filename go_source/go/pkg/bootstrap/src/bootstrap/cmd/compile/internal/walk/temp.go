// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/compile/internal/walk/temp.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/compile/internal/walk/temp.go:1
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bootstrap/cmd/compile/internal/base"
	"bootstrap/cmd/compile/internal/ir"
	"bootstrap/cmd/compile/internal/typecheck"
	"bootstrap/cmd/compile/internal/types"
)

// initStackTemp appends statements to init to initialize the given
// temporary variable to val, and then returns the expression &tmp.
func initStackTemp(init *ir.Nodes, tmp *ir.Name, val ir.Node) *ir.AddrExpr {
	if val != nil && !types.Identical(tmp.Type(), val.Type()) {
		base.Fatalf("bad initial value for %L: %L", tmp, val)
	}
	appendWalkStmt(init, ir.NewAssignStmt(base.Pos, tmp, val))
	return typecheck.Expr(typecheck.NodAddr(tmp)).(*ir.AddrExpr)
}

// stackTempAddr returns the expression &tmp, where tmp is a newly
// allocated temporary variable of the given type. Statements to
// zero-initialize tmp are appended to init.
func stackTempAddr(init *ir.Nodes, typ *types.Type) *ir.AddrExpr {
	return initStackTemp(init, typecheck.Temp(typ), nil)
}

// stackBufAddr returns thte expression &tmp, where tmp is a newly
// allocated temporary variable of type [len]elem. This variable is
// initialized, and elem must not contain pointers.
func stackBufAddr(len int64, elem *types.Type) *ir.AddrExpr {
	if elem.HasPointers() {
		base.FatalfAt(base.Pos, "%v has pointers", elem)
	}
	tmp := typecheck.Temp(types.NewArray(elem, len))
	return typecheck.Expr(typecheck.NodAddr(tmp)).(*ir.AddrExpr)
}
