package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenAddingAFileToKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given a file smaller than one chunk", func(t *testing.T) {
		t.Parallel()

		t.Run("Should pin it under a raw CIDv1", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)

			content, storedFile := suite.addRandomFile(t, 1024)

			assert.Regexp(t, "^bafkrei", storedFile.Cid)
			assert.Equal(t, "file.bin", storedFile.Name)
			assert.Equal(t, int64(len(content)), storedFile.SizeBytes)
			pins, err := suite.sut.ListPins(context.Background())
			require.NoError(t, err)
			assert.Contains(t, pins, recursivePin(storedFile.Cid))
		})
	})

	t.Run("Given a file larger than one chunk", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store it as a DAG whose size includes the encoding overhead", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)

			content, storedFile := suite.addRandomFile(t, 1<<20)

			assert.Regexp(t, "^bafybei", storedFile.Cid)
			assert.Greater(t, storedFile.SizeBytes, int64(len(content)))
		})
	})
}
