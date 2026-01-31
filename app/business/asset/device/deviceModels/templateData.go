package deviceModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
)

// SysModbusDeviceConfigData modbus数据配置
type SysModbusDeviceConfigData struct {
	DeviceConfigID        int64  `gorm:"column:device_config_id;primaryKey;comment:文档id" json:"device_config_id,string"` // 文档id
	TemplateID            int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id,string"`             // 模板id
	Name                  string `gorm:"column:name;comment:数据名称" json:"name"`                                           // 数据名称
	Type                  string `gorm:"column:type;comment:数据类型" json:"type"`                                           // 数据类型
	DataType              string `gorm:"column:data_type;comment:设备节点类型，到底是开关还是数值" json:"data_type"`                     // 设备节点类型，到底是开关还是数值
	Slave                 string `gorm:"column:slave;comment:从设备地址" json:"slave"`                                        // 从设备地址
	Register              *int64 `gorm:"column:register;not null;comment:寄存器/偏移量" json:"register,string"`                // 寄存器/偏移量
	StorageType           string `gorm:"column:storage_type;comment:存储策略" json:"storage_type"`                           // 存储策略
	PredictEnable         *bool  `gorm:"column:predict_enable;comment:是否开启预测" json:"predict_enable"`                     // 是否开启预测
	PerturbationVariables string `gorm:"column:perturbation_variables;comment:扰动变量" json:"perturbation_variables"`       // 扰动变量
	Annotation            string `gorm:"column:annotation;comment:注解" json:"annotation"`                                 // 扰动变量
	GraphEnable           *bool  `gorm:"column:graph_enable;comment:是否开启图表" json:"graph_enable"`                         // 是否开启图表
	AggFunction           string `gorm:"column:agg_function;comment:聚合函数" json:"agg_function"`                           // 聚合函数
	Unit                  string `gorm:"column:unit;comment:单位" json:"unit"`                                             // 单位
	Precision             int64  `gorm:"column:precision;not null;comment:数据精度" json:"precision,string"`                 // 数据精度
	FunctionCode          int    `gorm:"column:function_code;comment:功能码" json:"function_code,string"`                   // 功能码
	Expression            string `gorm:"column:expression;comment:表达式" json:"expression"`                                // 表达式
	Mode                  *int   `gorm:"column:mode;comment:功能码" json:"mode,string"`                                     // 功能码
	DataFormat            string `gorm:"column:data_format;comment:读写方式" json:"data_format"`                             // 读写方式
	Sort                  string `gorm:"column:sort;comment:数据排序" json:"sort"`                                           // 数据排序
	DeptID                int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                     // 部门ID
	State                 bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                               // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func OfSysModbusDeviceConfigData(req *SetSysModbusDeviceConfigDataReq) *SysModbusDeviceConfigData {
	var perturbationVariables []byte
	var err error
	if len(req.PerturbationVariables) != 0 {
		perturbationVariables, err = json.Marshal(req.PerturbationVariables)
		if err != nil {
			zap.L().Error("perturbationVariables json marshal error", zap.Error(err))
			return nil
		}
	}

	var annotations []byte
	if len(req.Annotations) != 0 {
		annotations, err = json.Marshal(req.Annotations)
		if err != nil {
			zap.L().Error("annotations json marshal error", zap.Error(err))
			return nil
		}
	}

	return &SysModbusDeviceConfigData{
		DeviceConfigID:        req.DeviceConfigID,
		TemplateID:            req.TemplateID,
		Name:                  req.Name,
		Type:                  req.Type,
		Slave:                 req.Slave,
		Register:              req.Register,
		StorageType:           req.StorageType,
		Unit:                  req.Unit,
		Precision:             req.Precision,
		FunctionCode:          req.FunctionCode,
		Mode:                  req.Mode,
		DataFormat:            req.DataFormat,
		AggFunction:           req.AggFunction,
		Sort:                  req.Sort,
		DataType:              req.DataType,
		PredictEnable:         req.PredictEnable,
		Annotation:            string(annotations),
		PerturbationVariables: string(perturbationVariables),
		GraphEnable:           req.GraphEnable,
		Expression:            req.Expression,
	}
}

// PerturbationVariableData 扰动变量影响预测值的数据
type PerturbationVariableData struct {
	DeviceId   int64 `json:"device_id,string"`
	TemplateId int64 `json:"template_id,string"`
	DataId     int64 `json:"data_id,string"`
}

