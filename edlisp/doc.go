// Package edlisp provides value types for texted's Lisp-like expression system.
//
// This package implements the core value types that can be used in texted scripts:
// symbols, numbers, strings, and lists. Each type implements the Value interface
// and has an associated ValueKind for type checking.
//
// Example usage:
//
//	// Create values
//	sym := edlisp.NewSymbol("search-forward")
//	str := edlisp.NewString("hello")
//	num := edlisp.NewNumber(42.0)
//	list := edlisp.NewList(sym, str)
//
//	// Type checking
//	if edlisp.IsA(sym, edlisp.TheSymbolKind) {
//		fmt.Println("It's a symbol:", sym.String())
//	}
//
//	// Working with lists
//	first := list.First()
//	rest := list.Rest()
//	extended := list.Append(num)
//
// The package follows the value system specification from the main texted
// documentation, providing a foundation for implementing the editor's
// script execution engine.
package edlisp