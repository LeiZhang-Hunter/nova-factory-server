package dashboardModels

import (
	"nova-factory-server/app/baize"
)

// SysDashboardData mapped from table <sys_dashboard_data>
type SysDashboardData struct {
	ID            int64  `gorm:"column:id;primaryKey;comment:面板id" json:"id"`                          // 面板id
	DatashboardID int64  `gorm:"column:datashboard_id;not null;comment:仪表盘数据id" json:"datashboard_id"` // 仪表盘数据id
	Data          string `gorm:"column:data;comment:数据" json:"data"`                                   // 数据
	DeptID        int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                           // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SetSysDashboardData struct {
	ID            int64  `gorm:"column:id;primaryKey;comment:面板id" json:"id"`                          // 面板id
	DatashboardID int64  `gorm:"column:datashboard_id;not null;comment:仪表盘数据id" json:"datashboard_id"` // 仪表盘数据id
	Data          string `gorm:"column:data;comment:数据" json:"data"`                                   // 数据
}

func ToSysDashboardData(set *SetSysDashboardData) *SysDashboardData {
	return &SysDashboardData{
		ID:            set.ID,
		DatashboardID: set.DatashboardID,
		Data:          set.Data,
	}
}

type SysDashboardDataReq struct {
	baize.BaseEntityDQL
}

type SysDashboardDataList struct {
	Rows  []*SysDashboardData `json:"rows"`
	Total int64               `json:"total"`
}
