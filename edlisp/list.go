package edlisp

import (
	"fmt"
	"strings"
)

// ListKind represents the kind for list values.
type ListKind struct{}

// KindName returns the unique name for list kind.
func (kind *ListKind) KindName() string {
	return "list"
}

// TheListKind is the singleton instance of ListKind.
var TheListKind = &ListKind{}

// List represents a list of values in texted expressions.
type List struct {
	Elements []Value
}

// Kind returns the ValueKind for lists.
func (list *List) Kind() ValueKind {
	return TheListKind
}

// NewList creates a new List with the given elements.
func NewList(elements ...Value) *List {
	return &List{Elements: elements}
}

// NewEmptyList creates a new empty List.
func NewEmptyList() *List {
	return &List{Elements: make([]Value, 0)}
}

// String returns the string representation of the list.
func (list *List) String() string {
	if len(list.Elements) == 0 {
		return "()"
	}

	var parts []string
	for _, element := range list.Elements {
		parts = append(parts, fmt.Sprintf("%v", element))
	}
	return "(" + strings.Join(parts, " ") + ")"
}

// Len returns the number of elements in the list.
func (list *List) Len() int {
	return len(list.Elements)
}

// IsEmpty returns true if the list has no elements.
func (list *List) IsEmpty() bool {
	return len(list.Elements) == 0
}

// First returns the first element of the list, or nil if empty.
func (list *List) First() Value {
	if len(list.Elements) == 0 {
		return nil
	}
	return list.Elements[0]
}

// Rest returns a new list containing all elements except the first.
func (list *List) Rest() *List {
	if len(list.Elements) <= 1 {
		return NewEmptyList()
	}
	return &List{Elements: list.Elements[1:]}
}

// Append returns a new list with the given element appended.
func (list *List) Append(element Value) *List {
	newElements := make([]Value, len(list.Elements)+1)
	copy(newElements, list.Elements)
	newElements[len(list.Elements)] = element
	return &List{Elements: newElements}
}

// Get returns the element at the given index, or nil if out of bounds.
func (list *List) Get(index int) Value {
	if index < 0 || index >= len(list.Elements) {
		return nil
	}
	return list.Elements[index]
}
