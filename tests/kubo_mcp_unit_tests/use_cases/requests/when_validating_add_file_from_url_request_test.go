package requests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
)

func TestWhenValidatingAddFileFromUrlRequest(t *testing.T) {
	t.Parallel()

	generateAddFileFromUrlRequest := func() *requests.AddFileFromUrlRequest {
		return &requests.AddFileFromUrlRequest{Url: "https://example.com/cat.png"}
	}

	t.Run("Given the URL is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileFromUrlRequest()
			request.Url = ""

			assertValidationErrorOn(t, request.Validate(), "url")
		})
	})

	t.Run("Given a relative URL", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileFromUrlRequest()
			request.Url = "/cat.png"

			assertValidationErrorOn(t, request.Validate(), "url")
		})
	})

	t.Run("Given a URL with an unsupported scheme", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileFromUrlRequest()
			request.Url = "file:///etc/passwd"

			assertValidationErrorOn(t, request.Validate(), "url")
		})
	})

	t.Run("Given the file name contains a path separator", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileFromUrlRequest()
			request.FileName = "images/cat.png"

			assertValidationErrorOn(t, request.Validate(), "fileName")
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			assert.NoError(t, generateAddFileFromUrlRequest().Validate())
		})
	})
}
