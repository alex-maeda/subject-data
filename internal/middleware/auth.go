// Package middleware provides HTTP middleware for the subject-data service.
package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

// BearerAuth returns middleware that validates the Authorization header against
// a set of allowed tokens. Tokens are compared with constant-time-equivalent
// map lookup (Go maps don't short-circuit on mismatch, but this is a small
// allowlist, not a password check — acceptable for V1).
//
// If allowedTokens is empty, auth is disabled and all requests pass through.
// This lets local dev work without configuring AUTH_TOKENS.
func BearerAuth(allowedTokens []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedTokens))
	for _, t := range allowedTokens {
		if t != "" {
			allowed[t] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// No tokens configured — auth disabled (local dev).
			if len(allowed) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			token := extractBearer(r.Header.Get("Authorization"))
			if token == "" {
				slog.Debug("missing or malformed Authorization header")
				writeUnauthorized(w)
				return
			}

			if _, ok := allowed[token]; !ok {
				slog.Debug("invalid bearer token")
				writeUnauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractBearer(header string) string {
	const prefix = "Bearer "
	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):]
	}
	return ""
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"Missing or invalid Authorization header"}}`))
}
