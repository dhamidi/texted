package writer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dhamidi/texted/edlisp"
)

// JSONWriter writes edlisp values in JSON format
type JSONWriter struct{}

// Write writes multiple expressions to the writer in JSON format
func (w *JSONWriter) Write(writer io.Writer, expressions []edlisp.Value) error {
	encoder := json.NewEncoder(writer)
	for _, expr := range expressions {
		jsonValue, err := w.valueToJSON(expr)
		if err != nil {
			return err
		}
		if err := encoder.Encode(jsonValue); err != nil {
			return err
		}
	}
	return nil
}

// WriteValue writes a single value to the writer in JSON format
func (w *JSONWriter) WriteValue(writer io.Writer, value edlisp.Value) error {
	jsonValue, err := w.valueToJSON(value)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	return encoder.Encode(jsonValue)
}

// valueToJSON converts an edlisp value to a JSON-compatible value
func (w *JSONWriter) valueToJSON(value edlisp.Value) (interface{}, error) {
	switch v := value.(type) {
	case *edlisp.List:
		jsonArray := make([]interface{}, 0, len(v.Elements))
		for _, element := range v.Elements {
			jsonValue, err := w.valueToJSON(element)
			if err != nil {
				return nil, err
			}
			jsonArray = append(jsonArray, jsonValue)
		}
		return jsonArray, nil
	case *edlisp.Symbol:
		return v.Name, nil
	case *edlisp.String:
		return v.Value, nil
	case *edlisp.Number:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("unsupported value type for JSON: %T", value)
	}
}
