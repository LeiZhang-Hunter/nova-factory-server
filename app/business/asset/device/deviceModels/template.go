package deviceModels

import (
	"nova-factory-server/app/baize"
)

// SysDeviceTemplate modbus数据模板
type SysDeviceTemplate struct {
	TemplateID   int64  `gorm:"column:template_id;primaryKey;comment:设备主键" json:"template_id,string"` // 设备主键
	Name         string `gorm:"column:name;comment:设备名称" json:"name"`                                 // 设备名称
	TemplateType int    `gorm:"column:template_type;comment:模板类型0是私有1 是共有" json:"template_type"`      // 模板类型0是私有1 是共有
	Vendor       string `gorm:"column:vendor;comment:供应商" json:"vendor"`                              // 供应商
	Protocol     string `gorm:"column:protocol;comment:协议类型" json:"protocol"`                         // 协议类型
	Remark       string `gorm:"column:remark;comment:备注" json:"remark"`                               // 备注
	DeptID       int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                           // 部门ID
	State        bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                     // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func ToSysDeviceTemplate(set *SysDeviceTemplateSetReq) *SysDeviceTemplate {
	return &SysDeviceTemplate{
		TemplateID: set.TemplateID,
		Name:       set.Name,
		Protocol:   set.Protocol,
		Vendor:     set.Vendor,
		Remark:     set.Remark,
	}
}

type SysDeviceTemplateSetReq struct {
	TemplateID int64  `gorm:"column:template_id;primaryKey;comment:设备主键"  json:"template_id,string"` // 设备主键
	Name       string `gorm:"column:name;comment:设备名称" binding:"required" json:"name"`               // 设备名称
	Vendor     string `gorm:"column:vendor;comment:供应商" binding:"required" json:"vendor"`            // 供应商
	Protocol   string `gorm:"column:protocol;comment:协议类型" binding:"required" json:"protocol"`       // 协议类型
	Remark     string `gorm:"column:remark;comment:备注" json:"remark"`                                // 备注
}

type SysDeviceTemplateDQL struct {
	Name     string `form:"name"`
	Protocol string `form:"protocol"` // 协议类型
	baize.BaseEntityDQL
}

type SysDeviceTemplateListData struct {
	Rows  []*SysDeviceTemplate `json:"rows"`
	Total int64                `json:"total"`
}
