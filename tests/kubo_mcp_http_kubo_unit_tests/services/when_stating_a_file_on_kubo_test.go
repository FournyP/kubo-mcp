package services_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenStatingAFileOnKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the stat request fails", func(t *testing.T) {
		t.Parallel()

		fileCid := common_tests.RandomCid()
		requestErr := errors.New("connection refused")

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				FilesStatWithResponse(anyCtx, &clients.FilesStatParams{Arg: "/ipfs/" + fileCid}).
				Return(nil, requestErr).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Stat(context.Background(), fileCid)

			assert.ErrorIs(t, err, requestErr)
		})
	})

	t.Run("Given Kubo fails", func(t *testing.T) {
		t.Parallel()

		fileCid := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				FilesStatWithResponse(anyCtx, &clients.FilesStatParams{Arg: "/ipfs/" + fileCid}).
				Return(&clients.FilesStatResponse{
					HTTPResponse: newHttpResponse(http.StatusInternalServerError),
					JSONDefault:  newKuboError("invalid path"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error carrying the Kubo message", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Stat(context.Background(), fileCid)

			assert.EqualError(t, err, "kubo files/stat failed with status 500: invalid path")
		})
	})

	t.Run("Given the file exists", func(t *testing.T) {
		t.Parallel()

		fileCid := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				FilesStatWithResponse(anyCtx, &clients.FilesStatParams{Arg: "/ipfs/" + fileCid}).
				Return(&clients.FilesStatResponse{
					HTTPResponse: newHttpResponse(http.StatusOK),
					JSON200: &clients.FilesStatResult{
						Hash:           fileCid,
						Size:           10,
						CumulativeSize: 12,
						Blocks:         1,
						Type:           "file",
					},
				}, nil).
				Times(1)
		}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			fileInfo, err := suite.sut.Stat(context.Background(), fileCid)

			require.NoError(t, err)
			assert.Equal(t, fileCid, fileInfo.Cid)
			assert.Equal(t, "file", fileInfo.Type)
			assert.Equal(t, int64(10), fileInfo.SizeBytes)
			assert.Equal(t, int64(12), fileInfo.CumulativeSizeBytes)
			assert.Equal(t, int64(1), fileInfo.Blocks)
		})
	})
}
