package gb

import "fmt"

type bus struct {
	rom  []uint8
	vram []uint8
	wram []uint8
	hram []uint8
	ie   uint8
}

const Size32Kb = 0x8000
const Size8Kb = 0x2000
const Size127b = 0x7F

func newBus() *bus {
	return &bus{
		rom:  make([]uint8, Size32Kb),
		vram: make([]uint8, Size8Kb),
		wram: make([]uint8, Size8Kb),
		hram: make([]uint8, Size127b),
		ie:   uint8(0),
	}
}

func (b *bus) Read(addr uint16) uint8 {
	switch {
	case addr == 0xFFFF:
		return b.ie
	case addr <= 0x7FFF:
		return b.rom[addr]
	case addr >= 0x8000 && addr <= 0x9FFF:
		return b.vram[addr-0x8000]
	case addr >= 0xC000 && addr <= 0xDFFF:
		return b.wram[addr-0xC000]
	case addr >= 0xE000 && addr <= 0xFDFF:
		return b.wram[addr-0xE000]
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return b.hram[addr-0xFF80]
	}

	return 0xFF
}

func (b *bus) Write(addr uint16, value uint8) error {
	switch {
	case addr == 0xFFFF:
		b.ie = value
	case addr <= 0x7FFF:
		return fmt.Errorf("attempt to write to ROM address: 0x%04X", addr)
	case addr >= 0x8000 && addr <= 0x9FFF:
		b.vram[addr-0x8000] = value
	case addr >= 0xC000 && addr <= 0xDFFF:
		b.wram[addr-0xC000] = value
	case addr >= 0xE000 && addr <= 0xFDFF:
		b.wram[addr-0xE000] = value
	case addr >= 0xFF80 && addr <= 0xFFFE:
		b.hram[addr-0xFF80] = value
	}

	return nil
}

func (b *bus) LoadROM(data []uint8) error {
	if len(data) > len(b.rom) {
		return fmt.Errorf("ROM size exceeds maximum of %d bytes", len(b.rom))
	}

	for i := 0; i < len(b.rom) && i < len(data); i++ {
		b.rom[i] = data[i]
	}

	return nil
}
