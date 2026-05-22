package models

type ShopVoiceSubmitReq struct {
	Type           string `json:"type"`
	ConversationID int64  `json:"conversation_id,string"`
	TabID          string `json:"tab_id"`
	FileName       string `json:"file_name"`
	AudioBase64    string `json:"audio_base64,omitempty"`
}

type ShopVoiceServerEvent struct {
	Type           string `json:"type"`
	ConversationID int64  `json:"conversation_id,string,omitempty"`
	Text           string `json:"text,omitempty"`
	Message        string `json:"message,omitempty"`
	Mime           string `json:"mime,omitempty"`
	SampleRate     int    `json:"sample_rate,omitempty"`
	Channels       int    `json:"channels,omitempty"`
	BitDepth       int    `json:"bit_depth,omitempty"`
	Seq            int    `json:"seq,omitempty"`
}
