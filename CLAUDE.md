# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Game Boy emulator written in Go, following the Pan Docs specification (https://gbdev.io/pandocs/). The project is in early stages with basic CPU structure in place.

## Build Commands

The project uses pnpm as a task runner with Go as the underlying language:

```bash
pnpm build          # Build binary to ./bin/cli
pnpm dev            # Run without building
pnpm test           # Run all tests with verbose output and shuffle
pnpm lint           # Run golangci-lint
pnpm lint:fix       # Run golangci-lint with auto-fix
pnpm format         # Check formatting with gofmt
pnpm format:fix     # Fix formatting with gofmt
```

Direct Go commands:
```bash
go test -v ./pkg/gb/...              # Run tests for a specific package
go test -v -run TestCpu ./pkg/gb/    # Run a single test
```

## Architecture

```
cmd/cli/        # CLI entrypoint
pkg/gb/         # Game Boy emulation core
  cpu.go        # CPU implementation (Sharp LR35902)
```

The emulator follows the Pan Docs hardware reference. Key Game Boy components to implement:
- **CPU**: Sharp LR35902 (SM83) - 8-bit CPU with 16-bit addressing
- **PPU**: Picture Processing Unit for graphics
- **APU**: Audio Processing Unit
- **Memory**: Memory bus with ROM/RAM banking

## Tool Versions

Managed via `mise.toml`:
- Go 1.25.5
- golangci-lint (latest compatible with Go version)

Run `pnpm postinstall` to verify Go/golangci-lint version compatibility.

## Testing

Uses `github.com/stretchr/testify` for assertions. Tests use the `_test` package suffix pattern (black-box testing).
