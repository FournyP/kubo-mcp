package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	common_errors "github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
	"github.com/FournyP/kubo-mcp/tests/kubo_mcp_mocks/services_mocks/interfaces_mocks"
)

const testMaxReadSizeBytes = int64(64)

var pngHeader = []byte("\x89PNG\r\n\x1a\n")

type WhenGettingAFileTestingSuite struct {
	sut *use_cases.GetFile

	mockIpfsStorage *interfaces_mocks.MockIpfsStorageInterface
}

func WhenGettingAFileBeforeEach(t *testing.T) *WhenGettingAFileTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)

	return &WhenGettingAFileTestingSuite{
		sut:             use_cases.NewGetFile(mockIpfsStorage, newTestIpfsSettings(0, testMaxReadSizeBytes)),
		mockIpfsStorage: mockIpfsStorage,
	}
}

func newGetFileRequest(path string) *requests.GetFileRequest {
	return &requests.GetFileRequest{PathRequest: requests.PathRequest{Path: path}}
}

func TestWhenGettingAFile(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid request", func(t *testing.T) {
		t.Parallel()

		request := newGetFileRequest("not-a-cid")

		initSuite := func(suite *WhenGettingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "path", validationError.Field)
		})
	})

	t.Run("Given reading the file fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		path := common_tests.RandomCid()
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenGettingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Read(ctx, path, testMaxReadSizeBytes+1).Return(nil, storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(ctx, newGetFileRequest(path))

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given a file larger than the read limit", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		path := common_tests.RandomCid()

		initSuite := func(suite *WhenGettingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().
				Read(ctx, path, testMaxReadSizeBytes+1).
				Return(make([]byte, testMaxReadSizeBytes+1), nil).
				Times(1)
		}

		t.Run("Should return a file too large error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(ctx, newGetFileRequest(path))

			var fileTooLargeError *common_errors.FileTooLargeError
			require.ErrorAs(t, err, &fileTooLargeError)
			assert.Equal(t, testMaxReadSizeBytes, fileTooLargeError.MaxSizeBytes)
		})
	})

	t.Run("Given a path without an extension", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		path := common_tests.RandomCid()

		initSuite := func(suite *WhenGettingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Read(ctx, path, testMaxReadSizeBytes+1).Return(pngHeader, nil).Times(1)
		}

		t.Run("Should return the content with its sniffed MIME type", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAFileBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(ctx, newGetFileRequest(path))

			require.NoError(t, err)
			assert.Equal(t, path, response.Path)
			assert.Equal(t, "image/png", response.MimeType)
			assert.Equal(t, int64(len(pngHeader)), response.SizeBytes)
			assert.Equal(t, pngHeader, response.Content)
		})
	})

	t.Run("Given a path with a known extension", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		path := common_tests.RandomCid() + "/logo.svg"
		content := []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`)

		initSuite := func(suite *WhenGettingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Read(ctx, path, testMaxReadSizeBytes+1).Return(content, nil).Times(1)
		}

		t.Run("Should use the MIME type of the extension", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAFileBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Get(ctx, newGetFileRequest(path))

			require.NoError(t, err)
			assert.Equal(t, "image/svg+xml", response.MimeType)
		})
	})
}
