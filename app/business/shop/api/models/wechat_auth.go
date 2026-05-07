package models

// WechatLoginReq 微信小程序登录参数
type WechatLoginReq struct {
	Code     string `json:"code" binding:"required"` // 微信登录code
	Nickname string `json:"nickname"`                // 用户昵称
	Avatar   string `json:"avatar"`                  // 用户头像
}

// RefreshTokenReq 刷新Token参数
type RefreshTokenReq struct {
	Token string `json:"token" binding:"required"` // 有效的JWT token
}

// WechatLoginResp 微信登录返回
type WechatLoginResp struct {
	Token  string `json:"token"`         // 访问令牌
	UserId int64  `json:"userId,string"` // 用户ID
}

// WechatUserCreate 微信用户创建参数
type WechatUserCreate struct {
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 用户昵称
	Avatar   string `json:"avatar"`   // 头像
	UserType int32  `json:"userType"` // 用户类型
	Status   *bool  `json:"status"`   // 状态
	Openid   string `json:"openid"`   // 微信openid
}
