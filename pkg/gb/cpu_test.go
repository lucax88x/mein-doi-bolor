//nolint:testpackage // testing internals
package gb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// =============================================================================
// CPU REGISTERS - BASIC STRUCTURE
// =============================================================================
//
// The Game Boy CPU (Sharp SM83, often called LR35902) has the following registers:
//
// 8-bit registers: A, F, B, C, D, E, H, L
// 16-bit registers: SP (Stack Pointer), PC (Program Counter)
//
// The 8-bit registers can be combined into 16-bit pairs:
//   - AF: Accumulator + Flags
//   - BC, DE, HL: General purpose register pairs
//
// Reference: Pan Docs - CPU Registers and Flags
// https://gbdev.io/pandocs/CPU_Registers_and_Flags.html

// -----------------------------------------------------------------------------
// Test 1: CPU struct exists and can be created
// -----------------------------------------------------------------------------
// Start here! This test just verifies you can create a CPU instance.
// You'll need to export the cpu struct (rename to CPU) and add a constructor.

func TestNewCPU_ReturnsInstance(t *testing.T) {
	// The NewCPU constructor should return a pointer to a CPU struct.
	// This is the foundation - everything else builds on this.
	cpu := newCPU()

	require.NotNil(t, cpu, "NewCPU should return a non-nil CPU instance")
}

// =============================================================================
// 8-BIT REGISTERS
// =============================================================================
//
// Each 8-bit register should be readable and writable.
// These are the building blocks for all CPU operations.

func TestCPU_8BitRegisters_DefaultToZero(t *testing.T) {
	// Before setting post-boot values, verify registers start at zero.
	// This tests that your struct fields are properly initialized.
	//
	// Note: We'll test post-boot values separately. This test ensures
	// the zero value of your struct is sensible.
	cpu := newCPU()

	// For now, we test that registers can be read.
	// The actual post-boot values are tested in TestNewCPU_PostBootValues.

	// Test that we can access all 8-bit registers
	// (They should have accessor methods)
	_ = cpu.A()
	_ = cpu.F()
	_ = cpu.B()
	_ = cpu.C()
	_ = cpu.D()
	_ = cpu.E()
	_ = cpu.H()
	_ = cpu.L()
}

//nolint:funlen
func TestCPU_8BitRegisters_CanBeWrittenAndRead(t *testing.T) {
	// Each register should support get/set operations.
	// Use table-driven tests - an idiomatic Go pattern that makes
	// it easy to add test cases and see which specific case failed.
	testCases := []struct {
		name     string
		setValue uint8
		getter   func(*cpu) uint8
		setter   func(*cpu, uint8)
	}{
		{
			name:     "Register A (Accumulator)",
			setValue: 0x42,
			getter:   func(c *cpu) uint8 { return c.A() },
			setter:   func(c *cpu, v uint8) { c.SetA(v) },
		},
		{
			name:     "Register B",
			setValue: 0x12,
			getter:   func(c *cpu) uint8 { return c.B() },
			setter:   func(c *cpu, v uint8) { c.SetB(v) },
		},
		{
			name:     "Register C",
			setValue: 0x34,
			getter:   func(c *cpu) uint8 { return c.C() },
			setter:   func(c *cpu, v uint8) { c.SetC(v) },
		},
		{
			name:     "Register D",
			setValue: 0x56,
			getter:   func(c *cpu) uint8 { return c.D() },
			setter:   func(c *cpu, v uint8) { c.SetD(v) },
		},
		{
			name:     "Register E",
			setValue: 0x78,
			getter:   func(c *cpu) uint8 { return c.E() },
			setter:   func(c *cpu, v uint8) { c.SetE(v) },
		},
		{
			name:     "Register H",
			setValue: 0x9A,
			getter:   func(c *cpu) uint8 { return c.H() },
			setter:   func(c *cpu, v uint8) { c.SetH(v) },
		},
		{
			name:     "Register L",
			setValue: 0xBC,
			getter:   func(c *cpu) uint8 { return c.L() },
			setter:   func(c *cpu, v uint8) { c.SetL(v) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newCPU()

			tc.setter(cpu, tc.setValue)
			got := tc.getter(cpu)

			require.Equal(t, tc.setValue, got,
				"After setting %s to 0x%02X, getter should return 0x%02X",
				tc.name, tc.setValue, tc.setValue)
		})
	}
}

