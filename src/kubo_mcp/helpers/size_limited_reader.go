package helpers

import (
	"io"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_common/errors"
)

type SizeLimitedReader struct {
	reader       io.Reader
	maxSizeBytes int64
	readBytes    int64
}

func NewSizeLimitedReader(reader io.Reader, maxSizeBytes int64) *SizeLimitedReader {
	return &SizeLimitedReader{reader: reader, maxSizeBytes: maxSizeBytes}
}

func (r *SizeLimitedReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.readBytes += int64(n)
	if r.Exceeded() {
		return n, errors.NewFileTooLargeError(r.maxSizeBytes)
	}
	return n, err
}

func (r *SizeLimitedReader) ReadBytes() int64 {
	return r.readBytes
}

// Exceeded survives transports that swallow the read error.
func (r *SizeLimitedReader) Exceeded() bool {
	return r.readBytes > r.maxSizeBytes
}
