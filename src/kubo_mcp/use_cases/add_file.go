package use_cases

import (
	"bytes"
	"context"
	"fmt"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

type AddFileInterface interface {
	Add(ctx context.Context, request *requests.AddFileRequest) (*responses.AddedFileResponse, error)
}

type AddFile struct {
	ipfsStorage  interfaces.IpfsStorageInterface
	ipfsSettings *settings.IpfsSettings
}

func NewAddFile(ipfsStorage interfaces.IpfsStorageInterface, ipfsSettings *settings.IpfsSettings) *AddFile {
	return &AddFile{ipfsStorage: ipfsStorage, ipfsSettings: ipfsSettings}
}

func (u *AddFile) Add(ctx context.Context, request *requests.AddFileRequest) (*responses.AddedFileResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	content, err := request.DecodeContent()
	if err != nil {
		return nil, err
	}

	if int64(len(content)) > u.ipfsSettings.MaxUploadSizeBytes {
		return nil, errors.NewFileTooLargeError(u.ipfsSettings.MaxUploadSizeBytes)
	}

	storedFile, err := u.ipfsStorage.Add(ctx, request.FileName, bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to add the file to IPFS: %w", err)
	}
	// Kubo's size includes DAG overhead on chunked files.
	storedFile.SizeBytes = int64(len(content))

	return responses.NewAddedFileResponse(storedFile, u.ipfsSettings.PublicGatewayUrl), nil
}
