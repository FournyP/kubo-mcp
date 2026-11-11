package use_cases_test

import (
	"context"
	"errors"
	"io"
	"strings"
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

const testRemoteUrl = "https://example.com/images/cat.png"

type WhenAddingAFileFromAUrlTestingSuite struct {
	sut *use_cases.AddFileFromUrl

	mockIpfsStorage          *interfaces_mocks.MockIpfsStorageInterface
	mockRemoteFileDownloader *interfaces_mocks.MockRemoteFileDownloaderInterface
}

func WhenAddingAFileFromAUrlBeforeEach(t *testing.T, maxUploadSizeBytes int64) *WhenAddingAFileFromAUrlTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)
	mockRemoteFileDownloader := interfaces_mocks.NewMockRemoteFileDownloaderInterface(mockController)

	return &WhenAddingAFileFromAUrlTestingSuite{
		sut: use_cases.NewAddFileFromUrl(
			mockIpfsStorage,
			mockRemoteFileDownloader,
			newTestIpfsSettings(maxUploadSizeBytes, 0),
		),
		mockIpfsStorage:          mockIpfsStorage,
		mockRemoteFileDownloader: mockRemoteFileDownloader,
	}
}

func newAddFileFromUrlRequest(rawUrl, fileName string) *requests.AddFileFromUrlRequest {
	return &requests.AddFileFromUrlRequest{
		FileNameRequest: requests.FileNameRequest{FileName: fileName},
		Url:             rawUrl,
	}
}

func newRemoteFile(name, content string, sizeBytes int64) *models.RemoteFile {
	return &models.RemoteFile{
		Name:        name,
		ContentType: "image/png",
		SizeBytes:   sizeBytes,
		Body:        io.NopCloser(strings.NewReader(content)),
	}
}

func TestWhenAddingAFileFromAUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid request", func(t *testing.T) {
		t.Parallel()

		request := newAddFileFromUrlRequest("ftp://example.com/cat.png", "")

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().Download(gomock.Any(), gomock.Any()).Times(0)
			suite.mockIpfsStorage.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 1024)
			initSuite(suite)

			_, err := suite.sut.Add(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "url", validationError.Field)
		})
	})

	t.Run("Given downloading the file fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "")
		downloadErr := errors.New("connection refused")

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().Download(ctx, testRemoteUrl).Return(nil, downloadErr).Times(1)
			suite.mockIpfsStorage.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 1024)
			initSuite(suite)

			_, err := suite.sut.Add(ctx, request)

			assert.ErrorIs(t, err, downloadErr)
		})
	})

	t.Run("Given an announced size above the upload limit", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "")

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().
				Download(ctx, testRemoteUrl).
				Return(newRemoteFile("cat.png", "hello", 5), nil).
				Times(1)
			suite.mockIpfsStorage.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a file too large error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 4)
			initSuite(suite)

			_, err := suite.sut.Add(ctx, request)

			var fileTooLargeError *common_errors.FileTooLargeError
			require.ErrorAs(t, err, &fileTooLargeError)
			assert.Equal(t, int64(4), fileTooLargeError.MaxSizeBytes)
		})
	})

	t.Run("Given an unannounced size above the upload limit", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "")
		uploadErr := errors.New("request body read failed")

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().
				Download(ctx, testRemoteUrl).
				Return(newRemoteFile("cat.png", "hello", -1), nil).
				Times(1)
			suite.mockIpfsStorage.EXPECT().
				Add(ctx, "cat.png", gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, content io.Reader) (*models.StoredFile, error) {
					_, _ = io.ReadAll(content)
					return nil, uploadErr
				}).
				Times(1)
		}

		t.Run("Should return a file too large error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 4)
			initSuite(suite)

			_, err := suite.sut.Add(ctx, request)

			var fileTooLargeError *common_errors.FileTooLargeError
			require.ErrorAs(t, err, &fileTooLargeError)
			assert.Equal(t, int64(4), fileTooLargeError.MaxSizeBytes)
		})
	})

	t.Run("Given adding the file to IPFS fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "")
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().
				Download(ctx, testRemoteUrl).
				Return(newRemoteFile("cat.png", "hello", 5), nil).
				Times(1)
			suite.mockIpfsStorage.EXPECT().Add(ctx, "cat.png", gomock.Any()).Return(nil, storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 1024)
			initSuite(suite)

			_, err := suite.sut.Add(ctx, request)

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given a file name is provided", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "renamed.png")
		fileCid := common_tests.RandomCid()

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().
				Download(ctx, testRemoteUrl).
				Return(newRemoteFile("cat.png", "hello", 5), nil).
				Times(1)
			suite.mockIpfsStorage.EXPECT().
				Add(ctx, "renamed.png", gomock.Any()).
				Return(&models.StoredFile{Cid: fileCid, Name: "renamed.png", SizeBytes: 5}, nil).
				Times(1)
		}

		t.Run("Should use the provided file name", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 1024)
			initSuite(suite)

			response, err := suite.sut.Add(ctx, request)

			require.NoError(t, err)
			assert.Equal(t, "renamed.png", response.Name)
		})
	})

	t.Run("Given no file name is provided", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "")
		fileCid := common_tests.RandomCid()

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().
				Download(ctx, testRemoteUrl).
				Return(newRemoteFile("cat.png", "hello", 5), nil).
				Times(1)
			suite.mockIpfsStorage.EXPECT().
				Add(ctx, "cat.png", gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, content io.Reader) (*models.StoredFile, error) {
					data, err := io.ReadAll(content)
					assert.NoError(t, err)
					assert.Equal(t, "hello", string(data))
					return &models.StoredFile{Cid: fileCid, Name: "cat.png", SizeBytes: 16}, nil
				}).
				Times(1)
		}

		t.Run(
			"Should upload the downloaded content under the remote file name with its content size",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenAddingAFileFromAUrlBeforeEach(t, 1024)
				initSuite(suite)

				response, err := suite.sut.Add(ctx, request)

				require.NoError(t, err)
				assert.Equal(t, fileCid, response.Cid)
				assert.Equal(t, int64(5), response.SizeBytes)
				assert.Equal(t, "ipfs://"+fileCid, response.IpfsUri)
				assert.Equal(t, testGatewayUrl+"/ipfs/"+fileCid+"?filename=cat.png", response.GatewayUrl)
			},
		)
	})

	t.Run("Given neither the request nor the remote file has a name", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		request := newAddFileFromUrlRequest(testRemoteUrl, "")
		fileCid := common_tests.RandomCid()

		initSuite := func(suite *WhenAddingAFileFromAUrlTestingSuite) {
			suite.mockRemoteFileDownloader.EXPECT().
				Download(ctx, testRemoteUrl).
				Return(newRemoteFile("", "hello", 5), nil).
				Times(1)
			suite.mockIpfsStorage.EXPECT().
				Add(ctx, "file", gomock.Any()).
				Return(&models.StoredFile{Cid: fileCid, Name: "file", SizeBytes: 5}, nil).
				Times(1)
		}

		t.Run("Should fall back to the default file name", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAFileFromAUrlBeforeEach(t, 1024)
			initSuite(suite)

			response, err := suite.sut.Add(ctx, request)

			require.NoError(t, err)
			assert.Equal(t, "file", response.Name)
		})
	})
}
