package middleware

import (
	"context"
	"github.com/leonelquinteros/gotext"
	"golang.org/x/text/language"
	"net/http"
)

var langMatcher = language.NewMatcher([]language.Tag{
	language.English,
	language.German,
})

const CurrentLanguageKey = "middleware.lang.current"

func WithLocalization(next http.Handler) http.Handler {
	gotext.Configure("web/locales", "en", "default")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang, _ := r.Cookie("user-language")
		acceptLang := r.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(langMatcher, lang.String(), acceptLang)
		ctx := context.WithValue(r.Context(), CurrentLanguageKey, tag.String())
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