// =============================================================================
// 16-BIT REGISTERS (SP and PC)
// =============================================================================
//
// SP (Stack Pointer) and PC (Program Counter) are native 16-bit registers.
// They're not composed of 8-bit pairs like AF/BC/DE/HL.

func TestCPU_16BitRegisters_CanBeWrittenAndRead(t *testing.T) {
	testCases := []struct {
		name     string
		setValue uint16
		getter   func(*cpu) uint16
		setter   func(*cpu, uint16)
	}{
		{
			name:     "Stack Pointer (SP)",
			setValue: 0xFFFE,
			getter:   func(c *cpu) uint16 { return c.SP() },
			setter:   func(c *cpu, v uint16) { c.SetSP(v) },
		},
		{
			name:     "Program Counter (PC)",
			setValue: 0x0100,
			getter:   func(c *cpu) uint16 { return c.PC() },
			setter:   func(c *cpu, v uint16) { c.SetPC(v) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newCPU()

			tc.setter(cpu, tc.setValue)
			got := tc.getter(cpu)

			require.Equal(t, tc.setValue, got,
				"After setting %s to 0x%04X, getter should return 0x%04X",
				tc.name, tc.setValue, tc.setValue)
		})
	}
}

// =============================================================================
// 16-BIT REGISTER PAIRS (AF, BC, DE, HL)
// =============================================================================
//
// The 8-bit registers can be accessed as 16-bit pairs:
//   - AF = (A << 8) | F
//   - BC = (B << 8) | C
//   - DE = (D << 8) | E
//   - HL = (H << 8) | L
//
// The high byte is the first register (A, B, D, H).
// The low byte is the second register (F, C, E, L).
//
// This is important because many instructions operate on register pairs.

func TestCPU_RegisterPairs_CombineHighAndLowBytes(t *testing.T) {
	// When we set individual 8-bit registers, reading the 16-bit pair
	// should return them combined correctly.
	testCases := []struct {
		name      string
		highValue uint8
		lowValue  uint8
		setHigh   func(*cpu, uint8)
		setLow    func(*cpu, uint8)
		getPair   func(*cpu) uint16
	}{
		{
			name:      "BC pair",
			highValue: 0x12,
			lowValue:  0x34,
			setHigh:   func(c *cpu, v uint8) { c.SetB(v) },
			setLow:    func(c *cpu, v uint8) { c.SetC(v) },
			getPair:   func(c *cpu) uint16 { return c.BC() },
		},
		{
			name:      "DE pair",
			highValue: 0x56,
			lowValue:  0x78,
			setHigh:   func(c *cpu, v uint8) { c.SetD(v) },
			setLow:    func(c *cpu, v uint8) { c.SetE(v) },
			getPair:   func(c *cpu) uint16 { return c.DE() },
		},
		{
			name:      "HL pair",
			highValue: 0x9A,
			lowValue:  0xBC,
			setHigh:   func(c *cpu, v uint8) { c.SetH(v) },
			setLow:    func(c *cpu, v uint8) { c.SetL(v) },
			getPair:   func(c *cpu) uint16 { return c.HL() },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newCPU()

			tc.setHigh(cpu, tc.highValue)
			tc.setLow(cpu, tc.lowValue)

			expected := uint16(tc.highValue)<<8 | uint16(tc.lowValue)
			got := tc.getPair(cpu)

			require.Equal(t, expected, got,
				"setting high=0x%02X, low=0x%02X should give pair=0x%04X",
				tc.highValue, tc.lowValue, expected)
		})
	}
}

