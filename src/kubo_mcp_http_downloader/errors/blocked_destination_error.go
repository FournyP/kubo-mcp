package errors

import "fmt"

type BlockedDestinationError struct {
	Address string `json:"address"`
}

func NewBlockedDestinationError(address string) *BlockedDestinationError {
	return &BlockedDestinationError{Address: address}
}

func (e *BlockedDestinationError) Error() string {
	return fmt.Sprintf("destination %s is not on the public internet", e.Address)
}
