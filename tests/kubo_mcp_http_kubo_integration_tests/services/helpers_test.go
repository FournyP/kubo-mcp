package services_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/clients"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/services"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_kubo/settings"
	"github.com/FournyP/kubo-mcp/tests/common_tests"
	"github.com/FournyP/kubo-mcp/tests/common_tests/container_builders"
)

type KuboIpfsStorageTestingSuite struct {
	sut *services.KuboIpfsStorage
}

func KuboIpfsStorageBeforeEach(t *testing.T) *KuboIpfsStorageTestingSuite {
	t.Helper()

	kuboSettings := &settings.KuboSettings{
		ApiUrl:         container_builders.GetKuboContainerBuilder().ApiUrl(),
		RequestTimeout: 30 * time.Second,
		UploadTimeout:  time.Minute,
	}
	client, err := clients.NewClientWithResponses(kuboSettings.ApiUrl)
	require.NoError(t, err)

	return &KuboIpfsStorageTestingSuite{sut: services.NewKuboIpfsStorage(client, kuboSettings)}
}

// addRandomFile stores unique content and returns it with its stored file.
func (suite *KuboIpfsStorageTestingSuite) addRandomFile(t *testing.T, size int) ([]byte, *models.StoredFile) {
	t.Helper()

	content := common_tests.RandomContent(size)
	storedFile, err := suite.sut.Add(context.Background(), "file.bin", bytes.NewReader(content))
	require.NoError(t, err)

	return content, storedFile
}

func recursivePin(cid string) models.Pin {
	return models.Pin{Cid: cid, Type: "recursive"}
}