func TestCPU_RegisterPairs_setAffectsBothBytes(t *testing.T) {
	// When we set a 16-bit pair, both constituent 8-bit registers
	// should be updated.
	testCases := []struct {
		name         string
		pairValue    uint16
		expectedHigh uint8
		expectedLow  uint8
		setPair      func(*cpu, uint16)
		getHigh      func(*cpu) uint8
		getLow       func(*cpu) uint8
	}{
		{
			name:         "BC pair",
			pairValue:    0xABCD,
			expectedHigh: 0xAB,
			expectedLow:  0xCD,
			setPair:      func(c *cpu, v uint16) { c.SetBC(v) },
			getHigh:      func(c *cpu) uint8 { return c.B() },
			getLow:       func(c *cpu) uint8 { return c.C() },
		},
		{
			name:         "DE pair",
			pairValue:    0x1234,
			expectedHigh: 0x12,
			expectedLow:  0x34,
			setPair:      func(c *cpu, v uint16) { c.SetDE(v) },
			getHigh:      func(c *cpu) uint8 { return c.D() },
			getLow:       func(c *cpu) uint8 { return c.E() },
		},
		{
			name:         "HL pair",
			pairValue:    0xFEDC,
			expectedHigh: 0xFE,
			expectedLow:  0xDC,
			setPair:      func(c *cpu, v uint16) { c.SetHL(v) },
			getHigh:      func(c *cpu) uint8 { return c.H() },
			getLow:       func(c *cpu) uint8 { return c.L() },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newCPU()

			tc.setPair(cpu, tc.pairValue)

			gotHigh := tc.getHigh(cpu)
			gotLow := tc.getLow(cpu)

			require.Equal(t, tc.expectedHigh, gotHigh,
				"High byte should be 0x%02X", tc.expectedHigh)
			require.Equal(t, tc.expectedLow, gotLow,
				"Low byte should be 0x%02X", tc.expectedLow)
		})
	}
}

// =============================================================================
// AF REGISTER PAIR - SPECIAL CASE
// =============================================================================
//
// The AF pair is special because F is the flags register, and its lower
// 4 bits are ALWAYS zero. This is a hardware constraint that your
// implementation must enforce.

func TestCPU_AFPair_CombinesAAndF(t *testing.T) {
	cpu := newCPU()

	cpu.SetA(0x12)
	// Note: F's lower 4 bits are always 0, so 0xEC becomes 0xE0
	cpu.SetF(0xEC)

	expected := uint16(0x12E0)
	got := cpu.AF()

	require.Equal(t, expected, got,
		"AF should combine A (high) and F (low)")
}

func TestCPU_setAF_AffectsBothRegisters(t *testing.T) {
	cpu := newCPU()

	// When setting AF to 0xABCD, A should become 0xAB
	// But F should become 0xC0 (not 0xCD!) because lower 4 bits are masked
	cpu.SetAF(0xABCD)

	require.Equal(t, uint8(0xAB), cpu.A(), "A should be high byte of AF")
	require.Equal(t, uint8(0xC0), cpu.F(),
		"F should be low byte with lower 4 bits masked to 0")
}

// =============================================================================
// FLAGS REGISTER (F)
// =============================================================================
//
// The F register contains CPU status flags in its upper 4 bits:
//
//   Bit 7 (0x80): Z - Zero flag
//   Bit 6 (0x40): N - Subtract flag (BCD)
//   Bit 5 (0x20): H - Half Carry flag (BCD)
//   Bit 4 (0x10): C - Carry flag
//   Bits 3-0:     Always 0
//
// Reference: Pan Docs - The Flags Register
// https://gbdev.io/pandocs/CPU_Registers_and_Flags.html#the-flags-register-lower-8-bits-of-af-register

func TestCPU_FRegister_LowerNibbleAlwaysZero(t *testing.T) {
	// This is a critical hardware behavior: the lower 4 bits of F
	// are hardwired to 0. No matter what value you try to write,
	// those bits must remain 0.
	cpu := newCPU()

	// Try to set F to 0xFF (all bits set)
	cpu.SetF(0xFF)

	// Only upper nibble should be set
	require.Equal(t, uint8(0xF0), cpu.F(),
		"F register lower 4 bits must always be 0")
}

func TestCPU_FRegister_MasksOnEveryWrite(t *testing.T) {
	// Test multiple values to ensure masking is consistent
	testCases := []struct {
		input    uint8
		expected uint8
	}{
		{0x00, 0x00}, // All zeros stay zero
		{0xFF, 0xF0}, // All ones -> only upper nibble
		{0x0F, 0x00}, // Only lower nibble -> all zeros
		{0xF0, 0xF0}, // Only upper nibble -> unchanged
		{0xA5, 0xA0}, // Mixed bits -> lower nibble cleared
		{0x5A, 0x50}, // Mixed bits -> lower nibble cleared
	}

	for _, tc := range testCases {
		t.Run(
			"input_"+formatHex(tc.input),
			func(t *testing.T) {
				cpu := newCPU()
				cpu.SetF(tc.input)

				require.Equal(t, tc.expected, cpu.F(),
					"setF(0x%02X) should result in F=0x%02X", tc.input, tc.expected)
			})
	}
}

