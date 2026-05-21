package gatewaymodels

import "nova-factory-server/app/baize"

// AISubAgent 子智能体配置。
type AISubAgent struct {
	ID                        int64               `json:"id,string" gorm:"column:id"`
	Name                      string              `json:"name" gorm:"column:name"`
	Type                      string              `json:"type" gorm:"column:type"`
	Description               string              `json:"description" gorm:"column:description"`
	Instruction               string              `json:"instruction" gorm:"column:instruction"`
	MCPEnabled                *bool               `json:"mcpEnabled" gorm:"column:mcp_enabled"`
	MCPServerIDs              string              `json:"mcpServerIds" gorm:"column:mcp_server_ids"`
	MCPServerEnabledIDs       string              `json:"mcpServerEnabledIds" gorm:"column:mcp_server_enabled_ids"`
	LocalToolEnabled          *bool               `json:"localToolEnabled" gorm:"column:local_tool_enabled"`
	LocalTools                string              `json:"localTools" gorm:"column:local_tools"`
	AllowMcpServerIdsToolsRaw string              `json:"-" gorm:"column:allow_mcp_server_ids_tools"`
	AllowMcpServerIdsTools    map[string][]string `json:"allowMcpServerIdsTools" gorm:"-"`
	Enable                    *bool               `json:"enable" gorm:"column:enable"`
	DeptID                    int64               `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// AISubAgentQuery 子智能体查询参数。
type AISubAgentQuery struct {
	Name             string `form:"name"`
	Type             string `form:"type"`
	MCPEnabled       *bool  `form:"mcpEnabled"`
	LocalToolEnabled *bool  `form:"localToolEnabled"`
	Enable           *bool  `form:"enable"`
	Page             int64  `form:"page"`
	Size             int64  `form:"size"`
}

// AISubAgentUpsert 子智能体新增修改参数。
type AISubAgentUpsert struct {
	ID                        int64               `json:"id,string"`
	Name                      string              `json:"name"`
	Type                      string              `json:"type"`
	Description               string              `json:"description"`
	Instruction               string              `json:"instruction"`
	MCPEnabled                *bool               `json:"mcpEnabled"`
	MCPServerIDs              string              `json:"mcpServerIds"`
	MCPServerEnabledIDs       string              `json:"mcpServerEnabledIds"`
	LocalToolEnabled          *bool               `json:"localToolEnabled"`
	LocalTools                string              `json:"localTools"`
	AllowMcpServerIdsToolsRaw string              `json:"-"`
	AllowMcpServerIdsTools    map[string][]string `json:"allowMcpServerIdsTools"`
	Enable                    *bool               `json:"enable"`
}

// AISubAgentListData 子智能体列表结果。
type AISubAgentListData struct {
	Rows  []*AISubAgent `json:"rows"`
	Total int64         `json:"total"`
}
