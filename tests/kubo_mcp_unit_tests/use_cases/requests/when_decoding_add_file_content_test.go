package requests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
)

func TestWhenDecodingAddFileContent(t *testing.T) {
	t.Parallel()

	generateAddFileRequest := func(contentBase64 string) *requests.AddFileRequest {
		return &requests.AddFileRequest{
			FileNameRequest: requests.FileNameRequest{FileName: "hello.txt"},
			ContentBase64:   contentBase64,
		}
	}

	t.Run("Given padded base64", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the decoded content", func(t *testing.T) {
			t.Parallel()

			content, err := generateAddFileRequest("aGVsbG8=").DecodeContent()

			require.NoError(t, err)
			assert.Equal(t, "hello", string(content))
		})
	})

	t.Run("Given unpadded base64", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the decoded content", func(t *testing.T) {
			t.Parallel()

			content, err := generateAddFileRequest("aGVsbG8").DecodeContent()

			require.NoError(t, err)
			assert.Equal(t, "hello", string(content))
		})
	})

	t.Run("Given a data URI", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the decoded content", func(t *testing.T) {
			t.Parallel()

			content, err := generateAddFileRequest("data:text/plain;base64,aGVsbG8=").DecodeContent()

			require.NoError(t, err)
			assert.Equal(t, "hello", string(content))
		})
	})

	t.Run("Given content that is not base64", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			_, err := generateAddFileRequest("not base64 !").DecodeContent()

			assertValidationErrorOn(t, err, "contentBase64")
		})
	})
}
