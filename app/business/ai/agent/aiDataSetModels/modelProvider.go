package aiDataSetModels

import "nova-factory-server/app/baize"

type SysAiModelProvider struct {
	ID         int64       `json:"id" db:"id"`
	Name       string      `json:"name" db:"name"`
	Logo       string      `json:"logo" db:"logo"`
	Tags       string      `json:"tags" db:"tags"`
	Status     int32       `json:"status" db:"status"`
	Rank       int32       `json:"rank" db:"rank"`
	DeptID     *int64      `json:"deptId" db:"dept_id"`
	LLMs       []*SysAiLLM `json:"llms" gorm:"-"`
	BaseEntity baize.BaseEntity
	State      int32 `json:"state" db:"state"`
}

type SysAiLLM struct {
	LlmName   string `json:"llmName" db:"llm_name"`
	ModelType string `json:"modelType" db:"model_type"`
	Fid       string `json:"fid" db:"fid"`
	MaxTokens int32  `json:"maxTokens" db:"max_tokens"`
	Tags      string `json:"tags" db:"tags"`
	IsTools   int32  `json:"isTools" db:"is_tools"`
	Status    string `json:"status" db:"status"`
	DeptID    *int64 `json:"deptId" db:"dept_id"`
	State     int32  `json:"state" db:"state"`
}

type SysAiModelProviderListReq struct {
	Name   string `form:"name"`
	Status *int32 `form:"status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

type SysAiModelProviderListData struct {
	Rows  []*SysAiModelProvider `json:"rows"`
	Total int64                 `json:"total"`
}
