package helpers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/helpers"
)

func TestWhenEditingARequestWithAnAuthorization(t *testing.T) {
	t.Parallel()

	t.Run("Given no authorization is configured", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not set an Authorization header", func(t *testing.T) {
			t.Parallel()

			request, err := http.NewRequest(http.MethodPost, "http://kubo:5001/api/v0/add", nil)
			require.NoError(t, err)

			err = helpers.NewAuthorizationRequestEditor("")(context.Background(), request)

			require.NoError(t, err)
			assert.Empty(t, request.Header.Values("Authorization"))
		})
	})

	t.Run("Given an authorization is configured", func(t *testing.T) {
		t.Parallel()

		t.Run("Should set it verbatim as the Authorization header", func(t *testing.T) {
			t.Parallel()

			request, err := http.NewRequest(http.MethodPost, "http://kubo:5001/api/v0/add", nil)
			require.NoError(t, err)

			err = helpers.NewAuthorizationRequestEditor("Basic dXNlcjpwYXNz")(context.Background(), request)

			require.NoError(t, err)
			assert.Equal(t, "Basic dXNlcjpwYXNz", request.Header.Get("Authorization"))
		})
	})
}
