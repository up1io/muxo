// Package middleware provides localization middleware for HTTP servers.
package middleware

import (
	"context"
	"github.com/leonelquinteros/gotext"
	"github.com/up1io/muxo/logger"
	"github.com/up1io/muxo/middleware"
	"golang.org/x/text/language"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Define supported languages
var supportedLanguages = []language.Tag{
	language.English,
	language.German,
	language.French,
	language.Spanish,
	language.Italian,
	language.Chinese,
	language.Japanese,
}

var langMatcher = language.NewMatcher(supportedLanguages)

// contextKey is a custom type to avoid collisions in the context values.
type contextKey string

// LanguageKey is the key used to store the language in the request context.
const LanguageKey contextKey = "current-language"

// LanguageFromContext returns the language value stored in ctx, if any.
func LanguageFromContext(ctx context.Context) (string, bool) {
	lang, ok := ctx.Value(LanguageKey).(string)
	return lang, ok
}

// NewLanguageContext returns a new Context that carries the language value.
func NewLanguageContext(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, LanguageKey, language)
}

// WithLocalization creates a middleware that configures localization based on the provided locales directory.
// It scans the locales directory for available locales and configures gotext with the default language.
func WithLocalization(localesDir string) middleware.Middleware {
	var availableLocales []string

	// Scan the locales directory for available locales
	if entries, err := os.ReadDir(localesDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				availableLocales = append(availableLocales, entry.Name())
			}
		}
	}

	if len(availableLocales) == 0 {
		availableLocales = []string{"en"}
	}

	defaultLang := availableLocales[0]
	domain := "default"

	gotext.Configure(localesDir, defaultLang, domain)

	logger.Info("Available locales: %s", strings.Join(availableLocales, ","))

	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var langStr string

			langCookie, err := r.Cookie("user-language")
			if err == nil && langCookie != nil {
				langStr = langCookie.Value
			}

			if langStr == "" {
				langStr = r.Header.Get("Accept-Language")
			}

			tag, _ := language.MatchStrings(langMatcher, langStr)

			ctx := NewLanguageContext(r.Context(), tag.String())
			req := r.WithContext(ctx)

			mu.Lock()
			gotext.SetLanguage(tag.String())
			mu.Unlock()

			next.ServeHTTP(w, req)
		})
	}
}
