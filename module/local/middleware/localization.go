package middleware

import (
	"context"
	"github.com/leonelquinteros/gotext"
	"github.com/up1io/muxo/middleware"
	"golang.org/x/text/language"
	"log"
	"net/http"
	"os"
	"strings"
)

var langMatcher = language.NewMatcher([]language.Tag{
	language.English,
	language.German,
})

type key int

const CurrentLanguageKey key = 0

// LanguageFromContext returns the language value stored in ctx, if any.
func LanguageFromContext(ctx context.Context) (string, bool) {
	lang, ok := ctx.Value(CurrentLanguageKey).(string)
	return lang, ok
}

// NewLanguageContext returns a new Context that carries the language value.
func NewLanguageContext(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, CurrentLanguageKey, language)
}

// WithLocalization creates a middleware that configures localization based on the provided locales directory
func WithLocalization(localesDir string) middleware.Middleware {
	var availableLocales []string
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

	log.Println("available locales", strings.Join(availableLocales, ","))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang, _ := r.Cookie("user-language")
			acceptLang := r.Header.Get("Accept-Language")
			tag, _ := language.MatchStrings(langMatcher, lang.String(), acceptLang)

			ctx := NewLanguageContext(r.Context(), tag.String())
			req := r.WithContext(ctx)

			gotext.SetLanguage(tag.String())

			next.ServeHTTP(w, req)
		})
	}
}
