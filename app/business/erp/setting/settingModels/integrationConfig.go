package settingModels

import "nova-factory-server/app/baize"

type IntegrationConfig struct {
	ID     uint64 `json:"id" db:"id"`
	Type   string `json:"type" db:"type"`
	Data   string `json:"data" db:"data"`
	Status *bool  `json:"status" db:"status"`
	DeptID int64  `json:"deptId" db:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" db:"state"`
}

type IntegrationConfigSet struct {
	Type   string `json:"type" binding:"required"`
	Data   string `json:"data"`
	Status *bool  `json:"status"`
}

type IntegrationConfigQuery struct {
	Type   string `form:"type"`
	Status *bool  `form:"status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

type IntegrationConfigListData struct {
	Rows  []*IntegrationConfig `json:"rows"`
	Total int64                `json:"total"`
}

type IntegrationConfigCheckLoginReq struct {
	Type        string `form:"type" binding:"required"`
	CheckURL    string `form:"checkUrl"`
	RedirectURL string `form:"redirectUrl"`
}

type IntegrationOAuthCallbackReq struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state"`
}

type IntegrationOAuthCallbackData struct {
	Code  string `json:"code"`
	State string `json:"state"`
}
