package edlisp

// SymbolKind represents the kind for symbol values.
type SymbolKind struct{}

// KindName returns the unique name for symbol kind.
func (kind *SymbolKind) KindName() string {
	return "symbol"
}

// TheSymbolKind is the singleton instance of SymbolKind.
var TheSymbolKind = &SymbolKind{}

// Symbol represents a symbolic name in texted expressions.
type Symbol struct {
	Name string
}

// Kind returns the ValueKind for symbols.
func (sym *Symbol) Kind() ValueKind {
	return TheSymbolKind
}

// NewSymbol creates a new Symbol with the given name.
func NewSymbol(name string) *Symbol {
	return &Symbol{Name: name}
}

// String returns the string representation of the symbol.
func (sym *Symbol) String() string {
	return sym.Name
}
