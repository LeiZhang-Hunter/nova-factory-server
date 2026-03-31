package deviceModels

import "nova-factory-server/app/baize"

// SysDeviceCheckMachinery 点检设备表
type SysDeviceCheckMachinery struct {
	RecordID       int64  `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:流水号" json:"record_id,string"` // 流水号
	PlanID         int64  `gorm:"column:plan_id;not null;comment:计划ID" json:"plan_id,string"`                         // 计划ID
	MachineryID    int64  `gorm:"column:machinery_id;not null;comment:设备ID" json:"machinery_id"`                      // 设备ID
	MachineryCode  string `gorm:"column:machinery_code;not null;comment:设备编码" json:"machinery_code"`                  // 设备编码
	MachineryName  string `gorm:"column:machinery_name;not null;comment:设备名称" json:"machinery_name"`                  // 设备名称
	MachineryBrand string `gorm:"column:machinery_brand;comment:品牌" json:"machinery_brand"`                           // 品牌
	MachinerySpec  string `gorm:"column:machinery_spec;comment:规格型号" json:"machinery_spec"`                           // 规格型号
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                                             // 备注
	Attr1          string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                            // 预留字段1
	Attr2          string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                            // 预留字段2
	Attr3          int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                            // 预留字段3
	Attr4          int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                            // 预留字段4
	DeptID         int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                         // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func ToSysDeviceCheckMachinery(vo *SysDeviceCheckMachineryVO) *SysDeviceCheckMachinery {
	return &SysDeviceCheckMachinery{
		RecordID:       vo.RecordID,
		PlanID:         vo.PlanID,
		MachineryID:    vo.MachineryID,
		MachineryCode:  vo.MachineryCode,
		MachineryName:  vo.MachineryName,
		MachineryBrand: vo.MachineryBrand,
		MachinerySpec:  vo.MachinerySpec,
		Remark:         vo.Remark,
		Attr1:          vo.Attr1,
		Attr2:          vo.Attr2,
		Attr3:          vo.Attr3,
		Attr4:          vo.Attr4,
	}
}

type SysDeviceCheckMachineryVO struct {
	RecordID       int64  `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:流水号" json:"record_id,string"` // 流水号
	PlanID         int64  `gorm:"column:plan_id;not null;comment:计划ID" json:"plan_id,string"`                         // 计划ID
	MachineryID    int64  `gorm:"column:machinery_id;not null;comment:设备ID" json:"machinery_id"`                      // 设备ID
	MachineryCode  string `gorm:"column:machinery_code;not null;comment:设备编码" json:"machinery_code"`                  // 设备编码
	MachineryName  string `gorm:"column:machinery_name;not null;comment:设备名称" json:"machinery_name"`                  // 设备名称
	MachineryBrand string `gorm:"column:machinery_brand;comment:品牌" json:"machinery_brand"`                           // 品牌
	MachinerySpec  string `gorm:"column:machinery_spec;comment:规格型号" json:"machinery_spec"`                           // 规格型号
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                                             // 备注
	Attr1          string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                            // 预留字段1
	Attr2          string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                            // 预留字段2
	Attr3          int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                            // 预留字段3
	Attr4          int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                            // 预留字段4
}

type SysDeviceCheckMachineryReq struct {
	PlanID int64 `gorm:"column:plan_id;not null;comment:计划ID" json:"plan_id,string"` // 计划ID
	baize.BaseEntityDQL
}

type SysDeviceCheckMachineryList struct {
	Rows  []*SysDeviceCheckMachinery `json:"rows"`
	Total int64                      `json:"total"`
}
