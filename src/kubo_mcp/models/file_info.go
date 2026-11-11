package models

type FileInfo struct {
	Cid                 string
	Type                string
	SizeBytes           int64
	CumulativeSizeBytes int64
	Blocks              int64
}
