package settingmodels

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/kind"
)

// IntegrationConfig 集成配置
type IntegrationConfig struct {
	ID     uint64 `json:"id" gorm:"id"`
	Type   string `json:"type" gorm:"type"`
	Data   string `json:"data" gorm:"data"`
	Status *bool  `json:"status" gorm:"status"`
	DeptID int64  `json:"deptId" gorm:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"`
}

func (i *IntegrationConfig) GetData() string {
	return i.Data
}

// GetType 返回集成类型，实现 observer.IntegrationConfig 接口
func (i *IntegrationConfig) GetType() string {
	return i.Type
}

func (i *IntegrationConfig) GetStatus() *bool {
	return i.Status
}
func (i *IntegrationConfig) GetOverrideURL() string {
	return ""
}
func (i *IntegrationConfig) GetMetadata() map[string]any {
	return make(map[string]any)
}

func (i *IntegrationConfig) Service() (api.Service, error) {
	return api.GetRegistry().Create(kind.Kind(i.Type))
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
	CheckURL    string `form:"checkUrl"`
	RedirectURL string `form:"redirectUrl"`
}

type IntegrationOAuthCallbackReq struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state"`
	Type  string `form:"type"`
}

type IntegrationOAuthCallbackData struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	Token      string `json:"token"`
	ExpireDate string `json:"expireDate"`
	IssueDate  string `json:"issueDate"`
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
	Message    string `json:"message"`
	ApiCode    int64  `json:"apiCode"`
}
