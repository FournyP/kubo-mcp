package helpers_test

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/helpers"
	common_errors "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

func TestWhenReadingWithASizeLimit(t *testing.T) {
	t.Parallel()

	t.Run("Given content of exactly the limit", func(t *testing.T) {
		t.Parallel()

		t.Run("Should read the whole content", func(t *testing.T) {
			t.Parallel()

			reader := helpers.NewSizeLimitedReader(strings.NewReader("hello"), 5)

			content, err := io.ReadAll(reader)

			require.NoError(t, err)
			assert.Equal(t, "hello", string(content))
			assert.Equal(t, int64(5), reader.ReadBytes())
			assert.False(t, reader.Exceeded())
		})
	})

	t.Run("Given content larger than the limit", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a file too large error", func(t *testing.T) {
			t.Parallel()

			reader := helpers.NewSizeLimitedReader(strings.NewReader("hello"), 4)

			_, err := io.ReadAll(reader)

			var fileTooLargeError *common_errors.FileTooLargeError
			require.ErrorAs(t, err, &fileTooLargeError)
			assert.Equal(t, int64(4), fileTooLargeError.MaxSizeBytes)
			assert.True(t, reader.Exceeded())
		})
	})
}
