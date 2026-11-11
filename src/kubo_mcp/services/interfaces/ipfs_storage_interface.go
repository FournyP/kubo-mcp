package interfaces

import (
	"context"
	"io"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
)

type IpfsStorageInterface interface {
	// Add pins content as one file under a CIDv1; SizeBytes includes DAG overhead.
	Add(ctx context.Context, fileName string, content io.Reader) (*models.StoredFile, error)
	// Read returns at most maxBytes of path (`<cid>[/sub/path]`).
	Read(ctx context.Context, path string, maxBytes int64) ([]byte, error)
	Stat(ctx context.Context, path string) (*models.FileInfo, error)
	Pin(ctx context.Context, cid string) error
	// Unpin keeps the blocks until Kubo garbage-collects them.
	Unpin(ctx context.Context, cid string) error
	// ListPins returns recursive pins only.
	ListPins(ctx context.Context) ([]models.Pin, error)
}
