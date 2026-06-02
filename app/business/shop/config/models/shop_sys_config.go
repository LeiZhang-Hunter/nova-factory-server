package models

import "nova-factory-server/app/baize"

// ShopSysConfig 商城系统配置（通用）
type ShopSysConfig struct {
	ID          int64  `json:"id,string" gorm:"id"`
	ConfigKey   string `json:"configKey" gorm:"config_key"`
	ConfigValue string `json:"configValue" gorm:"config_value"`
	ConfigType  string `json:"configType" gorm:"config_type"`
	Remark      string `json:"remark" gorm:"remark"`
	DeptID      int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// KeyValue key-value 结构，用于泛化配置接口
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// BatchConfigReq 通用批量配置请求（全量覆盖）
type BatchConfigReq struct {
	Configs []KeyValue `json:"configs"`
}

// BatchConfigResp 通用批量配置响应
type BatchConfigResp struct {
	Configs []KeyValue `json:"configs"`
}
