// 管家婆全渠道 Token 签发请求参数。
// IssueTokenReq 由 /oauth/token 接口接收，包含签发和刷新两种模式所需的全部字段。
package models

// IssueTokenReq Token 签发/刷新请求参数
// grant_type 支持 authorization_code 和 refresh_token 两种模式
type IssueTokenReq struct {
	AppKey       string `json:"appkey"`
	AppSecret    string `json:"appsecret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}
