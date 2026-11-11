package services_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenReadingAFileFromKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the cat request fails", func(t *testing.T) {
		t.Parallel()

		path := common_tests.RandomCid()
		requestErr := errors.New("connection refused")

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				CatWithResponse(anyCtx, &clients.CatParams{Arg: path, Length: lo.ToPtr(int64(10))}).
				Return(nil, requestErr).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Read(context.Background(), path, 10)

			assert.ErrorIs(t, err, requestErr)
		})
	})

	t.Run("Given the file does not exist", func(t *testing.T) {
		t.Parallel()

		path := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				CatWithResponse(anyCtx, &clients.CatParams{Arg: path, Length: lo.ToPtr(int64(10))}).
				Return(&clients.CatResponse{
					HTTPResponse: newHttpResponse(http.StatusInternalServerError),
					JSONDefault:  newKuboError("block was not found locally"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error carrying the Kubo message", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Read(context.Background(), path, 10)

			assert.EqualError(t, err, "kubo cat failed with status 500: block was not found locally")
		})
	})

	t.Run("Given the stream fails after the status was sent", func(t *testing.T) {
		t.Parallel()

		path := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			httpResponse := newHttpResponse(http.StatusOK)
			httpResponse.Trailer.Set("X-Stream-Error", "context deadline exceeded")
			suite.mockClient.EXPECT().
				CatWithResponse(anyCtx, &clients.CatParams{Arg: path, Length: lo.ToPtr(int64(10))}).
				Return(&clients.CatResponse{HTTPResponse: httpResponse, Body: []byte("hel")}, nil).
				Times(1)
		}

		t.Run("Should return the stream error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Read(context.Background(), path, 10)

			assert.EqualError(t, err, "kubo cat failed: context deadline exceeded")
		})
	})

	t.Run("Given the file exists", func(t *testing.T) {
		t.Parallel()

		path := common_tests.RandomCid() + "/logo.png"

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				CatWithResponse(anyCtx, &clients.CatParams{Arg: path, Length: lo.ToPtr(int64(10))}).
				Return(&clients.CatResponse{HTTPResponse: newHttpResponse(http.StatusOK), Body: []byte("hello")}, nil).
				Times(1)
		}

		t.Run("Should return the content", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			content, err := suite.sut.Read(context.Background(), path, 10)

			require.NoError(t, err)
			assert.Equal(t, "hello", string(content))
		})
	})
}
