package responses

import (
	"net/url"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp/models"
)

type AddedFileResponse struct {
	Cid        string `json:"cid"`
	Name       string `json:"name"`
	SizeBytes  int64  `json:"sizeBytes"`
	IpfsUri    string `json:"ipfsUri"`
	GatewayUrl string `json:"gatewayUrl"`
}

// NewAddedFileResponse adds `filename` so gateways serve the original name.
func NewAddedFileResponse(storedFile *models.StoredFile, publicGatewayUrl string) *AddedFileResponse {
	return &AddedFileResponse{
		Cid:        storedFile.Cid,
		Name:       storedFile.Name,
		SizeBytes:  storedFile.SizeBytes,
		IpfsUri:    "ipfs://" + storedFile.Cid,
		GatewayUrl: publicGatewayUrl + "/ipfs/" + storedFile.Cid + "?filename=" + url.QueryEscape(storedFile.Name),
	}
}
