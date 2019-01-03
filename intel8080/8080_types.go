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
