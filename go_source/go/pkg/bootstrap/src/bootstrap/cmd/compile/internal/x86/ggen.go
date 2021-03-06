// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/compile/internal/x86/ggen.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/compile/internal/x86/ggen.go:1
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"bootstrap/cmd/compile/internal/ir"
	"bootstrap/cmd/compile/internal/objw"
	"bootstrap/cmd/compile/internal/types"
	"bootstrap/cmd/internal/obj"
	"bootstrap/cmd/internal/obj/x86"
)

func zerorange(pp *objw.Progs, p *obj.Prog, off, cnt int64, ax *uint32) *obj.Prog {
	if cnt == 0 {
		return p
	}
	if *ax == 0 {
		p = pp.Append(p, x86.AMOVL, obj.TYPE_CONST, 0, 0, obj.TYPE_REG, x86.REG_AX, 0)
		*ax = 1
	}

	if cnt <= int64(4*types.RegSize) {
		for i := int64(0); i < cnt; i += int64(types.RegSize) {
			p = pp.Append(p, x86.AMOVL, obj.TYPE_REG, x86.REG_AX, 0, obj.TYPE_MEM, x86.REG_SP, off+i)
		}
	} else if cnt <= int64(128*types.RegSize) {
		p = pp.Append(p, x86.ALEAL, obj.TYPE_MEM, x86.REG_SP, off, obj.TYPE_REG, x86.REG_DI, 0)
		p = pp.Append(p, obj.ADUFFZERO, obj.TYPE_NONE, 0, 0, obj.TYPE_ADDR, 0, 1*(128-cnt/int64(types.RegSize)))
		p.To.Sym = ir.Syms.Duffzero
	} else {
		p = pp.Append(p, x86.AMOVL, obj.TYPE_CONST, 0, cnt/int64(types.RegSize), obj.TYPE_REG, x86.REG_CX, 0)
		p = pp.Append(p, x86.ALEAL, obj.TYPE_MEM, x86.REG_SP, off, obj.TYPE_REG, x86.REG_DI, 0)
		p = pp.Append(p, x86.AREP, obj.TYPE_NONE, 0, 0, obj.TYPE_NONE, 0, 0)
		p = pp.Append(p, x86.ASTOSL, obj.TYPE_NONE, 0, 0, obj.TYPE_NONE, 0, 0)
	}

	return p
}

func ginsnop(pp *objw.Progs) *obj.Prog {
	// See comment in ../amd64/ggen.go.
	p := pp.Prog(x86.AXCHGL)
	p.From.Type = obj.TYPE_REG
	p.From.Reg = x86.REG_AX
	p.To.Type = obj.TYPE_REG
	p.To.Reg = x86.REG_AX
	return p
}
