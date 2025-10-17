package deviceMonitorModel

import (
	"nova-factory-server/app/baize"
	"time"
)

// SysIotDbDevMap 数据和iotdb 对照表
type SysIotDbDevMap struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"` // 自增标识
	DeviceID   int64  `gorm:"column:device_id;not null;comment:设备id" json:"device_id"`        // 设备id
	TemplateID int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id"`    // 模板id
	DataID     int64  `gorm:"column:data_id;not null;comment:测点id" json:"data_id"`            // 测点id
	Device     string `gorm:"column:device;not null;comment:设备名字" json:"device"`              // 设备名字
	DataName   string `gorm:"column:data_name;not null;comment:数据名字" json:"data_name"`        // 数据名字
	Unit       string `gorm:"column:unit;not null;comment:数据单位" json:"unit"`
	DeptID     int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`      // 部门ID
	State      bool   `gorm:"column:state;comment:操作状态（0正常-1删除）" json:"state"` // 操作状态（0正常-1删除）
	baize.BaseEntity
}

// SysIotDbDevMapData 返回数据
type SysIotDbDevMapData struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id,string"` // 自增标识
	DeviceID   int64  `gorm:"column:device_id;not null;comment:设备id" json:"device_id"`               // 设备id
	TemplateID int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id"`           // 模板id
	DataID     int64  `gorm:"column:data_id;not null;comment:测点id" json:"data_id"`                   // 测点id
	Device     string `gorm:"column:device;not null;comment:设备名字" json:"device"`                     // 设备名字
	DataName   string `gorm:"column:data_name;not null;comment:数据名字" json:"data_name"`               // 数据名字
	Unit       string `gorm:"column:unit;not null;comment:单位" json:"unit"`                           // 数据名字
	DevName    string `gorm:"-" json:"dev_name"`                                                     //设备名字
}

// DevData 设备上报的数据
type DevData struct {
	Time       time.Time `json:"time"`
	Value      float64   `json:"value"`
	Name       string    `json:"name"`
	Unit       string    `json:"unit"`
	DeviceID   int64     `json:"device_id,string"`   // 设备id
	TemplateID int64     `json:"template_id,string"` // 模板id
	DataID     int64     `json:"data_id,string"`     // 测点id
	Dev        string    `json:"dev"`
	DevName    string    `json:"dev_name"` //设备名字
}

type DevDataReq struct {
	Dev       []string `json:"dev" form:"dev"`
	Start     uint64   `json:"start" form:"start"`
	End       uint64   `json:"end" form:"end"`
	DataScope string   `swaggerignore:"true"`
	OrderBy   string   `json:"orderBy" form:"orderBy"`                   //排序字段
	IsAsc     string   `json:"isAsc" form:"isAsc"`                       //排序规则  降序desc   asc升序
	Page      int64    `json:"pageNum" default:"1" form:"pageNum"`       //第几页
	Size      int64    `json:"pageSize" default:"10000" form:"pageSize"` //数量
}
type DevDataResp struct {
	Rows  []DevData `json:"rows"`
	Total uint64    `json:"total"`
}

type DevListReq struct {
	DataName string `gorm:"column:data_name;not null;comment:数据名字" form:"data_name"` // 数据名字
	baize.BaseEntityDQL
}

type DevListResp struct {
	Rows  []*SysIotDbDevMap `json:"rows"`
	Total uint64            `json:"total"`
}
