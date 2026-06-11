package api

// LoginState 登录状态
type LoginState interface {
	GetOnline() bool
	GetMessage() string
	GetType() string
	GetCheckURL() string
	GetRaw() string
}
