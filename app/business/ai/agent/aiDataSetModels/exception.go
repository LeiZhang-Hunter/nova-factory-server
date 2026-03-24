package aiDataSetModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
)

// SysAiPredictionException 预测列表
type SysAiPredictionException struct {
	ID            int64    `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                 // id
	ReasonID      int64    `gorm:"column:reason_id;not null;comment:模型推理id" json:"reason_id,string"` // 模型推理id
	ActionID      int64    `gorm:"column:action_id;not null;comment:处理通知id" json:"action_id,string"` // 处理通知id
	Threshold     int64    `gorm:"column:threshold;not null;comment:threshold" json:"threshold"`     // threshold
	Name          string   `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                  // 智能预警名称
	Dev           string   `gorm:"column:dev;comment:测点名称" json:"-"`                                 // 测点名称
	DevList       []string `gorm:"-" json:"dev"`                                                     // 智能预警名称
	Model         string   `gorm:"column:model;comment:预测模型" json:"model"`                           // 预测模型
	Interval      int64    `gorm:"column:interval;not null;comment:预测时间段" json:"interval,string"`    // 预测时间段
	PredictLength int64    `gorm:"column:predict_length;comment:预测长度" json:"predict_length,string"`  // 预测长度
	AggFunction   string   `gorm:"column:agg_function;comment:聚合函数，用来计算图表" json:"agg_function"`      // 聚合函数，用来计算图表
	DeptID        int64    `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                       // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SetSysAiPredictionException struct {
	ID            int64    `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                 // id
	ReasonID      int64    `gorm:"column:reason_id;not null;comment:模型推理id" json:"reason_id,string"` // 模型推理id
	ActionID      int64    `gorm:"column:action_id;not null;comment:处理通知id" json:"action_id,string"` // 处理通知id
	Threshold     int64    `gorm:"column:threshold;not null;comment:threshold" json:"threshold"`     // threshold
	Name          string   `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                  // 智能预警名称
	Dev           []string `gorm:"column:dev;comment:测点名称" json:"dev"`                               // 测点名称
	Model         string   `gorm:"column:model;comment:预测模型" json:"model"`                           // 预测模型
	Interval      int64    `gorm:"column:interval;not null;comment:预测时间段" json:"interval,string"`    // 预测时间段
	PredictLength int64    `gorm:"column:predict_length;comment:预测长度" json:"predict_length"`         // 预测长度
	AggFunction   string   `gorm:"column:agg_function;comment:聚合函数，用来计算图表" json:"agg_function"`      // 聚合函数，用来计算图表
}

func ToSysAiPredictionException(set *SetSysAiPredictionException) *SysAiPredictionException {
	dev, err := json.Marshal(set.Dev)
	if err != nil {
		zap.L().Error("json encode advanced failed", zap.Error(err))
		dev = []byte("")
	}
	return &SysAiPredictionException{
		ID:            set.ID,
		ReasonID:      set.ReasonID,
		ActionID:      set.ActionID,
		Threshold:     set.Threshold,
		Name:          set.Name,
		Dev:           string(dev),
		Model:         set.Model,
		Interval:      set.Interval,
		PredictLength: set.PredictLength,
		AggFunction:   set.AggFunction,
	}
}

type SysAiPredictionExceptionListReq struct {
	Name string `form:"name"` // 告警策略名称
	baize.BaseEntityDQL
}

type SysAiPredictionExceptionList struct {
	Rows  []*SysAiPredictionException `json:"rows"`
	Total uint64                      `json:"total"`
}
