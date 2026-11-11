package use_cases

import (
	"context"
	"fmt"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
)

type GetFileInfoInterface interface {
	Get(ctx context.Context, request *requests.GetFileInfoRequest) (*responses.GetFileInfoResponse, error)
}

type GetFileInfo struct {
	ipfsStorage  interfaces.IpfsStorageInterface
	ipfsSettings *settings.IpfsSettings
}

func NewGetFileInfo(ipfsStorage interfaces.IpfsStorageInterface, ipfsSettings *settings.IpfsSettings) *GetFileInfo {
	return &GetFileInfo{ipfsStorage: ipfsStorage, ipfsSettings: ipfsSettings}
}

func (u *GetFileInfo) Get(
	ctx context.Context,
	request *requests.GetFileInfoRequest,
) (*responses.GetFileInfoResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	fileInfo, err := u.ipfsStorage.Stat(ctx, request.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat the file on IPFS: %w", err)
	}

	return &responses.GetFileInfoResponse{
		Cid:                 fileInfo.Cid,
		Type:                fileInfo.Type,
		SizeBytes:           fileInfo.SizeBytes,
		CumulativeSizeBytes: fileInfo.CumulativeSizeBytes,
		Blocks:              fileInfo.Blocks,
		GatewayUrl:          u.ipfsSettings.PublicGatewayUrl + "/ipfs/" + request.Path,
	}, nil
}
