package responses

type PinStateResponse struct {
	Cid    string `json:"cid"`
	Pinned bool   `json:"pinned"`
}
