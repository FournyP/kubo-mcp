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
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
	"github.com/FournyP/kubo-mcp/tests/kubo_mcp_mocks/services_mocks/interfaces_mocks"
)

type WhenListingPinsTestingSuite struct {
	sut *use_cases.ListPins

	mockIpfsStorage *interfaces_mocks.MockIpfsStorageInterface
}

func WhenListingPinsBeforeEach(t *testing.T) *WhenListingPinsTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)

	return &WhenListingPinsTestingSuite{
		sut:             use_cases.NewListPins(mockIpfsStorage),
		mockIpfsStorage: mockIpfsStorage,
	}
}

func TestWhenListingPins(t *testing.T) {
	t.Parallel()

	t.Run("Given listing the pins fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenListingPinsTestingSuite) {
			suite.mockIpfsStorage.EXPECT().ListPins(ctx).Return(nil, storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingPinsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.List(ctx)

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given the node has no pins", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		initSuite := func(suite *WhenListingPinsTestingSuite) {
			suite.mockIpfsStorage.EXPECT().ListPins(ctx).Return(nil, nil).Times(1)
		}

		t.Run("Should return an empty list", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingPinsBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.List(ctx)

			require.NoError(t, err)
			assert.NotNil(t, response.Pins)
			assert.Empty(t, response.Pins)
		})
	})

	t.Run("Given the node has pins", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		firstCid := common_tests.RandomCid()
		secondCid := common_tests.RandomCid()

		initSuite := func(suite *WhenListingPinsTestingSuite) {
			suite.mockIpfsStorage.EXPECT().
				ListPins(ctx).
				Return([]models.Pin{{Cid: firstCid, Type: "recursive"}, {Cid: secondCid, Type: "recursive"}}, nil).
				Times(1)
		}

		t.Run("Should return every pin", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingPinsBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.List(ctx)

			require.NoError(t, err)
			assert.Equal(t, []responses.PinResponse{
				{Cid: firstCid, Type: "recursive"},
				{Cid: secondCid, Type: "recursive"},
			}, response.Pins)
		})
	})
}
