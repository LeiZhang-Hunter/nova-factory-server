package gatewaymodels

import "nova-factory-server/app/baize"

// AIAgentConfigPublishHistory 智能体配置发布历史。
type AIAgentConfigPublishHistory struct {
	ID                 int64  `json:"id,string" gorm:"column:id"`
	AgentID            int64  `json:"agentId,string" gorm:"column:agent_id"`
	Version            string `json:"version" gorm:"column:version"`
	ConfigMd5          string `json:"-"  gorm:"column:config_md5"`
	ConfigSnapshot     string `json:"configSnapshot" gorm:"column:config_snapshot"`
	PublishDescription string `json:"publishDescription" gorm:"column:publish_description"`
	DeptID             int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// AIAgentConfigPublishHistoryUpsert 智能体配置发布历史新增修改参数。
type AIAgentConfigPublishHistoryUpsert struct {
	ID                 int64  `json:"id,string"`
	AgentID            int64  `json:"agentId,string"`
	Version            string `json:"version"`
	ConfigSnapshot     string `json:"configSnapshot"`
	ConfigMd5          string `json:"-"`
	PublishDescription string `json:"publishDescription"`
}

// AIAgentConfigPublishHistoryQuery 智能体配置发布历史查询参数。
type AIAgentConfigPublishHistoryQuery struct {
	AgentID int64  `form:"agentId"`
	Version string `form:"version"`
	Page    int64  `form:"page"`
	Size    int64  `form:"size"`
}

// AIAgentConfigPublishHistoryListData 智能体配置发布历史列表结果。
type AIAgentConfigPublishHistoryListData struct {
	Rows  []*AIAgentConfigPublishHistory `json:"rows"`
	Total int64                          `json:"total"`
}
