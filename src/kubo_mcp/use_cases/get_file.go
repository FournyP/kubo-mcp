package use_cases

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"path"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/services/interfaces"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/requests"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp/use_cases/responses"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

type GetFileInterface interface {
	Get(ctx context.Context, request *requests.GetFileRequest) (*responses.GetFileResponse, error)
}

type GetFile struct {
	ipfsStorage  interfaces.IpfsStorageInterface
	ipfsSettings *settings.IpfsSettings
}

func NewGetFile(ipfsStorage interfaces.IpfsStorageInterface, ipfsSettings *settings.IpfsSettings) *GetFile {
	return &GetFile{ipfsStorage: ipfsStorage, ipfsSettings: ipfsSettings}
}

func (u *GetFile) Get(ctx context.Context, request *requests.GetFileRequest) (*responses.GetFileResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// One extra byte tells an over-limit file apart from one exactly at the limit.
	content, err := u.ipfsStorage.Read(ctx, request.Path, u.ipfsSettings.MaxReadSizeBytes+1)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file from IPFS: %w", err)
	}

	if int64(len(content)) > u.ipfsSettings.MaxReadSizeBytes {
		return nil, errors.NewFileTooLargeError(u.ipfsSettings.MaxReadSizeBytes)
	}

	return &responses.GetFileResponse{
		Path:      request.Path,
		MimeType:  u.detectMimeType(request.Path, content),
		SizeBytes: int64(len(content)),
		Content:   content,
	}, nil
}

// detectMimeType trusts the extension first: sniffing reports SVG as XML.
func (u *GetFile) detectMimeType(filePath string, content []byte) string {
	if mimeType := mime.TypeByExtension(path.Ext(filePath)); mimeType != "" {
		return mimeType
	}
	return http.DetectContentType(content)
}
