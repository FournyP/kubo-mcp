package responses

type PinResponse struct {
	Cid  string `json:"cid"`
	Type string `json:"type"`
}

type ListPinsResponse struct {
	Pins []PinResponse `json:"pins"`
}
