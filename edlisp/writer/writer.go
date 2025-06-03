package writer

import (
	"fmt"
	"io"

	"github.com/dhamidi/texted/edlisp"
)

// Writer defines the interface for format writers
type Writer interface {
	Write(w io.Writer, expressions []edlisp.Value) error
	WriteValue(w io.Writer, value edlisp.Value) error
}

// Format represents supported output formats
type Format string

const (
	FormatShell Format = "shell"
	FormatSExp  Format = "sexp"
	FormatJSON  Format = "json"
)

// NewWriter creates a writer for the specified format
func NewWriter(format Format) (Writer, error) {
	switch format {
	case FormatShell:
		return &ShellWriter{}, nil
	case FormatSExp:
		return &SExpWriter{}, nil
	case FormatJSON:
		return &JSONWriter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}