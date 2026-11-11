package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenReadingAFileFromKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the file does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a Kubo error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)

			_, err := suite.sut.Read(context.Background(), common_tests.RandomCid(), 10)

			assert.ErrorContains(t, err, "kubo cat failed with status 500")
		})
	})

	t.Run("Given a sub-path inside a file", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a Kubo error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			_, storedFile := suite.addRandomFile(t, 64)

			_, err := suite.sut.Read(context.Background(), storedFile.Cid+"/nested.png", 10)

			assert.ErrorContains(t, err, "kubo cat failed with status 500")
		})
	})

	t.Run("Given a length shorter than the file", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return only the first bytes", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			content, storedFile := suite.addRandomFile(t, 64)

			read, err := suite.sut.Read(context.Background(), storedFile.Cid, 10)

			require.NoError(t, err)
			assert.Equal(t, content[:10], read)
		})
	})

	t.Run("Given the file exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the whole content", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			content, storedFile := suite.addRandomFile(t, 1<<20)

			read, err := suite.sut.Read(context.Background(), storedFile.Cid, 2<<20)

			require.NoError(t, err)
			assert.Equal(t, content, read)
		})
	})
}
