package models

import "nova-factory-server/app/baize"

// ShopSysConfig 商城系统配置
type ShopSysConfig struct {
	ID          int64  `json:"id,string" db:"id"`                   // 主键ID
	ConfigKey   string `json:"configKey" db:"config_key"`           // 配置键名
	ConfigValue string `json:"configValue" db:"config_value"`       // 配置键值
	ConfigType  string `json:"configType" db:"config_type"`         // 配置类型
	Remark      string `json:"remark" db:"remark"`                  // 备注
	DeptID      int64  `json:"deptId" gorm:"column:dept_id" db:"-"` // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state" db:"state"` // 操作状态
}
