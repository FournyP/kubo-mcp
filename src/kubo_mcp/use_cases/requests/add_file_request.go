package requests

import (
	"encoding/base64"
	"strings"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

type AddFileRequest struct {
	FileNameRequest

	// ContentBase64 also accepts a data URI.
	ContentBase64 string
}

func (r *AddFileRequest) Validate() error {
	if strings.TrimSpace(r.FileName) == "" {
		return errors.NewValidationError("fileName", "fileName is required")
	}

	if err := r.FileNameRequest.Validate(); err != nil {
		return err
	}

	if strings.TrimSpace(r.ContentBase64) == "" {
		return errors.NewValidationError("contentBase64", "contentBase64 is required")
	}

	return nil
}

func (r *AddFileRequest) DecodeContent() ([]byte, error) {
	encoded := strings.TrimSpace(r.ContentBase64)
	if strings.HasPrefix(encoded, "data:") {
		_, encoded, _ = strings.Cut(encoded, ",")
	}

	if content, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		return content, nil
	}

	content, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.NewValidationError("contentBase64", "contentBase64 is not valid base64")
	}

	return content, nil
}
