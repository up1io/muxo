// Package locales provides utilities for working with localization files.
package locales

// Reader is an interface for reading localized text.
type Reader interface {
	// Text returns the localized version of the given string.
	Text(s string, vars ...interface{}) string
}
