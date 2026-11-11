package use_cases

import (
	"context"
	"fmt"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
)

type UnpinFileInterface interface {
	Unpin(ctx context.Context, request *requests.UnpinFileRequest) (*responses.PinStateResponse, error)
}

type UnpinFile struct {
	ipfsStorage interfaces.IpfsStorageInterface
}

func NewUnpinFile(ipfsStorage interfaces.IpfsStorageInterface) *UnpinFile {
	return &UnpinFile{ipfsStorage: ipfsStorage}
}

func (u *UnpinFile) Unpin(
	ctx context.Context,
	request *requests.UnpinFileRequest,
) (*responses.PinStateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	if err := u.ipfsStorage.Unpin(ctx, request.Cid); err != nil {
		return nil, fmt.Errorf("failed to unpin the file: %w", err)
	}

	return &responses.PinStateResponse{Cid: request.Cid, Pinned: false}, nil
}
