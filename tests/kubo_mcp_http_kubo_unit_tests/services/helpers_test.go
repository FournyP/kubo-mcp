package services_test

import (
	"net/http"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/services"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/settings"
	"github.com/FournyP/kubo-mcp/tests/kubo_mcp_http_kubo_mocks/clients_mocks"
)

// anyCtx: the storage derives a timeout context, so it never equals the caller's.
var anyCtx = gomock.Any()

type KuboIpfsStorageTestingSuite struct {
	sut *services.KuboIpfsStorage

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func KuboIpfsStorageBeforeEach(t *testing.T) *KuboIpfsStorageTestingSuite {
	mockController := gomock.NewController(t)
	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)
	kuboSettings := &settings.KuboSettings{RequestTimeout: 5 * time.Second, UploadTimeout: 5 * time.Second}

	return &KuboIpfsStorageTestingSuite{
		sut:        services.NewKuboIpfsStorage(mockClient, kuboSettings),
		mockClient: mockClient,
	}
}

func newHttpResponse(statusCode int) *http.Response {
	return &http.Response{StatusCode: statusCode, Trailer: http.Header{}}
}

func newKuboError(message string) *clients.KuboError {
	return &clients.KuboError{Message: message}
}
