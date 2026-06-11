package result

// OAuthTokenResponse 管家婆oauthcode换取token的返回结果
type OAuthTokenResponse interface {
	GetCode() int64
	GetMessage() string
	GetToken() string
	GetExpireDate() string
	GetIssueDate() string
	GetAppKey() string
	GetAppSecret() string
	base
}
