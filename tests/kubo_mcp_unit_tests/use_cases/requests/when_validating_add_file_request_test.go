package requests_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
)

func TestWhenValidatingAddFileRequest(t *testing.T) {
	t.Parallel()

	generateAddFileRequest := func() *requests.AddFileRequest {
		return &requests.AddFileRequest{
			FileNameRequest: requests.FileNameRequest{FileName: "logo.png"},
			ContentBase64:   "aGVsbG8=",
		}
	}

	t.Run("Given the file name is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileRequest()
			request.FileName = " "

			assertValidationErrorOn(t, request.Validate(), "fileName")
		})
	})

	t.Run("Given the file name is longer than 255 characters", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileRequest()
			request.FileName = strings.Repeat("a", 256)

			assertValidationErrorOn(t, request.Validate(), "fileName")
		})
	})

	t.Run("Given the file name contains a path separator", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileRequest()
			request.FileName = "../logo.png"

			assertValidationErrorOn(t, request.Validate(), "fileName")
		})
	})

	t.Run("Given the content is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := generateAddFileRequest()
			request.ContentBase64 = ""

			assertValidationErrorOn(t, request.Validate(), "contentBase64")
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			assert.NoError(t, generateAddFileRequest().Validate())
		})
	})
}
