package gb

import "fmt"

// Initial CPU register values after boot ROM (DMG).
const (
	initA  uint8  = 0x01
	initF  uint8  = 0xB0
	initB  uint8  = 0x00
	initC  uint8  = 0x13
	initD  uint8  = 0x00
	initE  uint8  = 0xD8
	initH  uint8  = 0x01
	initL  uint8  = 0x4D
	initSP uint16 = 0xFFFE
	initPC uint16 = 0x0100
)

const (
	lowByteMask    = 0xFF // Mask for extracting the low byte of a 16-bit value
	flagMask       = 0xF0 // F register lower 4 bits are always 0
	flagZ          = 0x80 // Zero flag (bit 7)
	flagN          = 0x40 // Subtract flag (bit 6)
	flagH          = 0x20 // Half Carry flag (bit 5)
	flagC          = 0x10 // Carry flag (bit 4)
	opCode_NOP     = 0x00
	opCode_LD_B_nr = 0x06
	opCode_LD_C_nr = 0x0E
	opCode_LD_D_nr = 0x16
	opCode_LD_E_nr = 0x1E
	opCode_LD_H_nr = 0x26
	opCode_LD_L_nr = 0x2E
	opCode_LD_A_nr = 0x3E
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

	sp uint16 // stack pointer
	pc uint16 // program counter

	cycles int

	bus *bus
}

func newCPU() *cpu {
	return &cpu{
		a:      initA,
		f:      initF,
		b:      initB,
		c:      initC,
		d:      initD,
		e:      initE,
		h:      initH,
		l:      initL,
		sp:     initSP,
		pc:     initPC,
		cycles: 0,
		bus:    newBus(),
	}
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
	return c.F()&flagZ != 0
}

func (c *cpu) SetFlagZ(value bool) {
	c.setFlag(value, flagZ)
}

func (c *cpu) FlagN() bool {
	return c.F()&flagN != 0
}

func (c *cpu) SetFlagN(value bool) {
	c.setFlag(value, flagN)
}

func (c *cpu) FlagH() bool {
	return c.F()&flagH != 0
}

func (c *cpu) SetFlagH(value bool) {
	c.setFlag(value, flagH)
}

func (c *cpu) FlagC() bool {
	return c.F()&flagC != 0
}

func (c *cpu) SetFlagC(value bool) {
	c.setFlag(value, flagC)
}

func (c *cpu) setFlag(value bool, flag uint8) {
	if value {
		c.SetF(c.F() | flag)
	} else {
		c.SetF(c.F() &^ flag)
	}
}

func (c *cpu) fetch() (uint8, error) {
	val, err := c.bus.Read(c.PC())

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetPC(c.PC() + 0x01)

	return val, nil
}

func (c *cpu) exec_LD_B_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetB(immediateValue)

	return 8, nil
}

func (c *cpu) exec_LD_C_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetC(immediateValue)

	return 8, nil
}

func (c *cpu) exec_LD_D_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetD(immediateValue)

	return 8, nil
}

func (c *cpu) exec_LD_E_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetE(immediateValue)

	return 8, nil
}

func (c *cpu) exec_LD_H_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetH(immediateValue)

	return 8, nil
}

func (c *cpu) exec_LD_L_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetL(immediateValue)

	return 8, nil
}

func (c *cpu) exec_LD_A_nr() (int, error) {
	immediateValue, err := c.fetch()

	if err != nil {
		return 0, fmt.Errorf("failed to read immediate value at PC+1: %v", err)
	}

	c.SetA(immediateValue)

	return 8, nil
}

func (c *cpu) exec(opCode uint8) (int, error) {
	switch opCode {
	case opCode_NOP:
		return 4, nil
	case opCode_LD_B_nr:
		return c.exec_LD_B_nr()
	case opCode_LD_C_nr:
		return c.exec_LD_C_nr()
	case opCode_LD_D_nr:
		return c.exec_LD_D_nr()
	case opCode_LD_E_nr:
		return c.exec_LD_E_nr()
	case opCode_LD_H_nr:
		return c.exec_LD_H_nr()
	case opCode_LD_L_nr:
		return c.exec_LD_L_nr()
	case opCode_LD_A_nr:
		return c.exec_LD_A_nr()
	default:
		return 0, fmt.Errorf("unimplemented opcode: 0x%02X", opCode)
	}
}

func debug(v uint) {
	// fmt.Printf("0x%X\n", v)
	// fmt.Printf("0x%02X\n", v) // 8-bit:  0x0f, 0xff
	fmt.Printf("0x%04X\n", v) // 16-bit: 0x0100, 0xffff
}

func (c *cpu) Step() (int, error) {
	opCode, err := c.fetch()

	if err != nil {
		return c.cycles, fmt.Errorf("failed to read opcode at PC: %v", err)
	}

	cycles, err := c.exec(opCode)
	c.cycles = cycles

	if err != nil {
		return c.cycles, fmt.Errorf("failed to read opcode at PC: %v", err)
	}

	return c.cycles, nil
}
