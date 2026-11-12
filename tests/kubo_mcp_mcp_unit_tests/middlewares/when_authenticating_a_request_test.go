package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/middlewares"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/settings"
)

type WhenAuthenticatingARequestTestingSuite struct {
	sut http.Handler
}

func WhenAuthenticatingARequestBeforeEach(apiKeys []string) *WhenAuthenticatingARequestTestingSuite {
	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	middleware := middlewares.NewBearerAuthMiddleware(&settings.AuthSettings{ApiKeys: apiKeys})

	return &WhenAuthenticatingARequestTestingSuite{sut: middleware.Wrap(next)}
}

func (suite *WhenAuthenticatingARequestTestingSuite) serve(authorization string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}
	recorder := httptest.NewRecorder()
	suite.sut.ServeHTTP(recorder, request)
	return recorder
}

func TestWhenAuthenticatingARequest(t *testing.T) {
	t.Parallel()

	apiKeys := []string{"first-key", "second-key"}

	t.Run("Given no API key is configured", func(t *testing.T) {
		t.Parallel()

		t.Run("Should let the request through", func(t *testing.T) {
			t.Parallel()

			suite := WhenAuthenticatingARequestBeforeEach(nil)

			assert.Equal(t, http.StatusOK, suite.serve("").Code)
		})
	})

	t.Run("Given the Authorization header is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized bearer challenge", func(t *testing.T) {
			t.Parallel()

			suite := WhenAuthenticatingARequestBeforeEach(apiKeys)

			response := suite.serve("")

			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, "Bearer", response.Header().Get("WWW-Authenticate"))
		})
	})

	t.Run("Given a scheme other than Bearer", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized status", func(t *testing.T) {
			t.Parallel()

			suite := WhenAuthenticatingARequestBeforeEach(apiKeys)

			assert.Equal(t, http.StatusUnauthorized, suite.serve("Basic first-key").Code)
		})
	})

	t.Run("Given an unknown key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized status", func(t *testing.T) {
			t.Parallel()

			suite := WhenAuthenticatingARequestBeforeEach(apiKeys)

			assert.Equal(t, http.StatusUnauthorized, suite.serve("Bearer first-ke").Code)
		})
	})

	t.Run("Given any of the configured keys", func(t *testing.T) {
		t.Parallel()

		t.Run("Should let the request through", func(t *testing.T) {
			t.Parallel()

			suite := WhenAuthenticatingARequestBeforeEach(apiKeys)

			assert.Equal(t, http.StatusOK, suite.serve("Bearer second-key").Code)
		})
	})
}
