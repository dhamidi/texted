# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`texted` is a scriptable, headless text editor written in Go, designed for making automated edits to files on disk. The editing language is based on Emacs commands and supports multiple input formats (shell-like syntax, S-expressions, and JSON).

## Development Commands

- **Build**: `go build`
- **Test**: `go test ./edlisp`
- Integration tests: `go run ./cmd/texted test`
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

## Technical Constraints and Architectural Choices

### Dependencies and External Libraries

- **Minimal external dependencies**: Only Cobra CLI framework (`github.com/spf13/cobra v1.9.1`) and its dependencies (`mousetrap`, `pflag`)
- **Core implementation uses only Go standard library** - no external libraries for parsing, text processing, or data structures
- **Go version requirement**: 1.24.3 (as specified in go.mod)

### Value System Design

- **Interface-based type system**: `Value` and `ValueKind` interfaces with singleton patterns for type kinds
- **Four core value types**: Symbol, String, Number (float64-based), List
- **Type checking via string comparison** using `IsA(value, kind)` function
- **Immutable value semantics**: Operations return new values rather than modifying existing ones
- **Slice-backed lists** with copy-on-write semantics for append operations
- **No memory pooling**: Relies on Go's garbage collector for memory management

### Parser Architecture

- **Multi-format parsing strategy**: Unified parsing to common AST (edlisp.Value tree) from three input formats
- **Manual tokenization**: No external lexer libraries, hand-written recursive descent parsers
- **Streaming parser design**: All parsers use `io.Reader` interface for input
- **Line-based parsing** with automatic format detection (parenthesis detection for S-expressions vs shell-like)
- **Comprehensive error handling**: Structured error types with detailed context

### JSON Format Constraints

- **Strict schema validation**: Only arrays allowed at top level
- **Type mapping rules**: JSON strings→String values, JSON numbers→Number values, nested arrays→List values
- **First array element must be string** (becomes Symbol for function calls)
- **No support for**: JSON booleans, null values, or objects

### I/O-Less Design Pattern

- **All core operations use standard interfaces**: `io.Reader`, `io.Writer`
- **No direct file system access** in core libraries
- **Streaming-friendly**: Parsers work with readers, not strings
- **Testable design**: Easy to mock I/O for testing

### Package Structure

- **`edlisp/`**: Core value system and types (symbols, numbers, strings, lists)
- **`edlisp/parser/`**: Multi-format parsing (line-based and JSON)
- **CLI implementation**: Uses Cobra framework for command-line interface
- **Root package**: Main texted package (per specification)

### Performance Design Decisions

- **Single-pass tokenization**: Character-by-character processing without backtracking
- **Immediate AST construction**: No separate parse tree intermediate representation
- **Float64 for all numbers**: No integer optimization
- **No type caching or interning**: Simple string-based type comparison

## Example Operations

- `search-forward "pattern"` - Search for text forward
- `set-mark` - Set the mark at current point
- `replace-region "text"` - Replace selected region
- `replace-match "text"` - Replace last search match

