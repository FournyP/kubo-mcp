package helpers_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_mcp/helpers"
)

func TestWhenBuildingATextResult(t *testing.T) {
	t.Parallel()

	t.Run("Given a payload that cannot be encoded", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a marshal error", func(t *testing.T) {
			t.Parallel()

			_, err := helpers.TextResult(make(chan int))

			assert.ErrorContains(t, err, "failed to marshal tool result")
		})
	})

	t.Run("Given a payload", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return it as one JSON text content", func(t *testing.T) {
			t.Parallel()

			result, err := helpers.TextResult(map[string]string{"cid": "bafy"})

			require.NoError(t, err)
			require.Len(t, result.Content, 1)
			assert.Equal(t, `{"cid":"bafy"}`, result.Content[0].(*mcp.TextContent).Text)
		})
	})
}
