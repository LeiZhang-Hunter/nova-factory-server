package models

import shopusermodels "nova-factory-server/app/business/shop/user/models"

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

// ShopGetInfoResp 商城用户鉴权信息
type ShopGetInfoResp struct {
	User *ShopAuthUserInfo `json:"user"` // 用户信息
}

// UserToAuthUserInfo converts a User model to ShopAuthUserInfo
func UserToAuthUserInfo(user *shopusermodels.User) *ShopAuthUserInfo {
	if user == nil {
		return nil
	}
	return &ShopAuthUserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Mobile:       user.Mobile,
		Avatar:       user.Avatar,
		UserType:     user.UserType,
		CompanyName:  user.CompanyName,
		ContactName:  user.ContactName,
		ContactPhone: user.ContactPhone,
		DeptID:       user.DeptID,
	}
}
