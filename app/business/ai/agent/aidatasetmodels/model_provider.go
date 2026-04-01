package aidatasetmodels

import (
	"nova-factory-server/app/baize"
)

type SysAiModelProvider struct {
	ID     int64       `json:"id" db:"id"`
	Name   string      `json:"name" db:"name"`
	Logo   string      `json:"logo" db:"logo"`
	Tags   string      `json:"tags" db:"tags"`
	Status int32       `json:"status" db:"status"`
	Rank   int32       `json:"rank" db:"rank"`
	DeptID int64       `json:"deptId" db:"dept_id"`
	LLMs   []*SysAiLLM `json:"llms" gorm:"-"`
	baize.BaseEntity
	State int32 `json:"state" db:"state"`
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

type FactoryProviderUpsert struct {
	Name string `json:"name"`
	Logo string `json:"logo"`
	Tags string `json:"tags"`
}

type FactoryLLMUpsert struct {
	LLMName   string `json:"llm_name"`
	Tags      string `json:"tags"`
	MaxTokens int64  `json:"max_tokens"`
	ModelType string `json:"model_type"`
	IsTools   bool   `json:"is_tools"`
}

type AiModelProviderEntity struct {
	ID     int64  `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	Logo   string `gorm:"column:logo"`
	Tags   string `gorm:"column:tags"`
	Status int32  `gorm:"column:status"`
	Rank   int32  `gorm:"column:rank"`
	DeptID int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"` // 部门ID
	baize.BaseEntity
	State int32 `gorm:"column:state"`
}

type AiLLMEntity struct {
	FID       string `gorm:"column:fid"`
	LlmName   string `gorm:"column:llm_name"`
	ModelType string `gorm:"column:model_type"`
	MaxTokens int64  `gorm:"column:max_tokens"`
	Tags      string `gorm:"column:tags"`
	IsTools   int32  `gorm:"column:is_tools"`
	Status    string `gorm:"column:status"`
	State     int32  `gorm:"column:state"`
	DeptID    int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"` // 部门ID
	baize.BaseEntity
}
