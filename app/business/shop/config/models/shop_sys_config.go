package models

import "nova-factory-server/app/baize"

// ShopSysConfig 商城系统配置
type ShopSysConfig struct {
	ID          int64  `json:"id,string" gorm:"id"`                   // 主键ID
	ConfigKey   string `json:"configKey" gorm:"config_key"`           // 配置键名
	ConfigValue string `json:"configValue" gorm:"config_value"`       // 配置键值
	ConfigType  string `json:"configType" gorm:"config_type"`         // 配置类型
	Remark      string `json:"remark" gorm:"remark"`                  // 备注
	DeptID      int64  `json:"deptId" gorm:"column:dept_id" gorm:"-"` // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state" gorm:"state"` // 操作状态
}

// WechatConfigReq 微信小程序配置请求
type WechatConfigReq struct {
	AppID          string `json:"appId"`          // 微信小程序应用ID
	AppSecret      string `json:"appSecret"`      // 微信小程序应用密钥
	Token          string `json:"token"`          // 微信令牌
	EncodingAESKey string `json:"encodingAESKey"` // 消息加密密钥
	MchID          string `json:"mchId"`          // 微信支付商户号
	MchKey         string `json:"mchKey"`         // 微信支付商户密钥
	NotifyURL      string `json:"notifyUrl"`      // 微信支付回调地址
	CertPath       string `json:"certPath"`       // 微信支付证书路径
}

// WechatConfigResp 微信小程序配置响应
type WechatConfigResp struct {
	AppID          string `json:"appId"`          // 微信小程序应用ID
	AppSecret      string `json:"appSecret"`      // 微信小程序应用密钥（脱敏）
	Token          string `json:"token"`          // 微信令牌
	EncodingAESKey string `json:"encodingAESKey"` // 消息加密密钥
	MchID          string `json:"mchId"`          // 微信支付商户号
	MchKey         string `json:"mchKey"`         // 微信支付商户密钥（脱敏）
	NotifyURL      string `json:"notifyUrl"`      // 微信支付回调地址
	CertPath       string `json:"certPath"`       // 微信支付证书路径
}
