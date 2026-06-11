package guanjiapo

import "encoding/json"

// OAuthTokenResponse 管家婆oauthcode换取token的返回结果
type OAuthTokenResponse struct {
	Code       int64  `json:"code"`
	Message    string `json:"message"`
	Token      string `json:"token"`
	ExpireDate string `json:"expiredate"`
	IssueDate  string `json:"issuedate"`
	AppKey     string `json:"appkey"`
	AppSecret  string `json:"appsecret"`
}

func (o *OAuthTokenResponse) GetCode() int64 {
	return o.Code
}
func (o *OAuthTokenResponse) GetMessage() string {
	return o.Message
}
func (o *OAuthTokenResponse) GetToken() string {
	return o.Token
}
func (o *OAuthTokenResponse) GetExpireDate() string {
	return o.ExpireDate
}
func (o *OAuthTokenResponse) GetIssueDate() string {
	return o.IssueDate
}
func (o *OAuthTokenResponse) GetAppKey() string {
	return o.AppKey
}
func (o *OAuthTokenResponse) GetAppSecret() string {
	return o.AppSecret
}
func (o *OAuthTokenResponse) Ptr() any {
	return o
}
func (o *OAuthTokenResponse) RawStr() (string, error) {
	marshal, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
func (o *OAuthTokenResponse) MetaData() map[string]any {
	return make(map[string]any)
}
