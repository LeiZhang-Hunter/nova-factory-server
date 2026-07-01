package gatewaymodels

import "nova-factory-server/app/baize"

// AgentConfigKey API Key 模型。
type AgentConfigKey struct {
	ID                  int64  `json:"id,string" gorm:"column:id"`
	Key                 string `json:"key" gorm:"column:key"`
	AllowMcpServerTools string `json:"-" gorm:"column:allow_mcp_server_tools"`
	DeptID              int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// AgentConfigKeyQuery API Key 查询参数。
type AgentConfigKeyQuery struct {
	Key  string `form:"key"`
	Page int64  `form:"page"`
	Size int64  `form:"size"`
}

// AgentConfigKeyUpsert API Key 保存参数。
type AgentConfigKeyUpsert struct {
	ID  int64  `json:"id,string"`
	Key string `json:"key"`
}

// AgentConfigKeyListData API Key 列表数据。
type AgentConfigKeyListData struct {
	Rows  []*AgentConfigKey `json:"rows"`
	Total int64             `json:"total"`
}

// AgentConfigKeyToolUpsert API Tool 保存参数。
type AgentConfigKeyToolUpsert struct {
	ID    int64    `json:"id,string"`
	Tools []string `json:"tools"`
}
