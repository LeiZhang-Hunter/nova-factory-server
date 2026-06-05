package gatewaymodels

import "nova-factory-server/app/baize"

// AIAgent 智能体配置
type AIAgent struct {
	ID                        int64               `json:"id,string" gorm:"column:id"`
	Name                      string              `json:"name" gorm:"column:name"`
	Type                      string              `json:"type"`
	Prompt                    string              `json:"prompt" gorm:"column:prompt"`
	DefaultLLMProviderID      string              `json:"defaultLlmProviderId" gorm:"column:default_llm_provider_id"`
	DefaultLLMModelID         string              `json:"defaultLlmModelId" gorm:"column:default_llm_model_id"`
	LLMTemperature            float64             `json:"llmTemperature" gorm:"column:llm_temperature"`
	LLMTopP                   float64             `json:"llmTopP" gorm:"column:llm_top_p"`
	LLMMaxTokens              int32               `json:"llmMaxTokens" gorm:"column:llm_max_tokens"`
	EnableLLMTemperature      *bool               `json:"enableLlmTemperature" gorm:"column:enable_llm_temperature"`
	EnableLLMTopP             *bool               `json:"enableLlmTopP" gorm:"column:enable_llm_top_p"`
	EnableLLMMaxTokens        *bool               `json:"enableLlmMaxTokens" gorm:"column:enable_llm_max_tokens"`
	LLMMaxContextCount        int32               `json:"llmMaxContextCount" gorm:"column:llm_max_context_count"`
	SandboxMode               string              `json:"sandboxMode" gorm:"column:sandbox_mode"`
	SandboxNetwork            *bool               `json:"sandboxNetwork" gorm:"column:sandbox_network"`
	WorkDir                   string              `json:"workDir" gorm:"column:work_dir"`
	MCPEnabled                *bool               `json:"mcpEnabled" gorm:"column:mcp_enabled"`
	SuppressPreToolContent    *bool               `json:"suppressPreToolContent" gorm:"column:suppress_pre_tool_content"` // Session-level strategy: buffer pre-tool assistant content until tool usage is known
	MCPServerIDs              string              `json:"mcpServerIds" gorm:"column:mcp_server_ids"`
	MCPServerEnabledIDs       string              `json:"mcpServerEnabledIds" gorm:"column:mcp_server_enabled_ids"`
	AllowMcpServerIdsToolsRaw string              `json:"-" gorm:"column:allow_mcp_server_ids_tools"`
	AllowMcpServerIdsTools    map[string][]string `json:"allowMcpServerIdsTools" gorm:"-"`
	ConfigVersion             string              `json:"configVersion" gorm:"column:config_version"`
	ActiveVersion             string              `json:"activeVersion" gorm:"-"`
	Enable                    *bool               `json:"enable"`
	DeptID                    int64               `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// AIAgentQuery 智能体查询参数
type AIAgentQuery struct {
	Name        string `form:"name"`
	SandboxMode string `form:"sandboxMode"`
	MCPEnabled  *bool  `form:"mcpEnabled"`
	Page        int64  `form:"page"`
	Size        int64  `form:"size"`
}

// AIAgentUpsert 智能体新增修改参数
type AIAgentUpsert struct {
	ID                        int64               `json:"id,string"`
	Name                      string              `json:"name"`
	Type                      string              `json:"type"`
	Prompt                    string              `json:"prompt"`
	DefaultLLMProviderID      string              `json:"defaultLlmProviderId"`
	DefaultLLMModelID         string              `json:"defaultLlmModelId"`
	LLMTemperature            float64             `json:"llmTemperature"`
	LLMTopP                   float64             `json:"llmTopP"`
	LLMMaxTokens              int32               `json:"llmMaxTokens"`
	EnableLLMTemperature      *bool               `json:"enableLlmTemperature"`
	EnableLLMTopP             *bool               `json:"enableLlmTopP"`
	EnableLLMMaxTokens        *bool               `json:"enableLlmMaxTokens"`
	LLMMaxContextCount        int32               `json:"llmMaxContextCount"`
	RetrievalTopK             int32               `json:"retrievalTopK"`
	RetrievalMatchThreshold   float64             `json:"retrievalMatchThreshold"`
	SandboxMode               string              `json:"sandboxMode"`
	SandboxNetwork            *bool               `json:"sandboxNetwork"`
	WorkDir                   string              `json:"workDir"`
	MCPEnabled                *bool               `json:"mcpEnabled"`
	MCPServerIDs              string              `json:"mcpServerIds"`
	MCPServerEnabledIDs       string              `json:"mcpServerEnabledIds"`
	AllowMcpServerIdsToolsRaw string              `json:"-"`
	Enable                    *bool               `json:"enable"`
	SuppressPreToolContent    *bool               `json:"suppressPreToolContent"`
	AllowMcpServerIdsTools    map[string][]string `json:"allowMcpServerIdsTools"` //允许描述的工具 mcp server id => tools
}

// AIAgentListData 智能体列表结果
type AIAgentListData struct {
	Rows  []*AIAgent `json:"rows"`
	Total int64      `json:"total"`
}
