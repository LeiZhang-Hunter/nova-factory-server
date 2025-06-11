package daemonizeModels

import (
	"nova-factory-server/app/baize"
)

type SysIotAgentConfig struct {
	ID            uint64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"`          // 自增标识
	AgentObjectID int64  `gorm:"column:agent_object_id;not null;comment:agent id" json:"agent_object_id"` // agent id
	ConfigVersion string `gorm:"column:config_version;not null;comment:配置版本" json:"config_version"`       // 配置版本
	Content       string `gorm:"column:content;not null;comment:配置内容" json:"content"`                     // 配置内容
	DeptID        int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                              // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func OfSysIotAgentConfig(set *SysIotAgentConfigSetReq) *SysIotAgentConfig {
	return &SysIotAgentConfig{
		AgentObjectID: set.AgentObjectID,
		ConfigVersion: set.ConfigVersion,
		Content:       set.Content,
	}
}

type SysIotAgentConfigSetReq struct {
	ID            uint64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"`          // 自增标识
	AgentObjectID int64  `gorm:"column:agent_object_id;not null;comment:agent id" json:"agent_object_id"` // agent id
	ConfigVersion string `gorm:"column:config_version;not null;comment:配置版本" json:"config_version"`       // 配置版本
	Content       string `gorm:"column:content;not null;comment:配置内容" json:"content"`                     // 配置内容
}

type GenerateGatewayConfigReq struct {
	ObjectID uint64 `gorm:"column:object_id;primaryKey;comment:agent uuid" json:"object_id"` // agent uuid
}

type GenerateGatewayConfigRes struct {
	Data string `json:"data,omitempty"` // 配置 uuid
}

type SysIotAgentConfigListReq struct {
	AgentObjectID int64 `gorm:"column:agent_object_id;not null;comment:agent id" json:"agent_object_id"` // agent id
	baize.BaseEntityDQL
}

// SysIotAgentConfigListData Agent配置列表
type SysIotAgentConfigListData struct {
	Rows  []*SysIotAgentConfig `json:"rows"`
	Total int64                `json:"total"`
}
