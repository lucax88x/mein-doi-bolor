package gb

import "fmt"

type cpu struct {
	a uint8
	f uint8
	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8

	sp uint16
	pc uint16
}

func newCPU() *cpu {
	return &cpu{}
}

func (c *cpu) A() uint8 {
	return c.a
}

func (c *cpu) setA(value uint8) {
	c.a = value
}

func (c *cpu) F() uint8 {
	return c.f
}

func (c *cpu) setF(value uint8) {
	c.f = value & 0xF0
}

func (c *cpu) B() uint8 {
	return c.b
}

func (c *cpu) setB(value uint8) {
	c.b = value
}

func (c *cpu) C() uint8 {
	return c.c
}

func (c *cpu) setC(value uint8) {
	c.c = value
}

func (c *cpu) D() uint8 {
	return c.d
}

func (c *cpu) setD(value uint8) {
	c.d = value
}

func (c *cpu) E() uint8 {
	return c.e
}

func (c *cpu) setE(value uint8) {
	c.e = value
}

func (c *cpu) H() uint8 {
	return c.h
}

func (c *cpu) setH(value uint8) {
	c.h = value
}

func (c *cpu) L() uint8 {
	return c.l
}

func (c *cpu) setL(value uint8) {
	c.l = value
}

func (c *cpu) SP() uint16 {
	return c.sp
}

func (c *cpu) setSP(value uint16) {
	c.sp = value
}

func (c *cpu) PC() uint16 {
	return c.pc
}

func (c *cpu) setPC(value uint16) {
	c.pc = value
}

func (c *cpu) AF() uint16 {
	maskedF := c.f & 0xF0

	return uint16(c.a)<<8 | uint16(maskedF)
}

func (c *cpu) setAF(value uint16) {
	c.a = uint8(value >> 8)
	f := uint8(value & 0xF0)

	c.f = f
}

func (c *cpu) BC() uint16 {
	return uint16(c.b)<<8 | uint16(c.c)
}

func (c *cpu) setBC(value uint16) {
	c.b = uint8(value >> 8)
	c.c = uint8(value & 0xFF)
}

func (c *cpu) DE() uint16 {
	return uint16(c.d)<<8 | uint16(c.e)
}

func (c *cpu) setDE(value uint16) {
	c.d = uint8(value >> 8)
	c.e = uint8(value & 0xFF)
}

func (c *cpu) HL() uint16 {
	return uint16(c.h)<<8 | uint16(c.l)
}

func (c *cpu) setHL(value uint16) {
	c.h = uint8(value >> 8)
	c.l = uint8(value & 0xFF)
}

func (c *cpu) FlagZ() bool {
	return c.F()&0x80 != 0
}

func (c *cpu) setFlagZ(value bool) {
	if value {
		c.setF(c.F() | 0x80)
	} else {
		c.setF(c.F() &^ 0x80)
	}
}

func debug(v uint) {
	fmt.Printf("0x%X\n", v)
}
