package devicemonitormodel

type CameraOfferReq struct {
	DeviceId  string `json:"device_id" binding:"required"`
	ChannelId string `json:"channel_id"`
	SDP64     string `json:"sdp64" binding:"required"`
	Node      string `json:"node"`
	TimeoutMS int64  `json:"timeout_ms"`
}

type CameraOfferRes struct {
	Token   string `json:"token"`
	PlayURL string `json:"play_url"`
	WhepURL string `json:"whep_url"`
	SDP64   string `json:"sdp64"`
}
