package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/settings"
)

const bearerPrefix = "Bearer "

type BearerAuthMiddleware struct {
	apiKeyHashes [][sha256.Size]byte
}

func NewBearerAuthMiddleware(authSettings *settings.AuthSettings) *BearerAuthMiddleware {
	apiKeyHashes := make([][sha256.Size]byte, 0, len(authSettings.ApiKeys))
	for _, apiKey := range authSettings.ApiKeys {
		apiKeyHashes = append(apiKeyHashes, sha256.Sum256([]byte(apiKey)))
	}
	return &BearerAuthMiddleware{apiKeyHashes: apiKeyHashes}
}

// Wrap passes every request through when no key is configured.
func (m *BearerAuthMiddleware) Wrap(next http.Handler) http.Handler {
	if len(m.apiKeyHashes) == 0 {
		return next
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !m.isAuthorized(request.Header.Get("Authorization")) {
			writer.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

// isAuthorized compares hashes in constant time, so neither key content nor length leaks.
func (m *BearerAuthMiddleware) isAuthorized(authorization string) bool {
	token, found := strings.CutPrefix(authorization, bearerPrefix)
	if !found || token == "" {
		return false
	}

	tokenHash := sha256.Sum256([]byte(token))
	authorized := 0
	for _, apiKeyHash := range m.apiKeyHashes {
		authorized |= subtle.ConstantTimeCompare(tokenHash[:], apiKeyHash[:])
	}
	return authorized == 1
}
