package aiDataSetModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
)

// SysAiPrediction 预测列表
type SysAiPrediction struct {
	ID                    int64              `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                 // id
	ReasonID              int64              `gorm:"column:reason_id;not null;comment:模型推理id" json:"reason_id,string"` // 模型推理id
	ActionID              int64              `gorm:"column:action_id;not null;comment:处理通知id" json:"action_id,string"` // 处理通知id
	Name                  string             `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                  // 智能预警名称
	Dev                   string             `gorm:"column:dev;not null;comment:智能预警名称" json:"-"`                      // 智能预警名称
	DevList               []string           `gorm:"-" json:"dev"`                                                     // 智能预警名称
	Advanced              string             `gorm:"column:advanced;comment:告警规则" json:"-"`                            // 告警规则
	AdvancedData          *logalert.Advanced `gorm:"-" json:"advanced"`
	Model                 string             `gorm:"column:model;comment:预测模型" json:"model"`                                   // 预测模型
	Field                 string             `gorm:"column:field;comment:预测字段" json:"field"`                                   // 预测字段
	Interval              int64              `gorm:"column:interval;not null;comment:预测时间段" json:"interval,string"`            // 预测时间段
	PredictLength         int64              `gorm:"column:predict_length;comment:预测长度" json:"predict_length"`                 // 预测长度
	Threshold             int64              `gorm:"column:threshold;comment:预测阀值" json:"threshold"`                           // 预测长度
	AggFunction           string             `gorm:"column:agg_function;comment:聚合函数，用来计算图表" json:"agg_function"`              // 聚合函数，用来计算图表
	PerturbationVariables string             `gorm:"column:perturbation_variables;comment:关联变量" json:"perturbation_variables"` // 关联变量
	DeptID                int64              `gorm:"column:dept_id;comment:部门ID" json:"dept_id,string"`                        // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SetSysAiPrediction struct {
	ID            int64             `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                 // id
	Name          string            `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                  // 智能预警名称
	Dev           []string          `gorm:"column:name;not null;comment:智能预警名称" json:"dev"`                   // 智能预警名称
	Model         string            `gorm:"column:model;comment:预测模型" json:"model"`                           // 预测模型
	AggFunction   string            `gorm:"column:agg_function;comment:聚合函数，用来计算图表" json:"agg_function"`      // 聚合函数，用来计算图表
	ReasonID      int64             `gorm:"column:reason_id;not null;comment:模型推理id" json:"reason_id,string"` // 模型推理id
	ActionID      int64             `gorm:"column:action_id;not null;comment:处理通知id" json:"action_id,string"` // 处理通知id
	PredictLength int64             `gorm:"column:predict_length;comment:预测长度" json:"predict_length"`         // 预测长度
	Threshold     int               `json:"threshold"`
	Advanced      logalert.Advanced `gorm:"column:advanced;comment:告警规则" json:"advanced"`                  // 告警规则
	Interval      int64             `gorm:"column:interval;not null;comment:预测时间段" json:"interval,string"` // 预测时间段
}

func ToSysAiPredictionList(data *SetSysAiPrediction) *SysAiPrediction {
	advanced, err := json.Marshal(data.Advanced)
	if err != nil {
		zap.L().Error("json encode advanced failed", zap.Error(err))
	}

	dev, err := json.Marshal(data.Dev)
	if err != nil {
		zap.L().Error("json encode advanced failed", zap.Error(err))
	}
	return &SysAiPrediction{
		ID:       data.ID,
		ReasonID: data.ReasonID,
		ActionID: data.ActionID,
		Name:     data.Name,
		Advanced: string(advanced),
		Model:    data.Model,
		Dev:      string(dev),
		//Field:         data.Field,
		Interval:      data.Interval,
		PredictLength: data.PredictLength,
		AggFunction:   data.AggFunction,
	}
}

type SysAiPredictionListReq struct {
	Name string `form:"name"` // 告警策略名称
	baize.BaseEntityDQL
}

type SysAiPredictionList struct {
	Rows  []*SysAiPrediction `json:"rows"`
	Total uint64             `json:"total"`
}
