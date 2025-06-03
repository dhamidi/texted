package edlisp

// Equal performs a deep equality check between two edlisp.Value objects.
// It returns true if both values are of the same kind and their contents are equal.
func Equal(a, b Value) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Check if they are the same kind
	if a.Kind().KindName() != b.Kind().KindName() {
		return false
	}

	// Delegate to kind-specific equality checks
	switch {
	case IsA(a, TheStringKind):
		aStr := a.(*String)
		bStr := b.(*String)
		return aStr.Value == bStr.Value

	case IsA(a, TheNumberKind):
		aNum := a.(*Number)
		bNum := b.(*Number)
		return aNum.Value == bNum.Value

	case IsA(a, TheSymbolKind):
		aSym := a.(*Symbol)
		bSym := b.(*Symbol)
		return aSym.Name == bSym.Name

	case IsA(a, TheListKind):
		aList := a.(*List)
		bList := b.(*List)

		// Check if lengths are different
		if aList.Len() != bList.Len() {
			return false
		}

		// Recursively check each element
		for i := 0; i < aList.Len(); i++ {
			if !Equal(aList.Get(i), bList.Get(i)) {
				return false
			}
		}
		return true

	default:
		// For unknown types, fall back to interface equality
		return a == b
	}
}
