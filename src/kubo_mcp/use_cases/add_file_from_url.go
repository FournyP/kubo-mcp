package use_cases

import (
	"context"
	"fmt"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/helpers"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

const defaultFileName = "file"

type AddFileFromUrlInterface interface {
	Add(ctx context.Context, request *requests.AddFileFromUrlRequest) (*responses.AddedFileResponse, error)
}

type AddFileFromUrl struct {
	ipfsStorage          interfaces.IpfsStorageInterface
	remoteFileDownloader interfaces.RemoteFileDownloaderInterface
	ipfsSettings         *settings.IpfsSettings
}

func NewAddFileFromUrl(
	ipfsStorage interfaces.IpfsStorageInterface,
	remoteFileDownloader interfaces.RemoteFileDownloaderInterface,
	ipfsSettings *settings.IpfsSettings,
) *AddFileFromUrl {
	return &AddFileFromUrl{
		ipfsStorage:          ipfsStorage,
		remoteFileDownloader: remoteFileDownloader,
		ipfsSettings:         ipfsSettings,
	}
}

func (u *AddFileFromUrl) Add(
	ctx context.Context,
	request *requests.AddFileFromUrlRequest,
) (*responses.AddedFileResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	remoteFile, err := u.remoteFileDownloader.Download(ctx, request.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to download the file: %w", err)
	}
	defer remoteFile.Body.Close()

	if remoteFile.SizeBytes > u.ipfsSettings.MaxUploadSizeBytes {
		return nil, errors.NewFileTooLargeError(u.ipfsSettings.MaxUploadSizeBytes)
	}

	// Content-Length can be missing or wrong, so the stream is capped too.
	sizeLimitedReader := helpers.NewSizeLimitedReader(remoteFile.Body, u.ipfsSettings.MaxUploadSizeBytes)

	storedFile, err := u.ipfsStorage.Add(ctx, u.resolveFileName(request, remoteFile), sizeLimitedReader)
	if sizeLimitedReader.Exceeded() {
		return nil, errors.NewFileTooLargeError(u.ipfsSettings.MaxUploadSizeBytes)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to add the file to IPFS: %w", err)
	}
	// Kubo's size includes DAG overhead on chunked files.
	storedFile.SizeBytes = sizeLimitedReader.ReadBytes()

	return responses.NewAddedFileResponse(storedFile, u.ipfsSettings.PublicGatewayUrl), nil
}

func (u *AddFileFromUrl) resolveFileName(
	request *requests.AddFileFromUrlRequest,
	remoteFile *models.RemoteFile,
) string {
	if request.FileName != "" {
		return request.FileName
	}
	if remoteFile.Name != "" {
		return remoteFile.Name
	}
	return defaultFileName
}
