// 管家婆全渠道 Token 签发相关的请求与响应结构。
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

// IssueTokenResponse Token 签发成功响应
// ReExpireDate 复用 ExpireDate，SelfMallAccount 为商家自建商城账号标识
type IssueTokenResponse struct {
	IsError         bool        `json:"iserror"`
	ErrorMsg        string      `json:"errormsg"`
	Token           interface{} `json:"token"`
	ExpireDate      interface{} `json:"expiredate"`
	RefreshToken    interface{} `json:"refresh_token"`
	ReExpireDate    interface{} `json:"re_expiredate"`
	AppKey          interface{} `json:"appkey"`
	AppSecret       interface{} `json:"appsecret"`
	SelfMallAccount interface{} `json:"selfmallaccount"`
}
