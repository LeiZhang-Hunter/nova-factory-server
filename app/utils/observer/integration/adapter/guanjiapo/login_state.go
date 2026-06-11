package guanjiapo

// LoginState 登录状态
type loginState struct {
	Online   bool   `json:"online"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	CheckURL string `json:"checkUrl"`
	Raw      string `json:"raw,omitempty"`
}

func (l *loginState) GetOnline() bool {
	return l.Online
}
func (l *loginState) GetMessage() string {
	return l.Message
}
func (l *loginState) GetType() string {
	return l.Type
}
func (l *loginState) GetCheckURL() string {
	return l.CheckURL
}
func (l *loginState) GetRaw() string {
	return l.Raw
}
