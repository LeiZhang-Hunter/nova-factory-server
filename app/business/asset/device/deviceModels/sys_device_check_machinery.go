package deviceModels

import "time"

// SysDeviceCheckMachinery 点检设备表
type SysDeviceCheckMachinery struct {
	RecordID       int64     `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:流水号" json:"record_id"` // 流水号
	PlanID         int64     `gorm:"column:plan_id;not null;comment:计划ID" json:"plan_id"`                         // 计划ID
	MachineryID    int64     `gorm:"column:machinery_id;not null;comment:设备ID" json:"machinery_id"`               // 设备ID
	MachineryCode  string    `gorm:"column:machinery_code;not null;comment:设备编码" json:"machinery_code"`           // 设备编码
	MachineryName  string    `gorm:"column:machinery_name;not null;comment:设备名称" json:"machinery_name"`           // 设备名称
	MachineryBrand string    `gorm:"column:machinery_brand;comment:品牌" json:"machinery_brand"`                    // 品牌
	MachinerySpec  string    `gorm:"column:machinery_spec;comment:规格型号" json:"machinery_spec"`                    // 规格型号
	Remark         string    `gorm:"column:remark;comment:备注" json:"remark"`                                      // 备注
	Attr1          string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                     // 预留字段1
	Attr2          string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                     // 预留字段2
	Attr3          int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                     // 预留字段3
	Attr4          int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                     // 预留字段4
	DeptID         int64     `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                  // 部门ID
	CreateBy       int64     `gorm:"column:create_by;comment:创建者" json:"create_by"`                               // 创建者
	CreateTime     time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`                          // 创建时间
	UpdateBy       int64     `gorm:"column:update_by;comment:更新者" json:"update_by"`                               // 更新者
	UpdateTime     time.Time `gorm:"column:update_time;comment:更新时间" json:"update_time"`                          // 更新时间
	State          bool      `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                            // 操作状态（0正常 -1删除）
}