type Annotation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetSysModbusDeviceConfigDataReq struct {
	DeviceConfigID        int64                      `gorm:"column:device_config_id;primaryKey;comment:文档id" json:"device_config_id,string"`        // 文档id
	TemplateID            int64                      `gorm:"column:template_id;not null;comment:模板id" binding:"required" json:"template_id,string"` // 模板id
	Name                  string                     `gorm:"column:name;comment:数据名称" binding:"required" json:"name"`                               // 数据名称
	Type                  string                     `gorm:"column:type;comment:数据类型" binding:"required" json:"type"`                               // 数据类型
	DataType              string                     `gorm:"column:data_type;comment:设备节点类型，到底是开关还是数值" json:"data_type"`                            // 设备节点类型，到底是开关还是数值
	Slave                 string                     `gorm:"column:slave;comment:从设备地址"  json:"slave"`                                              // 从设备地址
	Register              *int64                     `gorm:"column:register;not null;comment:寄存器/偏移量" json:"register,string"`                       // 寄存器/偏移量
	AggFunction           string                     `gorm:"column:agg_function;comment:聚合函数" json:"agg_function"`                                  // 聚合函数
	StorageType           string                     `gorm:"column:storage_type;comment:存储策略" json:"storage_type"`                                  // 存储策略
	Unit                  string                     `gorm:"column:unit;comment:单位" json:"unit"`                                                    // 单位
	Precision             int64                      `gorm:"column:precision;not null;comment:数据精度" json:"precision,string"`                        // 数据精度
	PredictEnable         *bool                      `gorm:"column:predict_enable;comment:是否开启预测" json:"predict_enable"`                            // 是否开启预测
	PerturbationVariables []PerturbationVariableData `gorm:"column:perturbation_variables;comment:扰动变量" json:"perturbation_variables"`              // 扰动变量
	Annotations           []Annotation               `gorm:"column:annotations;comment:注解" json:"annotation"`                                       // 扰动变量
	GraphEnable           *bool                      `gorm:"column:graph_enable;comment:是否开启图表" json:"graph_enable"`                                // 是否开启图表
	FunctionCode          int                        `gorm:"column:function_code;comment:功能码" json:"function_code,string"`                          // 功能码
	Mode                  *int                       `gorm:"column:mode;comment:功能码" json:"mode,string"`                                            // 功能码
	DataFormat            string                     `gorm:"column:data_format;comment:读写方式" json:"data_format"`                                    // 读写方式
	Expression            string                     `gorm:"column:expression;comment:表达式" json:"expression"`                                       // 读写方式
	Sort                  string                     `gorm:"column:sort;comment:数据排序" binding:"required" json:"sort"`                               // 数据排序
}

func ToSetSysModbusDeviceConfigDataReq(req *SysModbusDeviceConfigData) *SetSysModbusDeviceConfigDataReq {
	var perturbationVariables []PerturbationVariableData = make([]PerturbationVariableData, 0)
	var err error
	if len(req.PerturbationVariables) != 0 {
		err = json.Unmarshal([]byte(req.PerturbationVariables), &perturbationVariables)
		if err != nil {
			zap.L().Error("ToSetSysModbusDeviceConfigDataReq json Unmarshal err", zap.Error(err))
		}
	}

	var annotations []Annotation = make([]Annotation, 0)
	if len(req.Annotation) != 0 {
		err = json.Unmarshal([]byte(req.Annotation), &annotations)
		if err != nil {
			zap.L().Error("ToSetSysModbusDeviceConfigDataReq json Unmarshal err", zap.Error(err))
		}
	}

	return &SetSysModbusDeviceConfigDataReq{
		DeviceConfigID:        req.DeviceConfigID,
		TemplateID:            req.TemplateID,
		Name:                  req.Name,
		Type:                  req.Type,
		Slave:                 req.Slave,
		Register:              req.Register,
		StorageType:           req.StorageType,
		Unit:                  req.Unit,
		Precision:             req.Precision,
		FunctionCode:          req.FunctionCode,
		Mode:                  req.Mode,
		DataFormat:            req.DataFormat,
		AggFunction:           req.AggFunction,
		Expression:            req.Expression,
		Sort:                  req.Sort,
		DataType:              req.DataType,
		PredictEnable:         req.PredictEnable,
		PerturbationVariables: perturbationVariables,
		Annotations:           annotations,
		GraphEnable:           req.GraphEnable,
	}
}

type SysModbusDeviceConfigDataListReq struct {
	TemplateID int64 `gorm:"column:template_id;not null;comment:模板id" form:"template_id" binding:"required" json:"template_id,string"` // 模板id
	baize.BaseEntityDQL
}

type SysModbusDeviceConfigDataListData struct {
	Rows  []*SetSysModbusDeviceConfigDataReq `json:"rows"`
	Total int64                              `json:"total"`
}

type DeviceReqInfo struct {
	DeviceID int64 `gorm:"column:device_id;not null;comment:设备id" form:"device_id,string"` // 设备id
}
