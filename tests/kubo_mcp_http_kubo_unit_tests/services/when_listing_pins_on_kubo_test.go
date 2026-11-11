package services_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
)

var recursivePinLsParams = &clients.PinLsParams{Type: lo.ToPtr(clients.Recursive)}

func TestWhenListingPinsOnKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the pin/ls request fails", func(t *testing.T) {
		t.Parallel()

		requestErr := errors.New("connection refused")

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().PinLsWithResponse(anyCtx, recursivePinLsParams).Return(nil, requestErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ListPins(context.Background())

			assert.ErrorIs(t, err, requestErr)
		})
	})

	t.Run("Given Kubo fails without a JSON error body", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				PinLsWithResponse(anyCtx, recursivePinLsParams).
				Return(&clients.PinLsResponse{
					HTTPResponse: newHttpResponse(http.StatusForbidden),
					Body:         []byte("forbidden"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error carrying the raw body", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ListPins(context.Background())

			assert.EqualError(t, err, "kubo pin/ls failed with status 403: forbidden")
		})
	})

	t.Run("Given Kubo fails without any body", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				PinLsWithResponse(anyCtx, recursivePinLsParams).
				Return(&clients.PinLsResponse{HTTPResponse: newHttpResponse(http.StatusBadGateway)}, nil).
				Times(1)
		}

		t.Run("Should return an error carrying the status", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ListPins(context.Background())

			assert.EqualError(t, err, "kubo pin/ls failed with status 502")
		})
	})

	t.Run("Given the node has pins", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				PinLsWithResponse(anyCtx, recursivePinLsParams).
				Return(&clients.PinLsResponse{
					HTTPResponse: newHttpResponse(http.StatusOK),
					JSON200: &clients.PinLsResult{Keys: map[string]clients.PinLsEntry{
						"bafyb": {Type: "recursive"},
						"bafya": {Type: "recursive"},
					}},
				}, nil).
				Times(1)
		}

		t.Run("Should return the pins sorted by CID", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			pins, err := suite.sut.ListPins(context.Background())

			require.NoError(t, err)
			assert.Equal(t, []models.Pin{{Cid: "bafya", Type: "recursive"}, {Cid: "bafyb", Type: "recursive"}}, pins)
		})
	})
}
