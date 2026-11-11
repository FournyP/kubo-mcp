package interfaces

import (
	"context"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
)

type RemoteFileDownloaderInterface interface {
	// Download opens rawUrl; the caller closes Body. Non-public IPs are refused by default.
	Download(ctx context.Context, rawUrl string) (*models.RemoteFile, error)
}
