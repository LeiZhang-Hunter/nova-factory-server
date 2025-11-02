package deviceModels

import (
	"nova-factory-server/app/baize"
	"time"
)

// SysDeviceCheckPlan 设备点检保养计划头表
type SysDeviceCheckPlan struct {
	PlanID     int64     `gorm:"column:plan_id;primaryKey;autoIncrement:true;comment:计划ID" json:"plan_id"` // 计划ID
	PlanCode   string    `gorm:"column:plan_code;not null;comment:计划编码" json:"plan_code"`                  // 计划编码
	PlanName   string    `gorm:"column:plan_name;comment:计划名称" json:"plan_name"`                           // 计划名称
	PlanType   string    `gorm:"column:plan_type;not null;comment:计划类型" json:"plan_type"`                  // 计划类型
	StartDate  time.Time `gorm:"column:start_date;comment:开始日期" json:"start_date"`                         // 开始日期
	EndDate    time.Time `gorm:"column:end_date;comment:结束日期" json:"end_date"`                             // 结束日期
	CycleType  string    `gorm:"column:cycle_type;comment:频率" json:"cycle_type"`                           // 频率
	CycleCount int32     `gorm:"column:cycle_count;comment:次数" json:"cycle_count"`                         // 次数
	Status     string    `gorm:"column:status;comment:状态" json:"status"`                                   // 状态
	Remark     string    `gorm:"column:remark;comment:备注" json:"remark"`                                   // 备注
	Attr1      string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                  // 预留字段1
	Attr2      string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                  // 预留字段2
	Attr3      int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                  // 预留字段3
	Attr4      int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                  // 预留字段4
	DeptID     int64     `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                               // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func ToSysDeviceCheckPlan(vo *SysDeviceCheckPlanVO) *SysDeviceCheckPlan {
	return &SysDeviceCheckPlan{
		PlanID:     vo.PlanID,
		PlanCode:   vo.PlanCode,
		PlanName:   vo.PlanName,
		PlanType:   vo.PlanType,
		StartDate:  vo.StartDate,
		EndDate:    vo.EndDate,
		CycleType:  vo.CycleType,
		CycleCount: vo.CycleCount,
		Status:     vo.Status,
		Remark:     vo.Remark,
		Attr1:      vo.Attr1,
		Attr2:      vo.Attr2,
		Attr3:      vo.Attr3,
		Attr4:      vo.Attr4,
	}
}

type SysDeviceCheckPlanVO struct {
	PlanID     int64     `gorm:"column:plan_id;primaryKey;autoIncrement:true;comment:计划ID" json:"plan_id,string"` // 计划ID
	PlanCode   string    `gorm:"column:plan_code;not null;comment:计划编码" json:"plan_code"`                         // 计划编码
	PlanName   string    `gorm:"column:plan_name;comment:计划名称" json:"plan_name"`                                  // 计划名称
	PlanType   string    `gorm:"column:plan_type;not null;comment:计划类型" json:"plan_type"`                         // 计划类型
	StartDate  time.Time `gorm:"column:start_date;comment:开始日期" json:"start_date"`                                // 开始日期
	EndDate    time.Time `gorm:"column:end_date;comment:结束日期" json:"end_date"`                                    // 结束日期
	CycleType  string    `gorm:"column:cycle_type;comment:频率" json:"cycle_type"`                                  // 频率
	CycleCount int32     `gorm:"column:cycle_count;comment:次数" json:"cycle_count"`                                // 次数
	Status     string    `gorm:"column:status;comment:状态" json:"status"`                                          // 状态
	Remark     string    `gorm:"column:remark;comment:备注" json:"remark"`                                          // 备注
	Attr1      string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                         // 预留字段1
	Attr2      string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                         // 预留字段2
	Attr3      int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                         // 预留字段3
	Attr4      int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                         // 预留字段4
}

type SysDeviceCheckPlanReq struct {
	Name string `form:"name"`
	baize.BaseEntityDQL
}

type SysDeviceCheckPlanList struct {
	Rows  []*SysDeviceCheckPlan `json:"rows"`
	Total int64                 `json:"total"`
}
