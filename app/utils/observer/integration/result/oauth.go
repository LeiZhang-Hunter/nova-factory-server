// 定义 OAuth 授权相关的响应接口。
// 用于封装第三方系统（如管家婆）OAuth 授权流程中
// 通过 oauthcode 换取访问令牌的返回结果。
package result

// OAuthTokenResponse OAuth授权令牌响应接口。
// 管家婆等系统通过 oauthcode 换取 token 时返回此结构，
// 包含令牌本身、过期时间、签发时间及关联的 AppKey/AppSecret。
type OAuthTokenResponse interface {
	// GetCode 返回业务状态码，0 表示成功
	GetCode() int64
	// GetMessage 返回业务消息，成功或失败描述
	GetMessage() string
	// GetToken 访问令牌，用于后续 API 调用的身份认证
	GetToken() string
	// GetExpireDate 令牌过期时间
	GetExpireDate() string
	// GetIssueDate 令牌签发时间
	GetIssueDate() string
	// GetAppKey 返回当前授权关联的 app key
	GetAppKey() string
	// GetAppSecret 返回当前授权关联的 app secret
	GetAppSecret() string

	base
}
