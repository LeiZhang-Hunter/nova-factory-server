package dashboardModels

import (
	"nova-factory-server/app/baize"
)

// SysDashboard mapped from table <sys_dashboard>
type SysDashboard struct {
	ID          int64  `gorm:"column:id;primaryKey;comment:调度id" json:"id"` // 调度id
	Name        string `gorm:"column:name;not null" json:"name"`
	Type        string `gorm:"column:type;not null" json:"type"`
	Description string `gorm:"column:description;not null;comment:描述" json:"description"` // 描述
	DeptID      int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SetSysDashboard struct {
	ID          int64  `gorm:"column:id;primaryKey;comment:调度id"  json:"id"` // 调度id
	Name        string `gorm:"column:name;not null"  binding:"required" json:"name"`
	Type        string `gorm:"column:type;not null"  binding:"required" json:"type"`
	Description string `gorm:"column:description;not null;comment:描述"  binding:"required" json:"description"` // 描述
}

type SysDashboardReq struct {
	Name string `gorm:"column:name;not null" json:"name"`
	Type string `gorm:"column:type;not null" json:"type"`
	baize.BaseEntityDQL
}

func ToSysDashboardReq(set *SetSysDashboard) *SysDashboard {
	return &SysDashboard{
		ID:          set.ID,
		Name:        set.Name,
		Type:        set.Type,
		Description: set.Description,
	}
}

type SysDashboardList struct {
	Rows  []*SysDashboard `json:"rows"`
	Total int64           `json:"total"`
}
