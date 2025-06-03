// Package edlisp implements the value types and kinds for the texted editor's Lisp-like expressions.
package edlisp

// ValueKind represents a type classification for values.
type ValueKind interface {
	// KindName returns a unique name identifying this kind.
	KindName() string
}

// Value represents any value that can be used in texted expressions.
type Value interface {
	Kind() ValueKind
}

// IsA checks if a value is of a specific kind.
func IsA(value Value, kind ValueKind) bool {
	return value.Kind().KindName() == kind.KindName()
}