package requests

import (
	"net/url"
	"strings"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

type AddFileFromUrlRequest struct {
	FileNameRequest

	Url string
}

func (r *AddFileFromUrlRequest) Validate() error {
	if strings.TrimSpace(r.Url) == "" {
		return errors.NewValidationError("url", "url is required")
	}

	parsedUrl, err := url.Parse(r.Url)
	if err != nil || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") || parsedUrl.Host == "" {
		return errors.NewValidationError("url", "url must be an absolute http(s) URL")
	}

	return r.FileNameRequest.Validate()
}
