package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenUnpinningAFileOnKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the file is not pinned", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a Kubo error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			_, storedFile := suite.addRandomFile(t, 64)
			require.NoError(t, suite.sut.Unpin(context.Background(), storedFile.Cid))

			err := suite.sut.Unpin(context.Background(), storedFile.Cid)

			assert.ErrorContains(t, err, "not pinned")
		})
	})

	t.Run("Given the file is pinned", func(t *testing.T) {
		t.Parallel()

		t.Run("Should remove it from the pins", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			_, storedFile := suite.addRandomFile(t, 64)

			err := suite.sut.Unpin(context.Background(), storedFile.Cid)

			require.NoError(t, err)
			pins, err := suite.sut.ListPins(context.Background())
			require.NoError(t, err)
			assert.NotContains(t, pins, recursivePin(storedFile.Cid))
		})
	})
}
