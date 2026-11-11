package use_cases

import (
	"context"
	"fmt"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
)

type PinFileInterface interface {
	Pin(ctx context.Context, request *requests.PinFileRequest) (*responses.PinStateResponse, error)
}

type PinFile struct {
	ipfsStorage interfaces.IpfsStorageInterface
}

func NewPinFile(ipfsStorage interfaces.IpfsStorageInterface) *PinFile {
	return &PinFile{ipfsStorage: ipfsStorage}
}

func (u *PinFile) Pin(ctx context.Context, request *requests.PinFileRequest) (*responses.PinStateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	if err := u.ipfsStorage.Pin(ctx, request.Cid); err != nil {
		return nil, fmt.Errorf("failed to pin the file: %w", err)
	}

	return &responses.PinStateResponse{Cid: request.Cid, Pinned: true}, nil
}
