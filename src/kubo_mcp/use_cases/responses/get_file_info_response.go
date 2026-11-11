package responses

type GetFileInfoResponse struct {
	Cid                 string `json:"cid"`
	Type                string `json:"type"`
	SizeBytes           int64  `json:"sizeBytes"`
	CumulativeSizeBytes int64  `json:"cumulativeSizeBytes"`
	Blocks              int64  `json:"blocks"`
	GatewayUrl          string `json:"gatewayUrl"`
}
