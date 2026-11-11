package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	common_errors "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
	"github.com/FournyP/kubo-mcp/tests/kubo_mcp_mocks/services_mocks/interfaces_mocks"
)

type WhenGettingFileInfoTestingSuite struct {
	sut *use_cases.GetFileInfo

	mockIpfsStorage *interfaces_mocks.MockIpfsStorageInterface
}

func WhenGettingFileInfoBeforeEach(t *testing.T) *WhenGettingFileInfoTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)

	return &WhenGettingFileInfoTestingSuite{
		sut:             use_cases.NewGetFileInfo(mockIpfsStorage, newTestIpfsSettings(0, 0)),
		mockIpfsStorage: mockIpfsStorage,
	}
}

func newGetFileInfoRequest(path string) *requests.GetFileInfoRequest {
	return &requests.GetFileInfoRequest{PathRequest: requests.PathRequest{Path: path}}
}

func TestWhenGettingFileInfo(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid request", func(t *testing.T) {
		t.Parallel()

		request := newGetFileInfoRequest("")

		initSuite := func(suite *WhenGettingFileInfoTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Stat(gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingFileInfoBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "path", validationError.Field)
		})
	})

	t.Run("Given the file lookup fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		path := common_tests.RandomCid()
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenGettingFileInfoTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Stat(ctx, path).Return(nil, storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingFileInfoBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(ctx, newGetFileInfoRequest(path))

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		path := common_tests.RandomCid()
		fileInfo := &models.FileInfo{Cid: path, Type: "file", SizeBytes: 10, CumulativeSizeBytes: 12, Blocks: 1}

		initSuite := func(suite *WhenGettingFileInfoTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Stat(ctx, path).Return(fileInfo, nil).Times(1)
		}

		t.Run("Should populate every field", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingFileInfoBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(ctx, newGetFileInfoRequest(path))

			require.NoError(t, err)
			assert.Equal(t, path, response.Cid)
			assert.Equal(t, "file", response.Type)
			assert.Equal(t, int64(10), response.SizeBytes)
			assert.Equal(t, int64(12), response.CumulativeSizeBytes)
			assert.Equal(t, int64(1), response.Blocks)
			assert.Equal(t, testGatewayUrl+"/ipfs/"+path, response.GatewayUrl)
		})
	})
}
