package models

// ShopLoginReq 商城用户登录参数
type ShopLoginReq struct {
	Account  string `json:"account" binding:"required"`  // 登录账号，支持用户名或手机号
	Password string `json:"password" binding:"required"` // 登录密码
}

// ShopAuthUserInfo 商城登录用户信息
type ShopAuthUserInfo struct {
	ID           int64  `json:"id,string"`    // 用户ID
	Username     string `json:"username"`     // 用户名
	Nickname     string `json:"nickname"`     // 用户昵称
	Mobile       string `json:"mobile"`       // 手机号
	Avatar       string `json:"avatar"`       // 头像
	UserType     int32  `json:"userType"`     // 用户类型
	CompanyName  string `json:"companyName"`  // 企业名称
	ContactName  string `json:"contactName"`  // 联系人
	ContactPhone string `json:"contactPhone"` // 联系人手机号
	DeptID       int64  `json:"deptId"`       // 部门ID
}

// ShopLoginResp 商城用户登录返回
type ShopLoginResp struct {
	Token       string            `json:"token"`       // 访问令牌
	TokenType   string            `json:"tokenType"`   // 令牌类型
	User        *ShopAuthUserInfo `json:"user"`        // 用户信息
	Permissions []string          `json:"permissions"` // 权限列表
}

// ShopGetInfoResp 商城用户鉴权信息
type ShopGetInfoResp struct {
	User        *ShopAuthUserInfo `json:"user"`        // 用户信息
	Permissions []string          `json:"permissions"` // 权限列表
}
