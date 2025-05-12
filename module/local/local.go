// Package local provides localization utilities.
package local

import "github.com/leonelquinteros/gotext"

// Text returns the localized version of the given string.
// It uses the language set by the localization middleware.
func Text(s string) string {
	return gotext.Get(s)
}
