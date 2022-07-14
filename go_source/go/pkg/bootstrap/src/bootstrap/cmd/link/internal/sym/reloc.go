// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/link/internal/sym/reloc.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/link/internal/sym/reloc.go:1
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sym

import (
	"bootstrap/cmd/internal/objabi"
	"bootstrap/cmd/internal/sys"
	"bootstrap/debug/elf"
)

// RelocVariant is a linker-internal variation on a relocation.
type RelocVariant uint8

const (
	RV_NONE RelocVariant = iota
	RV_POWER_LO
	RV_POWER_HI
	RV_POWER_HA
	RV_POWER_DS

	// RV_390_DBL is a s390x-specific relocation variant that indicates that
	// the value to be placed into the relocatable field should first be
	// divided by 2.
	RV_390_DBL

	RV_CHECK_OVERFLOW RelocVariant = 1 << 7
	RV_TYPE_MASK      RelocVariant = RV_CHECK_OVERFLOW - 1
)

func RelocName(arch *sys.Arch, r objabi.RelocType) string {
	// We didn't have some relocation types at Go1.4.
	// Uncomment code when we include those in bootstrap code.

	switch {
	case r >= objabi.MachoRelocOffset: // Mach-O
		// nr := (r - objabi.MachoRelocOffset)>>1
		// switch ctxt.Arch.Family {
		// case sys.AMD64:
		// 	return macho.RelocTypeX86_64(nr).String()
		// case sys.ARM:
		// 	return macho.RelocTypeARM(nr).String()
		// case sys.ARM64:
		// 	return macho.RelocTypeARM64(nr).String()
		// case sys.I386:
		// 	return macho.RelocTypeGeneric(nr).String()
		// default:
		// 	panic("unreachable")
		// }
	case r >= objabi.ElfRelocOffset: // ELF
		nr := r - objabi.ElfRelocOffset
		switch arch.Family {
		case sys.AMD64:
			return elf.R_X86_64(nr).String()
		case sys.ARM:
			return elf.R_ARM(nr).String()
		case sys.ARM64:
			return elf.R_AARCH64(nr).String()
		case sys.I386:
			return elf.R_386(nr).String()
		case sys.MIPS, sys.MIPS64:
			return elf.R_MIPS(nr).String()
		case sys.PPC64:
			return elf.R_PPC64(nr).String()
		case sys.S390X:
			return elf.R_390(nr).String()
		default:
			panic("unreachable")
		}
	}

	return r.String()
}