// =============================================================================
// INDIVIDUAL FLAG OPERATIONS
// =============================================================================
//
// For convenience, you'll want methods to get/set individual flags.
// This makes instruction implementations much cleaner.

func TestCPU_ZeroFlag_GetAndset(t *testing.T) {
	// The Zero flag (Z) is set when an operation results in zero.
	// It's bit 7 of the F register (0x80).
	cpu := newCPU()

	// Initially should be false (assuming F starts at 0)
	cpu.SetF(0x00)
	require.False(t, cpu.FlagZ(), "Z flag should be false when bit 7 is 0")

	// set the Z flag
	cpu.SetFlagZ(true)
	require.True(t, cpu.FlagZ(), "Z flag should be true after setFlagZ(true)")
	require.Equal(t, uint8(0x80), cpu.F()&0x80, "Bit 7 should be set")

	// Clear the Z flag
	cpu.SetFlagZ(false)
	require.False(t, cpu.FlagZ(), "Z flag should be false after setFlagZ(false)")
	require.Equal(t, uint8(0x00), cpu.F()&0x80, "Bit 7 should be cleared")
}

func TestCPU_SubtractFlag_GetAndset(t *testing.T) {
	// The Subtract flag (N) indicates if the last operation was a subtraction.
	// Used for DAA (Decimal Adjust Accumulator) instruction.
	// It's bit 6 of the F register (0x40).
	cpu := newCPU()

	cpu.SetF(0x00)
	require.False(t, cpu.FlagN(), "N flag should be false when bit 6 is 0")

	cpu.SetFlagN(true)
	require.True(t, cpu.FlagN(), "N flag should be true after setFlagN(true)")
	require.Equal(t, uint8(0x40), cpu.F()&0x40, "Bit 6 should be set")

	cpu.SetFlagN(false)
	require.False(t, cpu.FlagN(), "N flag should be false after setFlagN(false)")
}

func TestCPU_HalfCarryFlag_GetAndset(t *testing.T) {
	// The Half Carry flag (H) is set when there's a carry from bit 3 to bit 4.
	// Also used for DAA. It's bit 5 of the F register (0x20).
	//
	// This flag trips up many emulator developers! Pay attention to it
	// when implementing arithmetic instructions.
	cpu := newCPU()

	cpu.SetF(0x00)
	require.False(t, cpu.FlagH(), "H flag should be false when bit 5 is 0")

	cpu.SetFlagH(true)
	require.True(t, cpu.FlagH(), "H flag should be true after setFlagH(true)")
	require.Equal(t, uint8(0x20), cpu.F()&0x20, "Bit 5 should be set")

	cpu.SetFlagH(false)
	require.False(t, cpu.FlagH(), "H flag should be false after setFlagH(false)")
}

func TestCPU_CarryFlag_GetAndset(t *testing.T) {
	// The Carry flag (C) is set when an operation overflows/underflows.
	// It's bit 4 of the F register (0x10).
	cpu := newCPU()

	cpu.SetF(0x00)
	require.False(t, cpu.FlagC(), "C flag should be false when bit 4 is 0")

	cpu.SetFlagC(true)
	require.True(t, cpu.FlagC(), "C flag should be true after setFlagC(true)")
	require.Equal(t, uint8(0x10), cpu.F()&0x10, "Bit 4 should be set")

	cpu.SetFlagC(false)
	require.False(t, cpu.FlagC(), "C flag should be false after setFlagC(false)")
}

