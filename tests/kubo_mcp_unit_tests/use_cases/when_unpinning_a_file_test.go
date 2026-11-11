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

type WhenUnpinningAFileTestingSuite struct {
	sut *use_cases.UnpinFile

	mockIpfsStorage *interfaces_mocks.MockIpfsStorageInterface
}

func WhenUnpinningAFileBeforeEach(t *testing.T) *WhenUnpinningAFileTestingSuite {
	mockController := gomock.NewController(t)
	mockIpfsStorage := interfaces_mocks.NewMockIpfsStorageInterface(mockController)

	return &WhenUnpinningAFileTestingSuite{
		sut:             use_cases.NewUnpinFile(mockIpfsStorage),
		mockIpfsStorage: mockIpfsStorage,
	}
}

func newUnpinFileRequest(cid string) *requests.UnpinFileRequest {
	return &requests.UnpinFileRequest{CidRequest: requests.CidRequest{Cid: cid}}
}

func TestWhenUnpinningAFile(t *testing.T) {
	t.Parallel()

	t.Run("Given an invalid request", func(t *testing.T) {
		t.Parallel()

		request := newUnpinFileRequest("not-a-cid")

		initSuite := func(suite *WhenUnpinningAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Unpin(gomock.Any(), gomock.Any()).Times(0)
		}

		t.Run("Should return a validation error without touching dependencies", func(t *testing.T) {
			t.Parallel()

			suite := WhenUnpinningAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Unpin(context.Background(), request)

			var validationError *common_errors.ValidationError
			require.ErrorAs(t, err, &validationError)
			assert.Equal(t, "cid", validationError.Field)
		})
	})

	t.Run("Given unpinning the file fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cid := common_tests.RandomCid()
		storageErr := errors.New("kubo unreachable")

		initSuite := func(suite *WhenUnpinningAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Unpin(ctx, cid).Return(storageErr).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUnpinningAFileBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Unpin(ctx, newUnpinFileRequest(cid))

			assert.ErrorIs(t, err, storageErr)
		})
	})

	t.Run("Given a valid request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cid := common_tests.RandomCid()

		initSuite := func(suite *WhenUnpinningAFileTestingSuite) {
			suite.mockIpfsStorage.EXPECT().Unpin(ctx, cid).Return(nil).Times(1)
		}

		t.Run("Should report the new pin state", func(t *testing.T) {
			t.Parallel()

			suite := WhenUnpinningAFileBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Unpin(ctx, newUnpinFileRequest(cid))

			require.NoError(t, err)
			assert.Equal(t, cid, response.Cid)
			assert.Equal(t, false, response.Pinned)
		})
	})
}
