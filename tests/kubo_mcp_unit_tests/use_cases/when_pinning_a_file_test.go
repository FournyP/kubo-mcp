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

type WhenPinningAFileTestingSuite struct {
	sut *use_cases.PinFile

	mockIpfsStorage *interfaces_mocks.MockIpfsStorageInterface
}

func WhenPinningAFileBeforeEach(t *testing.T) *WhenPinningAFileTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)

	return &WhenPinningAFileTestingSuite{
		sut:             use_cases.NewPinFile(mockIpfsStorage),
		mockIpfsStorage: mockIpfsStorage,
	}
}

func newPinFileRequest(cid string) *requests.PinFileRequest {
	return &requests.PinFileRequest{CidRequest: requests.CidRequest{Cid: cid}}
}

func TestWhenPinningAFile(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid request", func(t *testing.T) {
		t.Parallel()

		request := newPinFileRequest("not-a-cid")

		initSuite := func(suite *WhenPinningAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Pin(gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenPinningAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Pin(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "cid", validationError.Field)
		})
	})

	t.Run("Given pinning the file fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cid := common_tests.RandomCid()
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenPinningAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Pin(ctx, cid).Return(storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPinningAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Pin(ctx, newPinFileRequest(cid))

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cid := common_tests.RandomCid()

		initSuite := func(suite *WhenPinningAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Pin(ctx, cid).Return(nil).Times(1)
		}

		t.Run("Should report the new pin state", func(t *testing.T) {
			t.Parallel()

			suite := WhenPinningAFileBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Pin(ctx, newPinFileRequest(cid))

			require.NoError(t, err)
			assert.Equal(t, cid, response.Cid)
			assert.Equal(t, true, response.Pinned)
		})
	})
}
