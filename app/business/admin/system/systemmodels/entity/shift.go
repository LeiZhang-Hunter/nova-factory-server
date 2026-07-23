package entity

import "nova-factory-server/app/baize"

// SysWorkShiftSetting 班次设置
type SysWorkShiftSetting struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id,string"`
	Name         string `gorm:"column:name;not null;comment:班次名称" json:"name"`                        // 班次名称
	BeginTime    int32  `gorm:"column:begin_time;not null;comment:开始时间" json:"begin_time"`            // 开始时间
	BeginTimeStr string `gorm:"column:begin_time_str;not null;comment:开始时间字符串" json:"begin_time_str"` // 开始时间字符串
	EndTime      int32  `gorm:"column:end_time;not null;comment:结束时间" json:"end_time"`                // 结束时间
	EndTimeStr   string `gorm:"column:end_time_str;not null;comment:结束时间字符串" json:"end_time_str"`     // 结束时间字符串
	Status       *bool  `gorm:"column:status;not null;default:1;comment:是否启用班次设置" json:"status"`      // 是否启用班次设置
	DeptID       int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                           // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;not null;default:0" json:"state"`
}
