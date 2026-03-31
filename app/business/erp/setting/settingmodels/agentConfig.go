package settingmodels

import (
	"nova-factory-server/app/baize"
)

type AgentConfig struct {
	ID      uint64 `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	AgentID string `json:"agentId" db:"agent_id"`
	Remark  string `json:"remark" db:"remark"`
	Status  *bool  `json:"status" db:"status"`
	DeptID  int64  `json:"deptId" db:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" db:"state"`
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
