package settingmodels

import (
	"nova-factory-server/app/baize"
)

type AgentConfig struct {
	ID      uint64 `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	AgentID string `json:"agentId" gorm:"agent_id"`
	Remark  string `json:"remark" gorm:"remark"`
	Status  *bool  `json:"status" gorm:"status"`
	DeptID  int64  `json:"deptId" gorm:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"`
}

type AgentConfigUpsert struct {
	ID      uint64 `json:"id,string"`
	Name    string `json:"name" binding:"required"`
	AgentID string `json:"agentId" binding:"required"`
	Remark  string `json:"remark"`
	Status  *bool  `json:"status"`
}

type AgentConfigQuery struct {
	Name    string `form:"name"`
	AgentID string `form:"agentId"`
	Status  *bool  `form:"status"`
	Page    int64  `form:"page"`
	Size    int64  `form:"size"`
}

type AgentConfigListData struct {
	Rows  []*AgentConfig `json:"rows"`
	Total int64          `json:"total"`
}
