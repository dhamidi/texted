package edlisp

import (
	"sort"
	"sync"
)

var (
	// documentationRegistry stores all function documentation
	documentationRegistry = make(map[string]FunctionDoc)

	// registryMutex protects concurrent access to the documentation registry
	registryMutex sync.RWMutex
)

// RegisterDocumentation adds a function's documentation to the global registry.
// This function is typically called from init() functions in builtin_*.go files.
func RegisterDocumentation(doc FunctionDoc) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	documentationRegistry[doc.Name] = doc
}

// GetDocumentation retrieves documentation for a function by name.
// Returns the documentation and a boolean indicating whether it was found.
func GetDocumentation(name string) (FunctionDoc, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	doc, exists := documentationRegistry[name]
	return doc, exists
}

// GetAllDocumentation returns all registered function documentation,
// sorted alphabetically by function name.
func GetAllDocumentation() []FunctionDoc {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	docs := make([]FunctionDoc, 0, len(documentationRegistry))
	for _, doc := range documentationRegistry {
		docs = append(docs, doc)
	}

	// Sort by function name for consistent output
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Name < docs[j].Name
	})

	return docs
}

// GetDocumentationByCategory returns all functions in a specific category,
// sorted alphabetically by function name.
func GetDocumentationByCategory(category string) []FunctionDoc {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	var docs []FunctionDoc
	for _, doc := range documentationRegistry {
		if doc.Category == category {
			docs = append(docs, doc)
		}
	}

	// Sort by function name for consistent output
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Name < docs[j].Name
	})

	return docs
}

// GetCategories returns all unique categories with registered functions,
// sorted alphabetically.
func GetCategories() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	categorySet := make(map[string]bool)
	for _, doc := range documentationRegistry {
		if doc.Category != "" {
			categorySet[doc.Category] = true
		}
	}

	categories := make([]string, 0, len(categorySet))
	for category := range categorySet {
		categories = append(categories, category)
	}

	sort.Strings(categories)
	return categories
}

// GetRegisteredFunctionCount returns the number of documented functions.
func GetRegisteredFunctionCount() int {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return len(documentationRegistry)
}
