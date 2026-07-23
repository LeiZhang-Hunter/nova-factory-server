package response

import modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"

// LoginResp 登录返回
type LoginResp struct {
	Token      string `json:"token"`      // 登录 token
	ExpireTime int64  `json:"expireTime"` // 有效时间，单位秒
}

type GetInfo struct {
	User        *modelentity.User `json:"user"`
	Roles       []string          `json:"roles"`
	Permissions []string          `json:"permissions"`
}
