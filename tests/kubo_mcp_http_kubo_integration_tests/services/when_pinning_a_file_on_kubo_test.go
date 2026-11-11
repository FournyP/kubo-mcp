package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenPinningAFileOnKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given a local file that is not pinned", func(t *testing.T) {
		t.Parallel()

		t.Run("Should add it back to the pins", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			_, storedFile := suite.addRandomFile(t, 64)
			require.NoError(t, suite.sut.Unpin(context.Background(), storedFile.Cid))

			err := suite.sut.Pin(context.Background(), storedFile.Cid)

			require.NoError(t, err)
			pins, err := suite.sut.ListPins(context.Background())
			require.NoError(t, err)
			assert.Contains(t, pins, recursivePin(storedFile.Cid))
		})
	})
}
