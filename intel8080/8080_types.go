package intel8080

type ConditionCodes struct {
	z   uint8
	s   uint8
	p   uint8
	cy  uint8
	ac  uint8
	pad uint8
}

func (c *ConditionCodes) GetZ() uint8 {
	return c.z & 0x1
}

func (c *ConditionCodes) GetS() uint8 {
	return c.s & 0x1
}

func (c *ConditionCodes) GetP() uint8 {
	return c.p & 0x1
}

func (c *ConditionCodes) GetCy() uint8 {
	return c.cy & 0x1
}

func (c *ConditionCodes) GetAc() uint8 {
	return c.ac & 0x1
}

func (c *ConditionCodes) GetPad() uint8 {
	return c.pad & 0x3
}

/*
 * Register Pairs
 * B   = B and C
 * D   = D and E
 * H   = H and L
 * PSW = A and State of condition bits
 */

type State8080 struct {
	a          uint8
	b          uint8
	c          uint8
	d          uint8
	e          uint8
	h          uint8
	l          uint8
	sp         uint16
	pc         uint16
	memory     []uint8
	cc         ConditionCodes
	int_enable uint8
}

// TODO: Figure out how memory is handled
func (state *State8080) wb(addr uint16, val uint8) {

}

func (state *State8080) rb(addr uint16) uint8 {
	return 0
}

func (state *State8080) setBC(val uint16) {
	state.b, state.c = uint8(val>>8), uint8(val&0xff)
}

func (state *State8080) getBC() uint16 {
	return uint16(state.b<<8) | uint16(state.c)
}

func (state *State8080) setDE(val uint16) {
	state.d, state.e = uint8(val>>8), uint8(val&0xff)
}

func (state *State8080) getDE() uint16 {
	return uint16(state.d<<8) | uint16(state.e)
}

func (state *State8080) setHL(val uint16) {
	state.h, state.l = uint8(val>>8), uint8(val&0xff)
}

func (state *State8080) getHL() uint16 {
	return uint16(state.h<<8) | uint16(state.l)
}