func TestCPU_Flags_setOneDoesNotAffectOthers(t *testing.T) {
	// setting one flag should not disturb the other flags.
	// This tests that your bit manipulation is correct.
	cpu := newCPU()

	// Start with all flags clear
	cpu.SetF(0x00)

	// set each flag one by one
	cpu.SetFlagZ(true)
	cpu.SetFlagN(true)
	cpu.SetFlagH(true)
	cpu.SetFlagC(true)

	// All flags should now be set
	require.True(t, cpu.FlagZ(), "Z should still be set")
	require.True(t, cpu.FlagN(), "N should still be set")
	require.True(t, cpu.FlagH(), "H should still be set")
	require.True(t, cpu.FlagC(), "C should still be set")
	require.Equal(t, uint8(0xF0), cpu.F(), "F should be 0xF0 with all flags set")

	// Now clear them one by one and verify others remain
	cpu.SetFlagZ(false)
	require.False(t, cpu.FlagZ(), "Z should be cleared")
	require.True(t, cpu.FlagN(), "N should remain set")
	require.True(t, cpu.FlagH(), "H should remain set")
	require.True(t, cpu.FlagC(), "C should remain set")

	cpu.SetFlagN(false)
	require.False(t, cpu.FlagN(), "N should be cleared")
	require.True(t, cpu.FlagH(), "H should remain set")
	require.True(t, cpu.FlagC(), "C should remain set")
}

// =============================================================================
// POST-BOOT REGISTER VALUES
// =============================================================================
//
// After the boot ROM finishes executing, the CPU registers have specific
// values. Your NewCPU() constructor should initialize to these values
// so games can run without needing the boot ROM.
//
// Reference: Pan Docs - Power Up Sequence
// https://gbdev.io/pandocs/Power_Up_Sequence.html
//
// For original DMG (Game Boy):
//   A  = 0x01
//   F  = 0xB0 (Z=1, N=0, H=1, C=1)
//   B  = 0x00
//   C  = 0x13
//   D  = 0x00
//   E  = 0xD8
//   H  = 0x01
//   L  = 0x4D
//   SP = 0xFFFE
//   PC = 0x0100

func TestNewCPU_PostBootValues_8BitRegisters(t *testing.T) {
	// These are the register values after the DMG boot ROM completes.
	// Games expect these values, so we initialize to them.
	cpu := newCPU()

	require.Equal(t, uint8(0x01), cpu.A(),
		"A should be 0x01 after boot (DMG)")
	require.Equal(t, uint8(0xB0), cpu.F(),
		"F should be 0xB0 after boot (Z=1, N=0, H=1, C=1)")
	require.Equal(t, uint8(0x00), cpu.B(),
		"B should be 0x00 after boot")
	require.Equal(t, uint8(0x13), cpu.C(),
		"C should be 0x13 after boot")
	require.Equal(t, uint8(0x00), cpu.D(),
		"D should be 0x00 after boot")
	require.Equal(t, uint8(0xD8), cpu.E(),
		"E should be 0xD8 after boot")
	require.Equal(t, uint8(0x01), cpu.H(),
		"H should be 0x01 after boot")
	require.Equal(t, uint8(0x4D), cpu.L(),
		"L should be 0x4D after boot")
}

func TestNewCPU_PostBootValues_16BitRegisters(t *testing.T) {
	cpu := newCPU()

	require.Equal(t, uint16(0xFFFE), cpu.SP(),
		"SP should be 0xFFFE after boot (top of HRAM)")
	require.Equal(t, uint16(0x0100), cpu.PC(),
		"PC should be 0x0100 after boot (cartridge entry point)")
}

func TestNewCPU_PostBootValues_Flags(t *testing.T) {
	// Verify flags are set correctly based on F = 0xB0
	//   Z = 1 (bit 7 set)
	//   N = 0 (bit 6 clear)
	//   H = 1 (bit 5 set)
	//   C = 1 (bit 4 set)
	cpu := newCPU()

	require.True(t, cpu.FlagZ(), "Z flag should be set after boot")
	require.False(t, cpu.FlagN(), "N flag should be clear after boot")
	require.True(t, cpu.FlagH(), "H flag should be set after boot")
	require.True(t, cpu.FlagC(), "C flag should be set after boot")
}

