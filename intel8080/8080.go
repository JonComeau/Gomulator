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

func unimplementedInstruction(state *State8080) {
	fmt.Print("Error: Unimplemented instruction\n")
	os.Exit(1)
}

func Emulate8080Op(state *State8080) {
	disassemble8080Op(state.memory, state.pc)

	switch state.memory[state.pc] {
	case 0x00: // NOP
		break
	case 0x01: // LXI  B,D16
		state.c = state.memory[state.pc+1]
		state.b = state.memory[state.pc+2]
		state.pc += 2
		break
	case 0x05: // DCR  B
	case 0x06: // MVI  B,D8
	case 0x09: // DAD  B
	case 0x0d: // DCR  C
	case 0x0e: // MVI  C,D8
	case 0x0f: // RRC
		x := state.a
		state.a = ((x & 1) << 7) | (x >> 1)
		state.cc.cy = b2i8(1 == (x & 1))
		break
	case 0x11: // LXI  D,D16
	case 0x13: // INX  D
	case 0x19: // DAD  D
	case 0x1a: // LDAX  D
	case 0x1f: // RAR
		x := state.a
		state.a = (state.cc.cy << 7) | (x >> 1)
		state.cc.cy = b2i8(1 == (x & 1))
		break
	case 0x21: // LXI  H,D16
	case 0x23: // INX  H
	case 0x26: // MVI  H,D8
	case 0x2f: // CMA (not)
		state.a = state.a ^ state.a
		break
	case 0x29: // DAD  H
	case 0x31: // LXI  SP,D16
	case 0x32: // STA  adr
	case 0x36: // MVI  M,D16
	case 0x3a: // LDA  adr
	case 0x3e: // MVI  A,D8
	case 0x41:
		state.b = state.c
		break
	case 0x42:
		state.b = state.d
		break
	case 0x43:
		state.b = state.e
		break
	case 0x56: // MOV  D,M
	case 0x5e: // MOV  E,M
	case 0x66: // MOV  H,M
	case 0x6f: // MOV  L,A
	case 0x77: // MOV  M,A
	case 0x7a: // MOV  A,D
	case 0x7b: // MOV  A,E
	case 0x7c: // MOV  A,H
	case 0x7e: // MOV  A,M
	case 0x80: // Add  b
		add(state, state.b)
		break
	case 0x81: // Add  c
		add(state, state.c)
		break
	case 0xa7: // ANA  A
	case 0xaf: // XRA  A
	case 0xc1: // Pop  b
		state.c = state.memory[state.sp]
		state.b = state.memory[state.sp+1]
		state.sp += 2
		break
	case 0xc2: // JNZ  adr
		if state.cc.z == 0 {
			state.pc = uint16(
				(state.memory[state.pc+2] << 8) | state.memory[state.pc+1])
		} else {
			state.pc += 2
		}
		break
	case 0xc3: // JMP  adr
		state.pc = uint16(
			(state.memory[state.pc+2] << 8) | state.memory[state.pc+1])
		break
	case 0xc5: // PUSH B
		state.memory[state.sp-1] = state.b
		state.memory[state.sp-2] = state.c
		state.sp -= 2
		break
	case 0xc6: // ADI  D6
	case 0xc9: // RET
		state.pc = uint16(
			state.memory[state.sp] | (state.memory[state.sp+1] << 8))
		state.sp += 2
		break
	case 0xcd: // Call Address
		ret := uint16(state.pc + 2)
		state.memory[state.sp-1] = uint8((ret >> 8) & 0xff)
		state.memory[state.sp-2] = uint8(ret & 0xff)
		state.sp = state.sp - 2
		state.pc = uint16(
			(state.memory[state.pc+2] << 8) | state.memory[state.pc+1])
		break
	case 0xd1: // POP  D
	case 0xd3: // OUT  D8
	case 0xd5: // PUSH D
	case 0xe1: // POP  H
	case 0xe5: // PUSH H
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
	case 0xfb: // EI
	case 0xfe: // CPI
		x := state.a - state.memory[state.pc+1]
		state.cc.z = b2i8(x == 0)
		state.cc.s = b2i8(0x80 == (x & 0x80))
		state.cc.p = parity(x, 8)
		state.cc.cy = b2i8(state.a > state.memory[state.pc+1])
		state.pc++
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

// Increment Register or Memory
func inr(state *State8080, val uint8) uint8 {
	result :=

	state.pc++
}

// Decrement Register or Memory
func dcr(state *State8080) {

}

// Complement Accumulator
func cma(state *State8080) {

}

// Decimal Adjust Accumulator
func daa(state *State8080) {

}

//endregion

// region Data Transfer Instructions

// MOV Instruction
func mov(state *State8080) {

}

// Store Accumulator
func stax(state *State8080) {

}

// Load Accumulator
func ldax(state *State8080) {

}

// endregion

// region Register/Memory to Accumulator Instructions

// Add reg/mem to Accumulator
func add(state *State8080, toAdd uint8) {
	answer := uint16(state.a) + uint16(toAdd)
	state.cc.z = b2i8((answer & 0xff) == 0)
	state.cc.s = b2i8((answer & 0x80) != 0)
	state.cc.cy = b2i8(answer > 0xff)
	state.cc.p = parity(answer & 0xff)
	state.a = uint8(answer & 0xff)
}

// Add reg/mem to Accumulator w/carry
func adc() {}

// Sub reg/mem from Accumulator
func sub() {}

// Sub reg/mem from Accumulator w/carry
func sbb() {}

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

// Push Data Onto Stack
func push(state *State8080) {}

// Pop Data Off Stack
func pop(state *State8080) {}

// Double Add
func dad(state *State8080) {}

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
