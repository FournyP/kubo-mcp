package use_cases

import (
	"context"
	"fmt"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
)

type ListPinsInterface interface {
	List(ctx context.Context) (*responses.ListPinsResponse, error)
}

type ListPins struct {
	ipfsStorage interfaces.IpfsStorageInterface
}

func NewListPins(ipfsStorage interfaces.IpfsStorageInterface) *ListPins {
	return &ListPins{ipfsStorage: ipfsStorage}
}

func (u *ListPins) List(ctx context.Context) (*responses.ListPinsResponse, error) {
	pins, err := u.ipfsStorage.ListPins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list the pins: %w", err)
	}

	response := &responses.ListPinsResponse{Pins: make([]responses.PinResponse, 0, len(pins))}
	for _, pin := range pins {
		response.Pins = append(response.Pins, responses.PinResponse{Cid: pin.Cid, Type: pin.Type})
	}

	return response, nil
}