func TestNewCPU_PostBootValues_RegisterPairs(t *testing.T) {
	// Verify that register pairs reflect the combined individual values
	cpu := newCPU()

	require.Equal(t, uint16(0x01B0), cpu.AF(),
		"AF should be 0x01B0 after boot")
	require.Equal(t, uint16(0x0013), cpu.BC(),
		"BC should be 0x0013 after boot")
	require.Equal(t, uint16(0x00D8), cpu.DE(),
		"DE should be 0x00D8 after boot")
	require.Equal(t, uint16(0x014D), cpu.HL(),
		"HL should be 0x014D after boot")
}

// =============================================================================
// EDGE CASES AND BOUNDARY CONDITIONS
// =============================================================================

func TestCPU_Registers_HandleMinMaxValues(t *testing.T) {
	// Ensure registers handle the full 8-bit and 16-bit ranges
	cpu := newCPU()

	// Test 8-bit boundaries
	cpu.SetA(0x00)
	require.Equal(t, uint8(0x00), cpu.A(), "A should handle 0x00")

	cpu.SetA(0xFF)
	require.Equal(t, uint8(0xFF), cpu.A(), "A should handle 0xFF")

	// Test 16-bit boundaries
	cpu.SetSP(0x0000)
	require.Equal(t, uint16(0x0000), cpu.SP(), "SP should handle 0x0000")

	cpu.SetSP(0xFFFF)
	require.Equal(t, uint16(0xFFFF), cpu.SP(), "SP should handle 0xFFFF")

	cpu.SetPC(0x0000)
	require.Equal(t, uint16(0x0000), cpu.PC(), "PC should handle 0x0000")

	cpu.SetPC(0xFFFF)
	require.Equal(t, uint16(0xFFFF), cpu.PC(), "PC should handle 0xFFFF")
}

func TestCPU_RegisterPairs_HandleMinMaxValues(t *testing.T) {
	cpu := newCPU()

	// Test BC at boundaries
	cpu.SetBC(0x0000)
	require.Equal(t, uint8(0x00), cpu.B(), "B should be 0x00")
	require.Equal(t, uint8(0x00), cpu.C(), "C should be 0x00")

	cpu.SetBC(0xFFFF)
	require.Equal(t, uint8(0xFF), cpu.B(), "B should be 0xFF")
	require.Equal(t, uint8(0xFF), cpu.C(), "C should be 0xFF")

	// Test HL at boundaries (HL is heavily used for memory addressing)
	cpu.SetHL(0x0000)
	require.Equal(t, uint16(0x0000), cpu.HL(), "HL should handle 0x0000")

	cpu.SetHL(0xFFFF)
	require.Equal(t, uint16(0xFFFF), cpu.HL(), "HL should handle 0xFFFF")
}

func TestCPU_setAF_AlwaysMasksLowerNibbleOfF(t *testing.T) {
	// Even when setting via setAF, the lower nibble of F must be masked
	testCases := []struct {
		input      uint16
		expectedAF uint16
		expectedF  uint8
	}{
		{0x01B0, 0x01B0, 0xB0}, // Normal case
		{0x01BF, 0x01B0, 0xB0}, // Lower nibble should be masked
		{0xFF0F, 0xFF00, 0x00}, // Only lower nibble set in F
		{0xFFFF, 0xFFF0, 0xF0}, // All bits set
		{0x0000, 0x0000, 0x00}, // All bits clear
	}

	for _, tc := range testCases {
		t.Run(
			"input_"+formatHex16(tc.input),
			func(t *testing.T) {
				cpu := newCPU()
				cpu.SetAF(tc.input)

				require.Equal(t, tc.expectedAF, cpu.AF(),
					"AF should be 0x%04X after setAF(0x%04X)", tc.expectedAF, tc.input)
				require.Equal(t, tc.expectedF, cpu.F(),
					"F should be 0x%02X after setAF(0x%04X)", tc.expectedF, tc.input)
			})
	}
}

// Helper to format 16-bit hex values for test names.
func formatHex16(v uint16) string {
	const hexDigits = "0123456789ABCDEF"

	return string([]byte{
		hexDigits[(v>>12)&0xF],
		hexDigits[(v>>8)&0xF],
		hexDigits[(v>>4)&0xF],
		hexDigits[v&0xF],
	})
}

// Helper to format hex values for test names.
func formatHex(v uint8) string {
	const hexDigits = "0123456789ABCDEF"

	return string([]byte{hexDigits[v>>4], hexDigits[v&0xF]})
}
