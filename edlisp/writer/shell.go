package writer

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/dhamidi/texted/edlisp"
)

// ShellWriter writes edlisp values in shell-like format
type ShellWriter struct{}

// Write writes multiple expressions to the writer in shell format
func (w *ShellWriter) Write(writer io.Writer, expressions []edlisp.Value) error {
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

// WriteValue writes a single value to the writer in shell format
func (w *ShellWriter) WriteValue(writer io.Writer, value edlisp.Value) error {
	shellStr, err := w.valueToShellString(value)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(writer, shellStr)
	return err
}

// valueToShellString converts an edlisp value to shell-like string format
func (w *ShellWriter) valueToShellString(value edlisp.Value) (string, error) {
	list, ok := value.(*edlisp.List)
	if !ok {
		return "", fmt.Errorf("can only convert lists to shell format, got %T", value)
	}

	if list.IsEmpty() {
		return "", nil
	}

	var parts []string
	for _, element := range list.Elements {
		part, err := w.valueToToken(element)
		if err != nil {
			return "", err
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, " "), nil
}

// valueToToken converts a single value to a shell token
func (w *ShellWriter) valueToToken(value edlisp.Value) (string, error) {
	switch v := value.(type) {
	case *edlisp.Symbol:
		return v.Name, nil
	case *edlisp.String:
		return strconv.Quote(v.Value), nil
	case *edlisp.Number:
		return fmt.Sprintf("%v", v), nil
	case *edlisp.List:
		return "", fmt.Errorf("nested lists are not supported in shell format")
	default:
		return "", fmt.Errorf("unsupported value type for shell format: %T", value)
	}
}