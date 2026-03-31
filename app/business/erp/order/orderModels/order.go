package orderModels

type CheckLoginStateReq struct {
	CheckURL string `form:"checkUrl"`
}

type CheckLoginStateResp struct {
	Online   bool   `json:"online"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	CheckURL string `json:"checkUrl"`
}
