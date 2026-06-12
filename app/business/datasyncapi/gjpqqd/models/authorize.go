// 管家婆全渠道 OAuth 授权相关的请求与响应结构。
// 定义授权码请求参数及 Token 换取成功后的响应结构。
package models

// AuthorizeReq OAuth 授权码请求参数
// 第三方系统访问 /oauth/authorize 时传递 appkey、appsecret、redirect_uri 和 state
type AuthorizeReq struct {
	AppKey    string `form:"appkey"`
	AppSecret string `form:"appsecret"`
	Redirect  string `form:"redirect_uri"`
	State     string `form:"state"`
}

// TokenResponse 使用授权码或 refresh_token 换取 access_token 的返回结果
// 包含 token 本身、过期时间、refresh_token 及关联的应用凭据
type TokenResponse struct {
	Token           string `json:"token"`
	ExpireDate      string `json:"expire_date"`
	RefreshToken    string `json:"refresh_token"`
	RefreshExpireAt string `json:"refresh_expire_at"`
	AppKey          string `json:"app_key"`
	AppSecret       string `json:"app_secret"`
	SelfMallAccount string `json:"self_mall_account"`
}
