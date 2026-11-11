package services_test

import (
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

var pinnedCidV1AddParams = &clients.AddParams{CidVersion: lo.ToPtr(1), Pin: lo.ToPtr(true)}

func TestWhenAddingAFileToKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the add request fails", func(t *testing.T) {
		t.Parallel()

		requestErr := errors.New("connection refused")

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				AddWithBodyWithResponse(anyCtx, pinnedCidV1AddParams, gomock.Any(), gomock.Any()).
				Return(nil, requestErr).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Add(context.Background(), "logo.png", strings.NewReader("hello"))

			assert.ErrorIs(t, err, requestErr)
		})
	})

	t.Run("Given Kubo rejects the upload", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				AddWithBodyWithResponse(anyCtx, pinnedCidV1AddParams, gomock.Any(), gomock.Any()).
				Return(&clients.AddResponse{
					HTTPResponse: newHttpResponse(http.StatusInternalServerError),
					JSONDefault:  newKuboError("repo is full"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error carrying the Kubo message", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Add(context.Background(), "logo.png", strings.NewReader("hello"))

			assert.EqualError(t, err, "kubo add failed with status 500: repo is full")
		})
	})

	t.Run("Given a response whose size cannot be decoded", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				AddWithBodyWithResponse(anyCtx, pinnedCidV1AddParams, gomock.Any(), gomock.Any()).
				Return(&clients.AddResponse{
					HTTPResponse: newHttpResponse(http.StatusOK),
					JSON200:      &clients.AddResult{Name: "logo.png", Hash: common_tests.RandomCid(), Size: "many"},
				}, nil).
				Times(1)
		}

		t.Run("Should return a parse error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Add(context.Background(), "logo.png", strings.NewReader("hello"))

			assert.ErrorContains(t, err, "failed to parse the kubo add size")
		})
	})

	t.Run("Given Kubo accepts the upload", func(t *testing.T) {
		t.Parallel()

		fileCid := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				AddWithBodyWithResponse(anyCtx, pinnedCidV1AddParams, gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, _ *clients.AddParams, contentType string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.AddResponse, error) {
					_, params, err := mime.ParseMediaType(contentType)
					require.NoError(t, err)
					part, err := multipart.NewReader(body, params["boundary"]).NextPart()
					require.NoError(t, err)
					content, err := io.ReadAll(part)
					require.NoError(t, err)

					assert.Equal(t, "file", part.FormName())
					assert.Equal(t, "logo.png", part.FileName())
					assert.Equal(t, "hello", string(content))

					return &clients.AddResponse{
						HTTPResponse: newHttpResponse(http.StatusOK),
						JSON200:      &clients.AddResult{Name: "logo.png", Hash: fileCid, Size: "13"},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should stream the file as a multipart upload and return the stored file", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			storedFile, err := suite.sut.Add(context.Background(), "logo.png", strings.NewReader("hello"))

			require.NoError(t, err)
			assert.Equal(t, fileCid, storedFile.Cid)
			assert.Equal(t, "logo.png", storedFile.Name)
			assert.Equal(t, int64(13), storedFile.SizeBytes)
		})
	})
}
