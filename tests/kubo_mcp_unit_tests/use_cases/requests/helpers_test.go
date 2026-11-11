package requests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common_errors "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

func assertValidationErrorOn(t *testing.T, err error, field string) {
	t.Helper()

	var validationError *common_errors.ValidationError
	require.ErrorAs(t, err, &validationError)
	assert.Equal(t, field, validationError.Field)
}
