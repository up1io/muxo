// Package local provides localization utilities.
package local

import (
	"github.com/leonelquinteros/gotext"
	"net/http"
)

// Text returns the localized version of the given string.
// It uses the language set by the localization middleware.
func Text(s string) string {
	return gotext.Get(s)
}

// SetLocal sets the user's preferred language by setting a cookie.
// This allows users to switch the locale for future requests.
// The language code should be a valid language tag (e.g., "en", "de", "fr").
func SetLocal(w http.ResponseWriter, lang string) {
	cookie := &http.Cookie{
		Name:     "user-language",
		Value:    lang,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	gotext.SetLanguage(lang)
}
