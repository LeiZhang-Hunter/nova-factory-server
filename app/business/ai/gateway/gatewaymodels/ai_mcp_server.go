package gatewaymodels

import "nova-factory-server/app/baize"

type MCPServer struct {
	ID          string `gorm:"column:id;primaryKey" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	Transport   string `gorm:"column:transport" json:"transport"`
	Command     string `gorm:"column:command" json:"command"`
	Args        string `gorm:"column:args" json:"args"`
	Env         string `gorm:"column:env" json:"env"`
	URL         string `gorm:"column:url" json:"url"`
	Headers     string `gorm:"column:headers" json:"headers"`
	Timeout     int32  `gorm:"column:timeout" json:"timeout"`
	IsCommon    *bool  `gorm:"column:is_common" json:"isCommon"`
	Enabled     *bool  `gorm:"column:enabled" json:"enabled"`
	DeptID      int64  `gorm:"column:dept_id" json:"deptId"`
	baize.BaseEntity
	State int32 `gorm:"column:state" json:"state"`
}

type MCPServerQuery struct {
	Name      string `form:"name"`
	Transport string `form:"transport"`
	Enabled   *bool  `form:"enabled"`
	Page      int64  `form:"page"`
	Size      int64  `form:"size"`
}

type MCPServerUpsert struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Transport   string `json:"transport"`
	Command     string `json:"command"`
	Args        string `json:"args"`
	Env         string `json:"env"`
	URL         string `json:"url"`
	Headers     string `json:"headers"`
	Timeout     int32  `json:"timeout"`
	IsCommon    *bool  `json:"isCommon"`
	Enabled     *bool  `json:"enabled"`
}

type MCPServerListData struct {
	Rows  []*MCPServer `json:"rows"`
	Total int64        `json:"total"`
}
