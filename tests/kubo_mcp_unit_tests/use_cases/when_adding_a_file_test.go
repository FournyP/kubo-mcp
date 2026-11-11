package use_cases_test

import (
	"context"
	"errors"
	"io"
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

type WhenAddingAFileTestingSuite struct {
	sut *use_cases.AddFile

	mockIpfsStorage *interfaces_mocks.MockIpfsStorageInterface
}

func WhenAddingAFileBeforeEach(t *testing.T, maxUploadSizeBytes int64) *WhenAddingAFileTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)

	return &WhenAddingAFileTestingSuite{
		sut:             use_cases.NewAddFile(mockIpfsStorage, newTestIpfsSettings(maxUploadSizeBytes, 0)),
		mockIpfsStorage: mockIpfsStorage,
	}
}

func newAddFileRequest(fileName, contentBase64 string) *requests.AddFileRequest {
	return &requests.AddFileRequest{
		FileNameRequest: requests.FileNameRequest{FileName: fileName},
		ContentBase64:   contentBase64,
	}
}

func TestWhenAddingAFile(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid request", func(t *testing.T) {
		t.Parallel()

		request := newAddFileRequest("", "aGVsbG8=")

		initSuite := func(suite *WhenAddingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileBeforeEach(t, 1024)
			initSuite(suite)

			response, err := suite.sut.Add(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "fileName", validationError.Field)
			assert.Nil(t, response)
		})
	})

	t.Run("Given content that is not base64", func(t *testing.T) {
		t.Parallel()

		request := newAddFileRequest("hello.txt", "not base64 !")

		initSuite := func(suite *WhenAddingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileBeforeEach(t, 1024)
			initSuite(suite)

			_, err := suite.sut.Add(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "contentBase64", validationError.Field)
		})
	})

	t.Run("Given content larger than the upload limit", func(t *testing.T) {
		t.Parallel()

		// "hello" is five bytes, one past the limit.
		request := newAddFileRequest("hello.txt", "aGVsbG8=")

		initSuite := func(suite *WhenAddingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a file too large error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileBeforeEach(t, 4)
			initSuite(suite)

			_, err := suite.sut.Add(context.Background(), request)

			var fileTooLargeError *common_errors.FileTooLargeError
			require.ErrorAs(t, err, &fileTooLargeError)
			assert.Equal(t, int64(4), fileTooLargeError.MaxSizeBytes)
		})
	})

	t.Run("Given adding the file to IPFS fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileRequest("hello.txt", "aGVsbG8=")
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenAddingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Add(ctx, "hello.txt", gomock.Any()).Return(nil, storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileBeforeEach(t, 1024)
			initSuite(suite)

			_, err := suite.sut.Add(ctx, request)

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given a data URI", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileRequest("hello.txt", "data:text/plain;base64,aGVsbG8=")
		fileCid := common_tests.RandomCid()

		initSuite := func(suite *WhenAddingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().
				Add(ctx, "hello.txt", gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, content io.Reader) (*models.StoredFile, error) {
					data, err := io.ReadAll(content)
					assert.NoError(t, err)
					assert.Equal(t, "hello", string(data))
					return &models.StoredFile{Cid: fileCid, Name: "hello.txt", SizeBytes: 5}, nil
				}).
				Times(1)
		}

		t.Run("Should add the decoded content", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileBeforeEach(t, 1024)
			initSuite(suite)

			_, err := suite.sut.Add(ctx, request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileRequest("hello world.txt", "aGVsbG8=")
		fileCid := common_tests.RandomCid()

		initSuite := func(suite *WhenAddingAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().
				Add(ctx, "hello world.txt", gomock.Any()).
				Return(&models.StoredFile{Cid: fileCid, Name: "hello world.txt", SizeBytes: 16}, nil).
				Times(1)
		}

		t.Run("Should return the file links and its content size", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileBeforeEach(t, 1024)
			initSuite(suite)

			response, err := suite.sut.Add(ctx, request)

			require.NoError(t, err)
			assert.Equal(t, fileCid, response.Cid)
			assert.Equal(t, "hello world.txt", response.Name)
			assert.Equal(t, int64(5), response.SizeBytes)
			assert.Equal(t, "ipfs://"+fileCid, response.IpfsUri)
			assert.Equal(t, testGatewayUrl+"/ipfs/"+fileCid+"?filename=hello+world.txt", response.GatewayUrl)
		})
	})
}
