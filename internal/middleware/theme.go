package middleware

import (
	"context"
	"net/http"
	"strings"
)

type themeContextKey struct{}

var themeKey = themeContextKey{}

const (
	ThemeLight = "light"
	ThemeDark  = "dark"
)

func ThemeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		theme := ThemeLight
		if cookie, err := r.Cookie("theme"); err == nil {
			val := strings.ToLower(strings.TrimSpace(cookie.Value))
			if val == ThemeDark {
				theme = ThemeDark
			} else if val == ThemeLight {
				theme = ThemeLight
			}
		}
		ctx := context.WithValue(r.Context(), themeKey, theme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ThemeFromContext(ctx context.Context) string {
	if ctx == nil {
		return ThemeLight
	}
	if theme, ok := ctx.Value(themeKey).(string); ok && (theme == ThemeDark || theme == ThemeLight) {
		return theme
	}
	return ThemeLight
}
