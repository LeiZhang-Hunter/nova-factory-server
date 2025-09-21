package aiDataSetModels

import (
	"nova-factory-server/app/baize"
)

// SysAiPredictionControl 趋势控制
type SysAiPredictionControl struct {
	ID              int64  `gorm:"column:id;primaryKey;comment:id" json:"id"`                               // id
	DeviceGatewayID int64  `gorm:"column:device_gateway_id;not null;comment:网关id" json:"device_gateway_id"` // 网关id
	Name            string `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                         // 智能预警名称
	Parallelism     int64  `gorm:"column:parallelism;not null;comment:并发" json:"parallelism"`               // 并发
	Threshold       int64  `gorm:"column:threshold;not null;comment:threshold" json:"threshold"`            // threshold
	Model           string `gorm:"column:model;comment:预测模型" json:"model"`                                  // 预测模型
	Interval        int64  `gorm:"column:interval;not null;comment:预测时间段" json:"interval"`                  // 预测时间段
	PredictLength   int64  `gorm:"column:predict_length;comment:预测长度" json:"predict_length"`                // 预测长度
	AggFunction     string `gorm:"column:agg_function;comment:聚合函数，用来计算图表" json:"agg_function"`             // 聚合函数，用来计算图表
	DeptID          int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                              // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SetSysAiPredictionControl struct {
	ID              int64  `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                               // id
	DeviceGatewayID int64  `gorm:"column:device_gateway_id;not null;comment:网关id" json:"device_gateway_id,string"` // 网关id
	Name            string `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                                // 智能预警名称
	Parallelism     int64  `gorm:"column:parallelism;not null;comment:并发" json:"parallelism"`                      // 并发
	Threshold       int64  `gorm:"column:threshold;not null;comment:threshold" json:"threshold"`                   // threshold
	Model           string `gorm:"column:model;comment:预测模型" json:"model"`                                         // 预测模型
	Interval        int64  `gorm:"column:interval;not null;comment:预测时间段" json:"interval,string"`                  // 预测时间段
	PredictLength   int64  `gorm:"column:predict_length;comment:预测长度" json:"predict_length"`                       // 预测长度
	AggFunction     string `gorm:"column:agg_function;comment:聚合函数，用来计算图表" json:"agg_function"`                    // 聚合函数，用来计算图表
}

func ToSysAiPredictionControl(set *SetSysAiPredictionControl) *SysAiPredictionControl {
	return &SysAiPredictionControl{
		ID:              set.ID,
		DeviceGatewayID: set.DeviceGatewayID,
		Name:            set.Name,
		Parallelism:     set.Parallelism,
		Threshold:       set.Threshold,
		Model:           set.Model,
		Interval:        set.Interval,
		PredictLength:   set.PredictLength,
		AggFunction:     set.AggFunction,
	}
}

type SysAiPredictionControlListReq struct {
	Name string `form:"name"` // 告警策略名称
	baize.BaseEntityDQL
}

type SysAiPredictionControlList struct {
	Rows  []*SysAiPredictionControl `json:"rows"`
	Total uint64                    `json:"total"`
}
