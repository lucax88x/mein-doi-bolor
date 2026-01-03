package gb

import "fmt"

const (
	lowByteMask = 0xFF // Mask for extracting the low byte of a 16-bit value
	flagMask    = 0xF0 // F register lower 4 bits are always 0
	flagZ       = 0x80 // Zero flag (bit 7)
)

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

func (c *cpu) SetA(value uint8) {
	c.a = value
}

func (c *cpu) F() uint8 {
	return c.f
}

func (c *cpu) SetF(value uint8) {
	c.f = value & flagMask
}

func (c *cpu) B() uint8 {
	return c.b
}

func (c *cpu) SetB(value uint8) {
	c.b = value
}

func (c *cpu) C() uint8 {
	return c.c
}

func (c *cpu) SetC(value uint8) {
	c.c = value
}

func (c *cpu) D() uint8 {
	return c.d
}

func (c *cpu) SetD(value uint8) {
	c.d = value
}

func (c *cpu) E() uint8 {
	return c.e
}

func (c *cpu) SetE(value uint8) {
	c.e = value
}

func (c *cpu) H() uint8 {
	return c.h
}

func (c *cpu) SetH(value uint8) {
	c.h = value
}

func (c *cpu) L() uint8 {
	return c.l
}

func (c *cpu) SetL(value uint8) {
	c.l = value
}

func (c *cpu) SP() uint16 {
	return c.sp
}

func (c *cpu) SetSP(value uint16) {
	c.sp = value
}

func (c *cpu) PC() uint16 {
	return c.pc
}

func (c *cpu) SetPC(value uint16) {
	c.pc = value
}

func (c *cpu) AF() uint16 {
	maskedF := c.f & flagMask

	return uint16(c.a)<<8 | uint16(maskedF)
}

func (c *cpu) SetAF(value uint16) {
	c.a = uint8(value >> 8)
	f := uint8(value & flagMask)

	c.f = f
}

func (c *cpu) BC() uint16 {
	return uint16(c.b)<<8 | uint16(c.c)
}

func (c *cpu) SetBC(value uint16) {
	c.b = uint8(value >> 8)
	c.c = uint8(value & lowByteMask)
}

func (c *cpu) DE() uint16 {
	return uint16(c.d)<<8 | uint16(c.e)
}

func (c *cpu) SetDE(value uint16) {
	c.d = uint8(value >> 8)
	c.e = uint8(value & lowByteMask)
}

func (c *cpu) HL() uint16 {
	return uint16(c.h)<<8 | uint16(c.l)
}

func (c *cpu) SetHL(value uint16) {
	c.h = uint8(value >> 8)
	c.l = uint8(value & lowByteMask)
}

func (c *cpu) FlagZ() bool {
	return c.F()&0x80 != 0
}

func (c *cpu) SetFlagZ(value bool) {
	if value {
		c.SetF(c.F() | flagZ)
	} else {
		c.SetF(c.F() &^ flagZ)
	}
}

func debug(v uint) {
	fmt.Printf("0x%X\n", v)
}
