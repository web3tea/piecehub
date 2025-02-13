package api

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	tokens  map[string]struct{}
	enabled bool
}

func NewAuthenticator(tokens []string) *Authenticator {
	if len(tokens) == 0 {
		return &Authenticator{enabled: false}
	}
	tokenMap := make(map[string]struct{})
	for _, token := range tokens {
		tokenMap[token] = struct{}{}
	}
	return &Authenticator{tokens: tokenMap, enabled: true}
}

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.enabled {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		token := extractToken(auth)
		if token == "" {
			http.Error(w, "Unauthorized - Invalid token format", http.StatusUnauthorized)
			return
		}

		if _, ok := a.tokens[token]; !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractToken(auth string) string {
	// handle "Authorization: your-token" format
	if !strings.Contains(auth, " ") {
		return auth
	}

	// handle "Authorization: Bearer/Basic/Token your-token" format
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 {
		return auth // invalid format
	}

	switch strings.ToLower(parts[0]) {
	case "bearer", "token", "basic":
		return parts[1]
	default:
		return parts[1]
	}
}
