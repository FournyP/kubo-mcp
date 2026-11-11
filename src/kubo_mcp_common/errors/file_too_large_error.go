package errors

import "fmt"

type FileTooLargeError struct {
	MaxSizeBytes int64 `json:"maxSizeBytes"`
}

func NewFileTooLargeError(maxSizeBytes int64) *FileTooLargeError {
	return &FileTooLargeError{MaxSizeBytes: maxSizeBytes}
}

func (e *FileTooLargeError) Error() string {
	return fmt.Sprintf("file exceeds the %d bytes limit", e.MaxSizeBytes)
}
