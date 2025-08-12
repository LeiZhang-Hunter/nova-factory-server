package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

type SysProductScheduleMap struct {
	ID           int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ScheduleID   int64 `gorm:"column:schedule_id;not null;comment:计划标识" json:"schedule_id"`               // 计划标识
	BeginTime    int64 `gorm:"column:begin_time;not null;comment:开始时间" json:"begin_time"`                 // 开始时间
	EndTime      int64 `gorm:"column:end_time;not null;comment:结束时间" json:"end_time"`                     // 结束时间
	Date         int   `gorm:"column:date;not null;comment:日期 1 2 3 4 5 6 7分别代表周几" json:"date"`           // 日期 1 2 3 4 5 6 7分别代表周几
	CraftRouteID int64 `gorm:"column:craft_route_id;not null;comment:工艺路线id" json:"craft_route_id"`       // 工艺路线id
	ScheduleType int   `gorm:"column:schedule_type;not null;comment:0为 循环日程 1为特殊日程" json:"schedule_type"` // 0为 循环日程 1为特殊日程
	State        bool  `gorm:"column:state;not null;default:0" json:"state"`
	DeptID       int64 `gorm:"column:dept_id;comment:部门ID" json:"dept_id"` // 部门ID
	baize.BaseEntity
}
