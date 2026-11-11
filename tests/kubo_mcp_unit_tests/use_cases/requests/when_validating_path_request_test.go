package requests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenValidatingPathRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given the path is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PathRequest{Path: ""}

			assertValidationErrorOn(t, request.Validate(), "path")
		})
	})

	t.Run("Given a path that does not start with a CID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PathRequest{Path: "/ipfs/" + common_tests.RandomCid()}

			assertValidationErrorOn(t, request.Validate(), "path")
		})
	})

	t.Run("Given a CID followed by a sub-path", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PathRequest{Path: common_tests.RandomCid() + "/images/logo.png"}

			assert.NoError(t, request.Validate())
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.PathRequest{Path: common_tests.RandomCid()}

			assert.NoError(t, request.Validate())
		})
	})
}
