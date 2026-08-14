package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThemeMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		cookie        *http.Cookie
		expectedTheme string
	}{
		{
			name:          "no cookie defaults to light",
			cookie:        nil,
			expectedTheme: ThemeLight,
		},
		{
			name:          "dark cookie sets dark",
			cookie:        &http.Cookie{Name: "theme", Value: "dark"},
			expectedTheme: ThemeDark,
		},
		{
			name:          "light cookie sets light",
			cookie:        &http.Cookie{Name: "theme", Value: "light"},
			expectedTheme: ThemeLight,
		},
		{
			name:          "case-insensitive and trimmed dark",
			cookie:        &http.Cookie{Name: "theme", Value: " DARK "},
			expectedTheme: ThemeDark,
		},
		{
			name:          "invalid value defaults to light",
			cookie:        &http.Cookie{Name: "theme", Value: "invalid-theme"},
			expectedTheme: ThemeLight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			var capturedTheme string
			nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				capturedTheme = ThemeFromContext(r.Context())
			})

			handler := ThemeMiddleware(nextHandler)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if capturedTheme != tt.expectedTheme {
				t.Errorf("expected theme %q, got %q", tt.expectedTheme, capturedTheme)
			}
		})
	}
}

func TestThemeFromContext_NilContext(t *testing.T) {
	theme := ThemeFromContext(nil)
	if theme != ThemeLight {
		t.Errorf("expected default theme %q for nil context, got %q", ThemeLight, theme)
	}
}

func TestThemeFromContext_EmptyContext(t *testing.T) {
	theme := ThemeFromContext(context.Background())
	if theme != ThemeLight {
		t.Errorf("expected default theme %q for empty context, got %q", ThemeLight, theme)
	}
}
