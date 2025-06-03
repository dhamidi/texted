package edlisp

// StringKind represents the kind for string values.
type StringKind struct{}

// KindName returns the unique name for string kind.
func (kind *StringKind) KindName() string {
	return "string"
}

// TheStringKind is the singleton instance of StringKind.
var TheStringKind = &StringKind{}

// String represents a string value in texted expressions.
type String struct {
	Value string
}

// Kind returns the ValueKind for strings.
func (str *String) Kind() ValueKind {
	return TheStringKind
}

// NewString creates a new String with the given value.
func NewString(value string) *String {
	return &String{Value: value}
}

// String returns the string representation of the string value.
func (str *String) String() string {
	return str.Value
}