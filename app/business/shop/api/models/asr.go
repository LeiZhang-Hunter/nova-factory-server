package models

type ShopASRSubmitReq struct {
	Type        string `json:"type"`
	FileName    string `json:"file_name"`
	AudioBase64 string `json:"audio_base64,omitempty"`
}

type ShopASRServerEvent struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`
}
