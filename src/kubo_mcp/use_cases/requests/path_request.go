package requests

import (
	"strings"

	"github.com/ipfs/go-cid"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

// PathRequest is a CID with an optional sub-path (`<cid>/images/logo.png`).
type PathRequest struct {
	Path string
}

func (r *PathRequest) Validate() error {
	if strings.TrimSpace(r.Path) == "" {
		return errors.NewValidationError("path", "path is required")
	}

	root, _, _ := strings.Cut(r.Path, "/")
	if _, err := cid.Decode(root); err != nil {
		return errors.NewValidationError("path", "path must start with a valid CID")
	}

	return nil
}
