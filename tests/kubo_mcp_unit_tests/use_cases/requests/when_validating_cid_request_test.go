package requests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenValidatingCidRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given the CID is missing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.CidRequest{Cid: ""}

			assertValidationErrorOn(t, request.Validate(), "cid")
		})
	})

	t.Run("Given an invalid CID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.CidRequest{Cid: "bafy-not-a-cid"}

			assertValidationErrorOn(t, request.Validate(), "cid")
		})
	})

	t.Run("Given a CID followed by a sub-path", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			request := &requests.CidRequest{Cid: common_tests.RandomCid() + "/logo.png"}

			assertValidationErrorOn(t, request.Validate(), "cid")
		})
	})

	t.Run("Given a CIDv0", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.CidRequest{Cid: "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"}

			assert.NoError(t, request.Validate())
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.CidRequest{Cid: common_tests.RandomCid()}

			assert.NoError(t, request.Validate())
		})
	})
}
