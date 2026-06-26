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

type ShopSysConfigWechatPayConfigDTO struct {
	AppId                 string `json:"appId"`
	MchId                 string `json:"mchId"`
	ApiV3Key              string `json:"apiV3Key"`
	SerialNo              string `json:"serialNo"`
	PrivateKeyPath        string `json:"privateKeyPath"`
	NotifyUrl             string `json:"notifyUrl"`
	PlatformPublicKeyId   string `json:"platformPublicKeyId"`
	PlatformPublicKeyPath string `json:"platformPublicKeyPath"`
}
