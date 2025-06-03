package edlisp

import "fmt"

// NumberKind represents the kind for numeric values.
type NumberKind struct{}

// KindName returns the unique name for number kind.
func (kind *NumberKind) KindName() string {
	return "number"
}

// TheNumberKind is the singleton instance of NumberKind.
var TheNumberKind = &NumberKind{}

// Number represents a numeric value in texted expressions.
type Number struct {
	Value float64
}

// Kind returns the ValueKind for numbers.
func (num *Number) Kind() ValueKind {
	return TheNumberKind
}

// NewNumber creates a new Number with the given value.
func NewNumber(value float64) *Number {
	return &Number{Value: value}
}

// NewIntNumber creates a new Number from an integer value.
func NewIntNumber(value int) *Number {
	return &Number{Value: float64(value)}
}

// String returns the string representation of the number.
func (num *Number) String() string {
	if num.Value == float64(int(num.Value)) {
		return fmt.Sprintf("%.0f", num.Value)
	}
	return fmt.Sprintf("%g", num.Value)
}

// Int returns the integer representation of the number.
func (num *Number) Int() int {
	return int(num.Value)
}

// Float returns the float64 value of the number.
func (num *Number) Float() float64 {
	return num.Value
}
