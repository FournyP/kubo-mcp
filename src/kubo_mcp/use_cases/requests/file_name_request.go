package requests

import (
	"strings"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

const maxFileNameLength = 255

// FileNameRequest checks the shape only; requests needing a name check presence first.
type FileNameRequest struct {
	FileName string
}

func (r *FileNameRequest) Validate() error {
	if len(r.FileName) > maxFileNameLength {
		return errors.NewValidationError("fileName", "fileName must be at most 255 characters")
	}

	if strings.ContainsAny(r.FileName, "/\\\x00") {
		return errors.NewValidationError("fileName", "fileName must not contain path separators")
	}

	return nil
}
