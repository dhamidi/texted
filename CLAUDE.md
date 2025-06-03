# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`texted` is a scriptable, headless text editor written in Go, designed for making automated edits to files on disk. The editing language is based on Emacs commands and supports multiple input formats (shell-like syntax, S-expressions, and JSON).

## Development Commands

- **Build**: `go build`
- **Test**: `go test ./...`
- **Run**: `go run .`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`

## Architecture

### Core Components

- **Buffer**: The fundamental building block containing UTF-8 encoded text with a `point` (cursor position) and `mark` (secondary position forming a region with the point)
- **Values System**: Type-safe value representation with `Value` and `ValueKind` interfaces supporting symbols, numbers, lists, and strings
- **Script Execution**: Programs can be written as shell-like commands, S-expressions, or JSON arrays

### Key APIs

- `texted.NewBuffer("initial contents")` - Creates a new buffer
- `buf.Do(script)` - Executes a texted script on the buffer
- `buf.String()` - Returns buffer contents
- `buf.Region()` - Returns selected region contents
- `buf.Point()`, `buf.Mark()` - Position accessors

### Script Formats

1. **Shell-like**: `search-forward "text"`
2. **S-expression**: `(search-forward "text")`  
3. **JSON**: `["search-forward", "text"]`

## Technical Constraints

- Use only Go standard library
- Implement I/O-less designs using appropriate interfaces (io.Reader, io.Writer)
- The main `texted` package is defined in the repository root directory
- Go version: 1.24.3

## Example Operations

- `search-forward "pattern"` - Search for text forward
- `set-mark` - Set the mark at current point
- `replace-region "text"` - Replace selected region
- `replace-match "text"` - Replace last search match