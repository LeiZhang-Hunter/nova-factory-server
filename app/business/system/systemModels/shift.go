package systemModels

import (
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/time"
)

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

func ToSysWorkShiftSetting(vo *SysWorkShiftSettingVO) (*SysWorkShiftSetting, error) {
	begin, err := time.FormatTodayTIme(vo.BeginTimeStr)
	if err != nil {
		zap.L().Error("ToSysWorkShiftSetting FormatTodayTIme error", zap.Error(err))
		return nil, err
	}
	end, err := time.FormatTodayTIme(vo.EndTimeStr)
	if err != nil {
		zap.L().Error("ToSysWorkShiftSetting FormatTodayTIme error", zap.Error(err))
		return nil, err
	}
	// 隔日第二天
	if end <= begin {
		end = end + 86400
	}
	return &SysWorkShiftSetting{
		ID:           vo.ID,
		Name:         vo.Name,
		BeginTime:    begin,
		BeginTimeStr: vo.BeginTimeStr,
		EndTime:      end,
		EndTimeStr:   vo.EndTimeStr,
		Status:       &vo.Status,
	}, nil
}

// SysWorkShiftSettingVO 班次设置
type SysWorkShiftSettingVO struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id,string"`
	Name         string `gorm:"column:name;not null;comment:班次名称" json:"name"`                        // 班次名称
	BeginTime    int32  `gorm:"column:begin_time;not null;comment:开始时间" json:"begin_time"`            // 开始时间
	BeginTimeStr string `gorm:"column:begin_time_str;not null;comment:开始时间字符串" json:"begin_time_str"` // 开始时间字符串
	EndTime      int32  `gorm:"column:end_time;not null;comment:结束时间" json:"end_time"`                // 结束时间
	EndTimeStr   string `gorm:"column:end_time_str;not null;comment:结束时间字符串" json:"end_time_str"`     // 结束时间字符串
	Status       bool   `gorm:"column:status;not null;default:1;comment:是否启用班次设置" json:"status"`      // 是否启用班次设置
}

// SysWorkShiftSettingReq 班次配置列表
type SysWorkShiftSettingReq struct {
	Name string `form:"name" ` //排序规则  降序desc   asc升序
	baize.BaseEntityDQL
}

type SysWorkShiftSettingList struct {
	Rows  []*SysWorkShiftSetting `json:"rows"`
	Total int64                  `json:"total"`
}
