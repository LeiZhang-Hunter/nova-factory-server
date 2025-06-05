package daemonizeModels

import (
	"nova-factory-server/app/baize"
)

type SysIotAgentConfig struct {
	ID            int32  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"`    // 自增标识
	UUID          string `gorm:"column:uuid;not null;comment:配置 uuid" json:"uuid"`                  // 配置 uuid
	CompanyUUID   string `gorm:"column:company_uuid;not null;comment:公司 uuid" json:"company_uuid"`  // 公司 uuid
	ConfigVersion string `gorm:"column:config_version;not null;comment:配置版本" json:"config_version"` // 配置版本
	Content       string `gorm:"column:content;not null;comment:配置内容" json:"content"`               // 配置内容
	DeptID        int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                        // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}
