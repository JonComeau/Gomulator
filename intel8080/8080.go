package intel8080

import (
	"fmt"
	"os"
)

func b2i8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func b2i16(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

func i2b8(i uint8) bool {
	if i == 0 {
		return false
	}
	return true
}

func i2b16(i uint16) bool {
	if i == 0 {
		return false
	}
	return true
}

func wb(state *State8080, addr uint16, val uint8) {
	state.wb(addr, val)
}

func rb(state *State8080, addr uint16) uint8 {
	return state.rb(addr)
}

func unimplementedInstruction(state *State8080) {
	fmt.Print("Error: Unimplemented instruction\n")
	fmt.Printf("Machine Code: %x", state.memory[state.pc])
	os.Exit(1)
}

func Emulate8080Op(state *State8080) {
	// TODO: Write Disassembly
	// disassemble8080Op(state.memory, state.pc)

	switch state.memory[state.pc] {
	case 0x00: // NOP
		break

	// LXI
	case 0x01: // LXI  B,D16
		state.c = state.memory[state.pc+1]
		state.b = state.memory[state.pc+2]
		state.pc += 2
		break
	case 0x11: // LXI  D,D16
		state.d = state.memory[state.pc+1]
		state.e = state.memory[state.pc+2]
		state.pc += 2
		break
	case 0x21: // LXI  H,D16
		state.h = state.memory[state.pc+1]
		state.l = state.memory[state.pc+2]
		state.pc += 2
		break
	case 0x31: // LXI  SP,D16
		state.sp = uint16(state.memory[state.pc+1]<<8) |
			uint16(state.memory[state.pc+2]&0xff)
		state.pc += 2
		break

	// STAX
	case 0x02: // STAX B
		wb(state, state.getBC(), state.a)
		break
	case 0x12: // STAX D
		wb(state, state.getDE(), state.a)
		break
	case 0x32: // STA  word
		wb(state,
			uint16(state.memory[state.pc+1]<<8)|
				uint16(state.memory[state.pc+2]&0xff),
			state.a)
		break

	// INX
	case 0x03: // INX  B
		state.setBC(state.getBC() + 1)
		break
	case 0x13: // INX  D
		state.setDE(state.getDE() + 1)
		break
	case 0x23: // INX  H
		state.setHL(state.getHL() + 1)
		break
	case 0x33: // INX  SP
		state.sp++
		break

	// INR
	case 0x04: // INR  B
		inr(state, &state.b)
		break
	case 0x0c: // INR  C
		inr(state, &state.c)
		break
	case 0x14: // INR  D
		inr(state, &state.d)
		break
	case 0x1c: // INR  E
		inr(state, &state.e)
		break
	case 0x24: // INR  H
		inr(state, &state.h)
		break
	case 0x2c: // INR  L
		inr(state, &state.l)
		break
	case 0x34: // INR  M
		bite := rb(state, state.getHL())
		wb(state,
			state.getHL(),
			inr(state, &bite))
	case 0x3c: // INR  A
		inr(state, &state.a)
		break

	// DCR
	case 0x05: // DCR  B
		dcr(state, &state.b)
		break
	case 0x0d: // DCR  C
		dcr(state, &state.c)
		break
	case 0x15: // DCR  D
		dcr(state, &state.d)
		break
	case 0x1d: // DCR  E
		dcr(state, &state.e)
		break
	case 0x25: // DCR  H
		dcr(state, &state.h)
		break
	case 0x2d: // DCR  L
		dcr(state, &state.l)
		break
	case 0x35: // DCR  M
		bite := rb(state, state.getHL())
		wb(state, state.getHL(), inr(state, &bite))
	case 0x3d: // DCR  A
		dcr(state, &state.a)
		break

	// MVI
	case 0x06: // MVI  B,D8
		state.b = state.memory[state.pc+1]
		state.pc++
		break
	case 0x0e: // MVI  C,D8
		state.c = state.memory[state.pc+1]
		state.pc++
		break
	case 0x16: // MVI  D,D8
		state.d = state.memory[state.pc+1]
		state.pc++
		break
	case 0x1e: // MVI  E,D8
		state.b = state.memory[state.pc+1]
		state.pc++
		break
	case 0x26: // MVI  H,D8
		state.h = state.memory[state.pc+1]
		state.pc++
		break
	case 0x2e: // MVI  L,D8
		state.l = state.memory[state.pc+1]
		state.pc++
		break
	case 0x36: // MVI  M,D8
		wb(state, state.getHL(), state.memory[state.pc+1])
		state.pc++
		break
	case 0x3e: // MVI  A,D8
		state.a = state.memory[state.pc+1]
		state.pc++
		break

	// Rotation Functions
	case 0x07: // RLC
		state.cc.cy = state.a >> 7
		state.a = (state.a << 1) | (state.cc.GetCy() << 7)
		break
	case 0x0f: // RRC
		x := state.a
		state.a = ((x & 1) << 7) | (x >> 1)
		state.cc.cy = b2i8(1 == (x & 1))
		break
	case 0x17: // RAL
		cy := state.cc.GetCy()
		state.cc.cy = state.a >> 7
		state.a = (state.a << 1) | cy
	case 0x1f: // RAR
		x := state.a
		state.a = (state.cc.cy << 7) | (x >> 1)
		state.cc.cy = b2i8(1 == (x & 1))
		break

	// DAD
	case 0x09: // DAD  B
		dad(state, state.getBC())
		break
	case 0x19: // DAD  D
		dad(state, state.getDE())
		break
	case 0x29: // DAD  H
		dad(state, state.getHL())
		break
	case 0x39: // DAD  SP
		dad(state, state.sp)
		break

	// LDAX
	case 0x0a: // LDAX B
		state.a = rb(state, state.getBC())
		break
	case 0x1a: // LDAX D
		state.a = rb(state, state.getDE())
		break
	case 0x3a: // LDA  adr
		state.a = rb(state, uint16(state.memory[state.pc+1]<<8)|
			uint16(state.memory[state.pc+2]&0xff))
		break

	// DCX
	case 0x0b: // DCX  B
	case 0x1b: // DCX  D
	case 0x2b: // DCX  H
	case 0x3b: // DCX  SP

	// SHLD
	case 0x22: // SHLD adr

	// DAA
	case 0x27: // DAA

	// LHDH
	case 0x2a: // LHDH adr

	// CMA
	case 0x2f: // CMA (not)
		state.a = state.a ^ state.a
		break

	// STC
	case 0x37: // STC
		state.cc.cy = 1
		break

	// MOV
	case 0x40: // MOV  B,B
		state.b = state.b
		break
	case 0x41: // MOV  B,C
		state.b = state.c
		break
	case 0x42: // MOV  B,D
		state.b = state.d
		break
	case 0x43: // MOV  B,E
		state.b = state.e
		break
	case 0x44: // MOV  B,H
		state.b = state.h
		break
	case 0x45: // MOV  B,L
		state.b = state.l
		break
	case 0x46: // MOV  B,M
		state.b = rb(state, state.getHL())
		break
	case 0x47: // MOV  B,A
		state.b = state.a
		break
	case 0x48: // MOV  C,B
		state.c = state.b
		break
	case 0x49: // MOV  C,C
		state.c = state.c
		break
	case 0x4a: // MOV  C,D
		state.c = state.d
		break
	case 0x4b: // MOV  C,E
		state.c = state.e
		break
	case 0x4c: // MOV  C,H
		state.c = state.h
		break
	case 0x4d: // MOV  C,L
		state.c = state.l
		break
	case 0x4e: // MOV  C,M
		state.c = rb(state, state.getHL())
		break
	case 0x4f: // MOV  C,A
		state.c = state.a
		break
	case 0x50: // MOV  D,B
		state.d = state.b
		break
	case 0x51: // MOV  D,C
		state.d = state.c
		break
	case 0x52: // MOV  D,D
		state.d = state.d
		break
	case 0x53: // MOV  D,E
		state.d = state.e
		break
	case 0x54: // MOV  D,H
		state.d = state.h
		break
	case 0x55: // MOV  D,L
		state.d = state.l
		break
	case 0x56: // MOV  D,M
		state.d = rb(state, state.getHL())
		break
	case 0x57: // MOV  D,A
		state.d = state.a
		break
	case 0x58: // MOV  E,B
		state.e = state.b
		break
	case 0x59: // MOV  E,C
		state.e = state.c
		break
	case 0x5a: // MOV  E,D
		state.e = state.d
		break
	case 0x5b: // MOV  E,E
		state.e = state.e
		break
	case 0x5c: // MOV  E,H
		state.e = state.h
		break
	case 0x5d: // MOV  E,L
		state.e = state.l
		break
	case 0x5e: // MOV  E,M
		state.e = rb(state, state.getHL())
		break
	case 0x5f: // MOV  E,A
		state.e = state.a
		break
	case 0x60: // MOV  H,B
		state.h = state.b
		break
	case 0x61: // MOV  H,C
		state.h = state.c
		break
	case 0x62: // MOV  H,D
		state.h = state.d
		break
	case 0x63: // MOV  H,E
		state.h = state.e
		break
	case 0x64: // MOV  H,H
		state.h = state.h
		break
	case 0x65: // MOV  H,L
		state.h = state.l
		break
	case 0x66: // MOV  H,M
		state.h = rb(state, state.getHL())
		break
	case 0x67: // MOV  H,A
		state.h = state.a
		break
	case 0x68: // MOV  L,B
		state.l = state.b
		break
	case 0x69: // MOV  L,C
		state.l = state.c
		break
	case 0x6a: // MOV  L,D
		state.l = state.d
		break
	case 0x6b: // MOV  L,E
		state.l = state.e
		break
	case 0x6c: // MOV  L,H
		state.l = state.h
		break
	case 0x6d: // MOV  L,L
		state.l = state.l
		break
	case 0x6e: // MOV  L,M
		state.l = rb(state, state.getHL())
		break
	case 0x6f: // MOV  L,A
		state.l = state.a
		break
	case 0x70: // MOV  M,B
		wb(state, state.getHL(), state.b)
		break
	case 0x71: // MOV  M,C
		wb(state, state.getHL(), state.c)
		break
	case 0x72: // MOV  M,D
		wb(state, state.getHL(), state.d)
		break
	case 0x73: // MOV  M,E
		wb(state, state.getHL(), state.e)
		break
	case 0x74: // MOV  M,H
		wb(state, state.getHL(), state.h)
		break
	case 0x75: // MOV  M,L
		wb(state, state.getHL(), state.l)
		break
	case 0x77: // MOV  M,A
		wb(state, state.getHL(), state.a)
		break
	case 0x78: // MOV  A,B
		state.a = state.b
		break
	case 0x79: // MOV  A,C
		state.a = state.c
		break
	case 0x7a: // MOV  A,D
		state.a = state.d
		break
	case 0x7b: // MOV  A,E
		state.a = state.e
		break
	case 0x7c: // MOV  A,H
		state.a = state.h
		break
	case 0x7d: // MOV  A,L
		state.a = state.l
		break
	case 0x7e: // MOV  A,M
		state.a = rb(state, state.getHL())
		break
	case 0x7f: // MOV  A,A
		state.a = state.a
		break

	// HLT
	case 0x76: // HLT
		state.pc--
		break

	// ADD
	case 0x80: // ADD  B
		add(state, state.b, 0)
		break
	case 0x81: // ADD  C
		add(state, state.c, 0)
		break
	case 0x82: // ADD  D
		add(state, state.d, 0)
		break
	case 0x83: // ADD  E
		add(state, state.e, 0)
		break
	case 0x84: // ADD  H
		add(state, state.h, 0)
		break
	case 0x85: // ADD  L
		add(state, state.l, 0)
		break
	case 0x86: // ADD  M
		add(state, rb(state, state.getHL()), 0)
		break
	case 0x87: // ADD  A
		add(state, state.a, 0)
		break

	// ADC
	case 0x88: // ADC  B
		add(state, state.b, 1)
		break
	case 0x89: // ADC  C
		add(state, state.c, 1)
		break
	case 0x8a: // ADC  D
		add(state, state.d, 1)
		break
	case 0x8b: // ADC  E
		add(state, state.e, 1)
		break
	case 0x8c: // ADC  H
		add(state, state.h, 1)
		break
	case 0x8d: // ADC  L
		add(state, state.l, 1)
		break
	case 0x8e: // ADC  M
		add(state, rb(state, state.getHL()), 1)
		break
	case 0x8f: // ADC  A
		add(state, state.a, 1)
		break

	// SUB
	case 0x90: // SUB  B
		sub(state, state.b, 0)
	case 0x91: // SUB  C
		sub(state, state.c, 0)
	case 0x92: // SUB  D
		sub(state, state.d, 0)
	case 0x93: // SUB  E
		sub(state, state.e, 0)
	case 0x94: // SUB  H
		sub(state, state.h, 0)
	case 0x95: // SUB  L
		sub(state, state.l, 0)
	case 0x96: // SUB  M
		sub(state, rb(state, state.getHL()), 0)
	case 0x97: // SUB  A
		sub(state, state.a, 0)

	// SBB
	case 0x98: // SBB  B
	case 0x99: // SBB  C
	case 0x9a: // SBB  D
	case 0x9b: // SBB  E
	case 0x9c: // SBB  H
	case 0x9d: // SBB  L
	case 0x9e: // SBB  M
	case 0x9f: // SBB  A

	// ANA
	case 0xa0: // ANA  B
	case 0xa1: // ANA  C
	case 0xa2: // ANA  D
	case 0xa3: // ANA  E
	case 0xa4: // ANA  H
	case 0xa5: // ANA  L
	case 0xa6: // ANA  M
	case 0xa7: // ANA  A

	// XRA
	case 0xa8: // ANA  B
	case 0xa9: // ANA  C
	case 0xaa: // ANA  D
	case 0xab: // ANA  E
	case 0xac: // ANA  H
	case 0xad: // ANA  L
	case 0xae: // ANA  M
	case 0xaf: // ANA  A

	// ORA
	case 0xb0: // ORA  B
	case 0xb1: // ORA  C
	case 0xb2: // ORA  D
	case 0xb3: // ORA  E
	case 0xb4: // ORA  H
	case 0xb5: // ORA  L
	case 0xb6: // ORA  M
	case 0xb7: // ORA  A

	// CMP
	case 0xb8: // CMP  B
	case 0xb9: // CMP  C
	case 0xba: // CMP  D
	case 0xbb: // CMP  E
	case 0xbc: // CMP  H
	case 0xbd: // CMP  L
	case 0xbe: // CMP  M
	case 0xbf: // CMP  A

	// POP
	case 0xc1: // POP  B
		state.c = state.memory[state.sp]
		state.b = state.memory[state.sp+1]
		state.sp += 2
		break
	case 0xd1: // POP  D
	case 0xe1: // POP  H
	case 0xf1: // Pop  PSW
		state.a = state.memory[state.sp+1]
		psw := state.memory[state.sp]
		state.cc.z = b2i8(0x01 == (psw & 0x01))
		state.cc.s = b2i8(0x02 == (psw & 0x02))
		state.cc.p = b2i8(0x04 == (psw & 0x04))
		state.cc.cy = b2i8(0x08 == (psw & 0x08))
		state.cc.ac = b2i8(0x10 == (psw & 0x10))
		state.sp += 2
		break

	// Jump Commands
	case 0xc2: // JNZ
		if state.cc.z == 0 {
			state.pc = uint16(
				(state.memory[state.pc+2] << 8) | state.memory[state.pc+1])
		} else {
			state.pc += 2
		}
		break
	case 0xc3: // JMP
		state.pc = uint16(
			(state.memory[state.pc+2] << 8) | state.memory[state.pc+1])
		break
	case 0xca: // JZ
	case 0xd2: // JNC
	case 0xda: // JC
	case 0xe2: // JPO
	case 0xea: // JPE
	case 0xf2: // JP
	case 0xfa: // JM

	// PUSH
	case 0xc5: // PUSH B
		state.memory[state.sp-1] = state.b
		state.memory[state.sp-2] = state.c
		state.sp -= 2
		break
	case 0xd5: // PUSH D
	case 0xe5: // PUSH H
	case 0xf5: // Push PSW
		state.memory[state.sp-1] = state.a
		psw := state.cc.z |
			state.cc.s<<1 |
			state.cc.p<<2 |
			state.cc.cy<<3 |
			state.cc.ac<<4
		state.memory[state.sp-2] = psw
		state.sp -= 2
		break

	// ADI
	case 0xc6: // ADI  D8

	// RST
	case 0xc7: // RST  0
	case 0xcf: // RST  1
	case 0xd7: // RST  2
	case 0xdf: // RST  3
	case 0xe7: // RST  4
	case 0xef: // RST  5
	case 0xf7: // RST  6
	case 0xff: // RST  7

	// Return Functions
	case 0xc0: // RNZ
	case 0xc8: // RZ
	case 0xc9: // RET
		state.pc = uint16(
			state.memory[state.sp] | (state.memory[state.sp+1] << 8))
		state.sp += 2
		break
	case 0xd0: // RNC
	case 0xd8: // RC
	case 0xe0: // RPO
	case 0xe8: // RPE
	case 0xf0: // RP
	case 0xf8: // RM

	// Call Functions
	case 0xc4: // CNZ
	case 0xcc: // CZ
	case 0xcd: // Call Address
		ret := uint16(state.pc + 2)
		state.memory[state.sp-1] = uint8((ret >> 8) & 0xff)
		state.memory[state.sp-2] = uint8(ret & 0xff)
		state.sp = state.sp - 2
		state.pc = uint16(
			(state.memory[state.pc+2] << 8) | state.memory[state.pc+1])
		break
	case 0xd4: // CNC
	case 0xdc: // CC
	case 0xe4: // CPO
	case 0xec: // CPE
	case 0xf4: // CP
	case 0xfc: // CM

	case 0xd3: // OUT  D8
	case 0xe6: // ANI
		x := state.a & state.memory[state.pc+1]
		state.cc.z = b2i8(x == 0)
		state.cc.s = b2i8(0x80 == (x & 0x80))
		state.cc.p = parity(x, 8)
		state.cc.cy = 0
		state.a = x
		state.pc++
		break
	case 0xeb: // XCHG
	case 0xfb: // EI
	case 0xfe: // CPI
		x := state.a - state.memory[state.pc+1]
		state.cc.z = b2i8(x == 0)
		state.cc.s = b2i8(0x80 == (x & 0x80))
		state.cc.p = parity(x, 8)
		state.cc.cy = b2i8(state.a > state.memory[state.pc+1])
		state.pc++
		break
	default:
		unimplementedInstruction(state)
		break
	}

	fmt.Printf("\tC=%d, P=%d, S=%d, Z=%d",
		state.cc.cy, state.cc.p,
		state.cc.s, state.cc.z)
	fmt.Printf("\tA $%02x B $%02x C $%02x D $%02x E $%02x H $%02x L $%02x SP %04x\n",
		state.a, state.b, state.c, state.d,
		state.e, state.h, state.l, state.sp)

	state.pc++
}

// region Carry Bit Instructions

// Set carry
func stc(state *State8080) {
	state.cc.cy = 0x01
}

// Complement Carry
func cmc(state *State8080) {
	state.cc.cy = state.cc.cy ^ 0x01
}

//endregion

//region Single Register Instructions

// TODO: CMA
// Increment Register or Memory
func inr(state *State8080, val *uint8) uint8 {
	*val++

	// Set Flags
	state.cc.ac = b2i8((*val & 0x0f) == 0)
	state.cc.z = b2i8((*val) == 0)
	state.cc.s = b2i8((*val & 0x80) != 0)
	state.cc.p = parity(*val)

	return *val
}

// Decrement Register or Memory
func dcr(state *State8080, val *uint8) {
	*val--

	// Set Flags
	state.cc.ac = b2i8(!((*val & 0x0f) == 0x0f))
	state.cc.z = b2i8((*val) == 0)
	state.cc.s = b2i8((*val & 0x80) != 0)
	state.cc.p = parity(*val)
}

// Decimal Adjust Accumulator
// adjusts register A to form two 4bit binary coded decimal digits.
// example: we want to add 93 and 8 (decimal operation):
//     MOV A, 0x93
//     ADI 0x08
//     ; now, A = 0x9B (0b10011011)
//     DAA
//     ; now, A = 0x01 (because 93 + 8 = 101)
//     ; and carry flag is set
func daa(state *State8080) {
	var cy, val uint8 = state.cc.GetCy(), 0
	lsb, msb := state.a&0x0f, state.a>>4

	if state.cc.GetAc() == 1 || lsb > 9 {
		val += 0x06
	}
	if state.cc.GetAc() == 1 || msb > 9 ||
		(msb >= 9 && lsb > 9) {
		val += 0x60
		cy = 1
	}
	add(state, val)
	state.cc.p = parity(state.a)
	state.cc.cy = cy
}

//endregion

// region Data Transfer Instructions

// TODO: MOV
// Store Accumulator
func stax(state *State8080) {
	wb(state, state.getBC(), state.a)
}

// Load Accumulator
func ldax(state *State8080, addr uint16) uint8 {
	return rb(state, addr)
}

// endregion

// region Register/Memory to Accumulator Instructions

// Add reg/mem to Accumulator with optional carry
func add(state *State8080, toAdd uint8, cy uint8) {
	answer := uint16(state.a) + uint16(toAdd) + uint16(cy)
	state.cc.z = b2i8((answer & 0xff) == 0)
	state.cc.s = b2i8((answer & 0x80) != 0)
	state.cc.cy = b2i8(answer > 0xff)
	state.cc.p = parity(answer & 0xff)
	state.a = uint8(answer & 0xff)
}

// Sub reg/mem from Accumulator with optional carry
func sub(state *State8080, toSub uint8, cy uint8) {
	var res uint16 = uint16(state.a) - uint16(toSub) -
}

// Logical & reg/mem with Accumulator
func ana() {}

// Logical ^ reg/mem with Accumulator
func xra() {}

// Logical | rem/mem with Accumulator
func ora() {}

// Compare Reg/mem with Accumulator
func cmp() {}

// endregion

// region Rotate Accumulator Instructions

// Rotate Accumulator Left
func rlc() {}

// Rotate Accumulator Right
func rrc() {}

// Rotate Accumulator Left w/carry
func ral() {}

// Rotate Accumulator Right w/carry
func rar() {}

// endregion

// region Register Pair Instructions

// Double Add
func dad(state *State8080, val uint16) {
	var res uint32 = uint32(state.getHL()) + uint32(val)
	state.setHL(uint16(res & 0xffff))
	state.cc.cy = b2i8((res & 0x10000) != 0)
}

// Increment Reg Pair
func inx(state *State8080) {}

// Decrement Reg Pair
func dcx(state *State8080) {}

// Exchange Regs
func xchg(state *State8080) {}

// Exchange Stack
func xthl(state *State8080) {}

// Load SP From H and L
func sphl(state *State8080) {}

// endregion

// region Immediate Instructions

// Load reg pair Immediate
func lxi(state *State8080) {}

// Move Immediate Data
func mvi(state *State8080) {}

// Add Immediate to Accumulator
func adi(state *State8080) {}

// Add Immediate to Accumulator w/carry
func aci(state *State8080) {}

// Sub Immediate from Accumulator
func sui(state *State8080) {}

// Sub Immediate from Accumulator w/borrow
func sbi(state *State8080) {}

// & Immediate with Accumulator
func ani(state *State8080) {}

// ^ Immediate with Accumulator
func xri(state *State8080) {}

// | Immediate with Accumulator
func ori(state *State8080) {}

// Compare Immediate with Accumulator
func cpi(state *State8080) {}

// endregion

// region Direct Addressing Instructions

// Store Accumulator Direct
func sta(state *State8080) {}

// Load Accumulator Direct
func lda(state *State8080) {}

// Store H and L Direct
func shld(state *State8080) {}

// Load H and L Direct
func lhld(state *State8080) {}

// endregion

// region Jump Instructions

// Load Program Counter
func pchl(state *State8080) {}

// Jump
func jmp(state *State8080) {}

// Jump If Carry
func jc(state *State8080) {}

// Jump If No Carry
func jnc(state *State8080) {}

// Jump If Zero
func jz(state *State8080) {}

// Jump If Not Zero
func jnz(state *State8080) {}

// Jump If Minus
func jm(state *State8080) {}

// Jump If Positive
func jp(state *State8080) {}

// Jump if Parity Even
func jpe(state *State8080) {}

// Jump If Parity Odd
func jpo(state *State8080) {}

// endregion

// region Call Subroutine Instructions

// Call
func call(state *State8080) {}

// Call if Carry
func cc(state *State8080) {}

// Call If No Carry
func cnc(state *State8080) {}

// Call If Zero
func cz(state *State8080) {}

// Call If Not Zero
func cnz(state *State8080) {}

// Call If Minus
func cm(state *State8080) {}

// Call If Plus
func cp(state *State8080) {}

// Call If Parity Even
func cpe(state *State8080) {}

// Call If Parity Odd
func cpo(state *State8080) {}

// endregion

// region Return From Subroutine Instructions

// Return
func ret(state *State8080) {}

// Return If Carry
func rn(state *State8080) {}

// Return if No Carry
func rnc(state *State8080) {}

// Return If Zero
func rz(state *State8080) {}

// Return If Not Zero
func rnz(state *State8080) {}

// Return If Minus
func rm(state *State8080) {}

// Return If Plus
func rp(state *State8080) {}

// Return If Parity Even
func rpe(state *State8080) {}

// Return If Parity Odd
func rpo(state *State8080) {}

// endregion

// RST Instruction
func rst(state *State8080) {}

// region Interrupt Flip-Flop Instructions

// Enable Interrupts
func ei(state *State8080) {}

// Disable Interrupts
func di(state *State8080) {}

// endregion

// region Input/Output Instructions

// Input
func in(state *State8080) {}

// Output
func out(state *State8080) {}

// endregion

// Halt Instruction
func hlt(state *State8080) {}

// region Pseudo-Instructions

// Origin
func org(state *State8080) {}

// Equate
func equ(state *State8080) {}

// Set
func set(state *State8080) {}

// End of Assembly
func end(state *State8080) {}

// Conditional Assembly

// Macro Definition

// endregion

func adm(state *State8080) {
	offset := uint16((state.h << 8) | state.l)
	answer := uint16(state.a + state.memory[offset])
	state.cc.z = b2i8((answer & 0xff) == 0)
	state.cc.s = b2i8((answer & 0x80) != 0)
	state.cc.cy = b2i8(answer > 0xff)
	state.cc.p = parity(answer & 0xff)
	state.a = uint8(answer & 0xff)
}

func disassemble8080Op(mem []uint, pc uint8) {

}
