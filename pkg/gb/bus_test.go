package gb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBus_NewBus(t *testing.T) {
	bus := newBus()
	require.NotNil(t, bus)
}

func TestBus_ReadWrite_WRAM(t *testing.T) {
	// WRAM: $C000-$DFFF (8 KB)
	tests := []struct {
		name  string
		addr  uint16
		value uint8
	}{
		{"WRAM start", 0xC000, 0x42},
		{"WRAM middle", 0xD000, 0xAB},
		{"WRAM end", 0xDFFF, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := newBus()
			bus.Write(tt.addr, tt.value)
			value, _ := bus.Read(tt.addr)
			require.Equal(t, tt.value, value)
		})
	}
}

func TestBus_ReadWrite_HRAM(t *testing.T) {
	// HRAM: $FF80-$FFFE (127 bytes)
	tests := []struct {
		name  string
		addr  uint16
		value uint8
	}{
		{"HRAM start", 0xFF80, 0x12},
		{"HRAM middle", 0xFFB0, 0x34},
		{"HRAM end", 0xFFFE, 0x56},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := newBus()
			bus.Write(tt.addr, tt.value)
			value, _ := bus.Read(tt.addr)
			require.Equal(t, tt.value, value)
		})
	}
}

func TestBus_ReadWrite_VRAM(t *testing.T) {
	// VRAM: $8000-$9FFF (8 KB)
	tests := []struct {
		name  string
		addr  uint16
		value uint8
	}{
		{"VRAM start", 0x8000, 0xAA},
		{"VRAM middle", 0x9000, 0xBB},
		{"VRAM end", 0x9FFF, 0xCC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := newBus()
			bus.Write(tt.addr, tt.value)
			value, _ := bus.Read(tt.addr)
			require.Equal(t, tt.value, value)
		})
	}
}

func TestBus_LoadROM(t *testing.T) {
	bus := newBus()
	rom := []byte{0x00, 0xC3, 0x50, 0x01} // NOP, JP $0150

	err := bus.LoadROM(rom)

	require.NoError(t, err)

	value, _ := bus.Read(0x0000)
	require.Equal(t, uint8(0x00), value)
	value, _ = bus.Read(0x0001)
	require.Equal(t, uint8(0xC3), value)
	value, _ = bus.Read(0x0002)
	require.Equal(t, uint8(0x50), value)
	value, _ = bus.Read(0x0003)
	require.Equal(t, uint8(0x01), value)
}

func TestBus_LoadROM_TooBig_Gives_Error(t *testing.T) {
	bus := newBus()
	rom := make([]byte, 0x8001) // One byte too large

	err := bus.LoadROM(rom)

	require.Error(t, err)
}

func TestBus_ROM_IsReadOnly(t *testing.T) {
	// ROM: $0000-$7FFF should not be writable
	bus := newBus()
	rom := []byte{0x42}
	bus.LoadROM(rom)

	// Attempt to write to ROM
	err := bus.Write(0x0000, 0xFF)

	// Has error
	require.Error(t, err, "should not write to ROM")

	// Value should remain unchanged
	value, _ := bus.Read(0x0000)
	require.Equal(t, uint8(0x42), value)
}

func TestBus_EchoRAM(t *testing.T) {
	// Echo RAM: $E000-$FDFF mirrors $C000-$DDFF
	bus := newBus()

	// Write to WRAM
	bus.Write(0xC000, 0x42)
	bus.Write(0xC123, 0xAB)

	// Should be readable from Echo RAM
	value, _ := bus.Read(0xE000)
	require.Equal(t, uint8(0x42), value)
	value, _ = bus.Read(0xE123)
	require.Equal(t, uint8(0xAB), value)

	// Write to Echo RAM should affect WRAM
	bus.Write(0xE456, 0xCD)
	value, _ = bus.Read(0xC456)
	require.Equal(t, uint8(0xCD), value)
}

func TestBus_IERegister(t *testing.T) {
	// IE Register: $FFFF (Interrupt Enable)
	bus := newBus()

	bus.Write(0xFFFF, 0x1F)
	value, _ := bus.Read(0xFFFF)
	require.Equal(t, uint8(0x1F), value)
}
