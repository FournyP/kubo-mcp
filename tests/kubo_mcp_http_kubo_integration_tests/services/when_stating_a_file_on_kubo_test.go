package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenStatingAFileOnKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the file does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a Kubo error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)

			_, err := suite.sut.Stat(context.Background(), common_tests.RandomCid())

			assert.ErrorContains(t, err, "kubo files/stat failed with status 500")
		})
	})

	t.Run("Given the file exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should describe it", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			content, storedFile := suite.addRandomFile(t, 1<<20)

			fileInfo, err := suite.sut.Stat(context.Background(), storedFile.Cid)

			require.NoError(t, err)
			assert.Equal(t, storedFile.Cid, fileInfo.Cid)
			assert.Equal(t, "file", fileInfo.Type)
			assert.Equal(t, int64(len(content)), fileInfo.SizeBytes)
			assert.Equal(t, storedFile.SizeBytes, fileInfo.CumulativeSizeBytes)
			assert.Equal(t, int64(4), fileInfo.Blocks)
		})
	})
}
