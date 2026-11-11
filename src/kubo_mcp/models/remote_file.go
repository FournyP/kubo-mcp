package models

import "io"

// RemoteFile is an open download; SizeBytes is -1 when unknown.
type RemoteFile struct {
	Name        string
	ContentType string
	SizeBytes   int64
	Body        io.ReadCloser
}
