package craftRouteModels

import (
	"nova-factory-server/app/baize"
	"time"
)

const (
	DAILY   = 1
	SPECIAL = 2
)

type SysProductSchedule struct {
	ID           int64  `gorm:"column:id;primaryKey;comment:调度id" json:"id"`                                                  // 调度id
	GatewayID    int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`                                    // 网关id
	ScheduleName string `gorm:"column:schedule_name;not null;comment:计划名称" json:"schedule_name"`                              // 计划名称
	Time         string `gorm:"column:time;not null;comment:时间序列化格式,普通日程,1,2,3,4,5;特殊日程:2025-04-04 ~ 2025-04-04" json:"time"` // 时间序列化格式,普通日程,1,2,3,4,5;特殊日程:2025-04-04 ~ 2025-04-04
	ScheduleType int    `gorm:"column:schedule_type;not null;comment:0为普通日程 1为特殊日程" json:"schedule_type"`                     // 0为普通日程 1为特殊日程
	Status       bool   `gorm:"column:status;not null;comment:操作状态（0正常 1启动）" json:"status"`                                   // 操作状态（0正常 1启动）
	DeptID       int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                                   // 部门ID
	State        bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                             // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func ToSysProductSchedule(data *SetSysProductSchedule) *SysProductSchedule {
	return &SysProductSchedule{
		ID:           data.Id,
		GatewayID:    data.GatewayID,
		ScheduleName: data.ScheduleName,
		Time:         data.Time,
		ScheduleType: data.Type,
	}
}

type SysProductScheduleReq struct {
	Year  string `form:"year" json:"year" binding:"required"`
	Month string `form:"month" json:"month" binding:"required"`
}

type ScheduleStatusData struct {
	Time   time.Time
	Type   int
	Status string
}

type SysProductScheduleListReq struct {
	Name string `form:"name" json:"name"`
	baize.BaseEntityDQL
}

type SysProductScheduleListData struct {
	Rows  []*SysProductSchedule `json:"rows"`
	Total int64                 `json:"total"`
}

type TimeManagerData struct {
	BeginTime string `json:"begin_time"`
	EndTime   string `json:"end_time"`
	RoueId    int64  `json:"run_status"`
}

type SetSysProductSchedule struct {
	Id           int64              `json:"id,string" binding:"required"`
	GatewayID    int64              `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"` // 网关id
	ScheduleName string             `json:"schedule_name" binding:"required"`
	Time         string             `json:"time" binding:"required"`
	TimeManager  []*TimeManagerData `json:"time_manager" binding:"required"`
	Type         int                `json:"type" binding:"required"`
}
