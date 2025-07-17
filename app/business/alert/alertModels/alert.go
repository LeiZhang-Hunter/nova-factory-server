package alertModels

import (
	"nova-factory-server/app/baize"
)

// SysAlert 告警策略配置
type SysAlert struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"` // 自增标识
	GatewayID   int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`      // 网关id
	TemplateID  int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id"`    // 模板id
	Name        string `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                // 告警策略名称
	Additions   string `gorm:"column:additions;comment:注解" json:"additions"`                   // 注解
	Advanced    string `gorm:"column:advanced;comment:告警规则" json:"advanced"`                   // 告警规则
	Description string `gorm:"column:description;not null;comment:配置版本" json:"description"`    // 配置版本
	DeptID      int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                     // 部门ID
	baize.BaseEntity
	State  bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`  // 操作状态（0正常 -1删除）
	Status bool `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"` // 操作状态（0正常 1异常）
}

type SetSysAlert struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"` // 自增标识
	GatewayID   int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`      // 网关id
	TemplateID  int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id"`    // 模板id
	Name        string `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                // 告警策略名称
	Additions   string `gorm:"column:additions;comment:注解" json:"additions"`                   // 注解
	Advanced    string `gorm:"column:advanced;comment:告警规则" json:"advanced"`                   // 告警规则
	Description string `gorm:"column:description;not null;comment:配置版本" json:"description"`    // 配置版本
	Status      bool   `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`              // 操作状态（0正常 1异常）
}

func ToSysAlert(data *SetSysAlert) *SysAlert {
	return &SysAlert{
		ID:          data.ID,
		GatewayID:   data.GatewayID,
		TemplateID:  data.TemplateID,
		Name:        data.Name,
		Additions:   data.Additions,
		Advanced:    data.Advanced,
		Description: data.Description,
		Status:      data.Status,
	}
}

type SysAlertListReq struct {
	GatewayID  int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`   // 网关id
	TemplateID int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id"` // 模板id
	Name       string `gorm:"column:name;not null;comment:告警策略名称" json:"name"`             // 告警策略名称
	Status     *bool  `gorm:"column:status;comment:启用状态" json:"status"`
	baize.BaseEntityDQL
}

type SysAlertList struct {
	Rows  []*SysAlert `json:"rows"`
	Total uint64      `json:"total"`
}

type ChangeSysAlert struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"` // 自增标识
	Status bool  `gorm:"column:status;comment:启用状态" json:"status"`
}
