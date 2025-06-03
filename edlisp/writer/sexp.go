package writer

import (
	"fmt"
	"io"

	"github.com/dhamidi/texted/edlisp"
)

// SExpWriter writes edlisp values in S-expression format
type SExpWriter struct{}

// Write writes multiple expressions to the writer in S-expression format
func (w *SExpWriter) Write(writer io.Writer, expressions []edlisp.Value) error {
	for _, expr := range expressions {
		if err := w.WriteValue(writer, expr); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}
	return nil
}

// WriteValue writes a single value to the writer in S-expression format
func (w *SExpWriter) WriteValue(writer io.Writer, value edlisp.Value) error {
	_, err := fmt.Fprint(writer, value)
	return err
}
