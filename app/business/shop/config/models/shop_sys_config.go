package models

import "nova-factory-server/app/baize"

// ShopSysConfig 商城系统配置
type ShopSysConfig struct {
	ID          int64  `json:"id,string" db:"id"`                    // 主键ID
	ConfigKey   string `json:"configKey" db:"config_key"`            // 配置键名
	ConfigValue string `json:"configValue" db:"config_value"`       // 配置键值
	ConfigType  string `json:"configType" db:"config_type"`         // 配置类型
	Remark      string `json:"remark" db:"remark"`                  // 备注
	DeptID      int64  `json:"deptId" gorm:"column:dept_id" db:"-"` // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state" db:"state"` // 操作状态
}

// WechatConfigReq 微信小程序配置请求
type WechatConfigReq struct {
	AppID           string `json:"appId"`            // 微信小程序应用ID
	AppSecret       string `json:"appSecret"`        // 微信小程序应用密钥
	Token           string `json:"token"`            // 微信令牌
	EncodingAESKey  string `json:"encodingAESKey"`   // 消息加密密钥
}

// WechatConfigResp 微信小程序配置响应
type WechatConfigResp struct {
	AppID          string `json:"appId"`           // 微信小程序应用ID
	AppSecret      string `json:"appSecret"`       // 微信小程序应用密钥（脱敏）
	Token          string `json:"token"`           // 微信令牌
	EncodingAESKey string `json:"encodingAESKey"`  // 消息加密密钥
}
