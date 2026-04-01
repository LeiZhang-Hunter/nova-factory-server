package gatewaymodels

import "nova-factory-server/app/baize"

type AIGateway struct {
	ID      int64  `json:"id" gorm:"column:id"`
	Name    string `json:"name" gorm:"column:name"`
	BaseURL string `json:"baseUrl" gorm:"column:base_url"`
	APIKey  string `json:"apiKey" gorm:"column:api_key"`
	Enabled int32  `json:"enabled" gorm:"column:enabled"`
	Active  int32  `json:"active" gorm:"column:active"`
	DeptID  int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

type AIGatewayQuery struct {
	Name    string `form:"name"`
	Enabled *int32 `form:"enabled"`
	Active  *int32 `form:"active"`
	Page    int64  `form:"page"`
	Size    int64  `form:"size"`
}

type AIGatewayUpsert struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BaseURL string `json:"baseUrl"`
	APIKey  string `json:"apiKey"`
	Enabled int32  `json:"enabled"`
	Active  int32  `json:"active"`
}

type AIGatewayListData struct {
	Rows  []*AIGateway `json:"rows"`
	Total int64        `json:"total"`
}
