package requests

import (
	"strings"

	"github.com/ipfs/go-cid"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

type CidRequest struct {
	Cid string
}

func (r *CidRequest) Validate() error {
	if strings.TrimSpace(r.Cid) == "" {
		return errors.NewValidationError("cid", "cid is required")
	}

	if _, err := cid.Decode(r.Cid); err != nil {
		return errors.NewValidationError("cid", "cid is not a valid CID")
	}

	return nil
}
