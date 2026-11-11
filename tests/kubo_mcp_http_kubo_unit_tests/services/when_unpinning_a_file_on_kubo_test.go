package services_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
)

func TestWhenUnpinningAFileOnKubo(t *testing.T) {
	t.Parallel()

	t.Run("Given the pin/rm request fails", func(t *testing.T) {
		t.Parallel()

		cid := common_tests.RandomCid()
		requestErr := errors.New("connection refused")

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				PinRmWithResponse(anyCtx, &clients.PinRmParams{Arg: cid}).
				Return(nil, requestErr).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Unpin(context.Background(), cid)

			assert.ErrorIs(t, err, requestErr)
		})
	})

	t.Run("Given Kubo fails", func(t *testing.T) {
		t.Parallel()

		cid := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				PinRmWithResponse(anyCtx, &clients.PinRmParams{Arg: cid}).
				Return(&clients.PinRmResponse{
					HTTPResponse: newHttpResponse(http.StatusInternalServerError),
					JSONDefault:  newKuboError("not pinned or pinned indirectly"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error carrying the Kubo message", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Unpin(context.Background(), cid)

			assert.EqualError(t, err, "kubo pin/rm failed with status 500: not pinned or pinned indirectly")
		})
	})

	t.Run("Given Kubo succeeds", func(t *testing.T) {
		t.Parallel()

		cid := common_tests.RandomCid()

		initSuite := func(suite *KuboIpfsStorageTestingSuite) {
			suite.mockClient.EXPECT().
				PinRmWithResponse(anyCtx, &clients.PinRmParams{Arg: cid}).
				Return(&clients.PinRmResponse{
					HTTPResponse: newHttpResponse(http.StatusOK),
					JSON200:      &clients.PinResult{Pins: []string{cid}},
				}, nil).
				Times(1)
		}

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := KuboIpfsStorageBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Unpin(context.Background(), cid)

			assert.NoError(t, err)
		})
	})
}
