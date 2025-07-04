package deviceModels

import "nova-factory-server/app/baize"

// SysModbusDeviceConfigData modbus数据配置
type SysModbusDeviceConfigData struct {
	DeviceConfigID int64  `gorm:"column:device_config_id;primaryKey;comment:文档id" json:"device_config_id,string"` // 文档id
	TemplateID     int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id,string"`             // 模板id
	Name           string `gorm:"column:name;comment:数据名称" json:"name"`                                           // 数据名称
	Type           string `gorm:"column:type;comment:数据类型" json:"type"`                                           // 数据类型
	DataType       string `gorm:"column:data_type;comment:设备节点类型，到底是开关还是数值" json:"data_type"`                     // 设备节点类型，到底是开关还是数值
	Slave          string `gorm:"column:slave;comment:从设备地址" json:"slave"`                                        // 从设备地址
	Register       int64  `gorm:"column:register;not null;comment:寄存器/偏移量" json:"register,string"`                // 寄存器/偏移量
	StorageType    string `gorm:"column:storage_type;comment:存储策略" json:"storage_type"`                           // 存储策略
	Unit           string `gorm:"column:unit;comment:单位" json:"unit"`                                             // 单位
	Precision      int64  `gorm:"column:precision;not null;comment:数据精度" json:"precision,string"`                 // 数据精度
	FunctionCode   int    `gorm:"column:function_code;comment:功能码" json:"function_code,string"`                   // 功能码
	Mode           int    `gorm:"column:mode;comment:功能码" json:"mode,string"`                                     // 功能码
	DataFormat     string `gorm:"column:data_format;comment:读写方式" json:"data_format"`                             // 读写方式
	Sort           string `gorm:"column:sort;comment:数据排序" json:"sort"`                                           // 数据排序
	DeptID         int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                     // 部门ID
	State          bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                               // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func OfSysModbusDeviceConfigData(req *SetSysModbusDeviceConfigDataReq) *SysModbusDeviceConfigData {
	return &SysModbusDeviceConfigData{
		DeviceConfigID: req.DeviceConfigID,
		TemplateID:     req.TemplateID,
		Name:           req.Name,
		Type:           req.Type,
		Slave:          req.Slave,
		Register:       req.Register,
		StorageType:    req.StorageType,
		Unit:           req.Unit,
		Precision:      req.Precision,
		FunctionCode:   req.FunctionCode,
		Mode:           req.Mode,
		DataFormat:     req.DataFormat,
		Sort:           req.Sort,
	}
}

type SetSysModbusDeviceConfigDataReq struct {
	DeviceConfigID int64  `gorm:"column:device_config_id;primaryKey;comment:文档id" json:"device_config_id,string"`        // 文档id
	TemplateID     int64  `gorm:"column:template_id;not null;comment:模板id" binding:"required" json:"template_id,string"` // 模板id
	Name           string `gorm:"column:name;comment:数据名称" binding:"required" json:"name"`                               // 数据名称
	Type           string `gorm:"column:type;comment:数据类型" binding:"required" json:"type"`                               // 数据类型
	DataType       string `gorm:"column:data_type;comment:设备节点类型，到底是开关还是数值" json:"data_type"`                            // 设备节点类型，到底是开关还是数值
	Slave          string `gorm:"column:slave;comment:从设备地址" binding:"required" json:"slave"`                            // 从设备地址
	Register       int64  `gorm:"column:register;not null;comment:寄存器/偏移量" json:"register,string"`                       // 寄存器/偏移量
	StorageType    string `gorm:"column:storage_type;comment:存储策略" json:"storage_type"`                                  // 存储策略
	Unit           string `gorm:"column:unit;comment:单位" json:"unit"`                                                    // 单位
	Precision      int64  `gorm:"column:precision;not null;comment:数据精度" json:"precision,string"`                        // 数据精度
	FunctionCode   int    `gorm:"column:function_code;comment:功能码" json:"function_code,string"`                          // 功能码
	Mode           int    `gorm:"column:mode;comment:功能码" json:"mode,string"`                                            // 功能码
	DataFormat     string `gorm:"column:data_format;comment:读写方式" json:"data_format"`                                    // 读写方式
	Sort           string `gorm:"column:sort;comment:数据排序" binding:"required" json:"sort"`                               // 数据排序
}

type SysModbusDeviceConfigDataListReq struct {
	TemplateID int64 `gorm:"column:template_id;not null;comment:模板id" form:"template_id" binding:"required" json:"template_id,string"` // 模板id
	baize.BaseEntityDQL
}

type SysModbusDeviceConfigDataListData struct {
	Rows  []*SysModbusDeviceConfigData `json:"rows"`
	Total int64                        `json:"total"`
}
