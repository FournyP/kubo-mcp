package responses

type GetFileResponse struct {
	Path      string `json:"path"`
	MimeType  string `json:"mimeType"`
	SizeBytes int64  `json:"sizeBytes"`
	Content   []byte `json:"-"`
}
