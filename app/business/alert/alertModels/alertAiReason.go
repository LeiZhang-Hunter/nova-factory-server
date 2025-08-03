package alertModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
)

type SysAlertAiReason struct {
	ID            int64    `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id,string"` // 主键
	Name          string   `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                     // 告警策略名称
	Prompt        string   `gorm:"column:prompt;comment:提示词设置" json:"prompt"`                           // 提示词设置
	Message       string   `gorm:"column:message;comment:提问消息模板" json:"message"`                        // 提问消息模板
	DatasetIds    string   `gorm:"column:dataset_ids;comment:知识库id列表" json:"-"`                         // 知识库id列表
	DatasetIdList []string `gorm:"-" json:"dataset_ids"`                                                // 知识库id列表
	baize.BaseEntity
	DeptID int64 `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`       // 部门ID
	State  bool  `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func FromSetAlertReasonToData(reason *SetAlertAiReason) *SysAlertAiReason {
	var err error
	var datasetIds []byte
	if len(reason.DatasetIds) != 0 {
		datasetIds, err = json.Marshal(reason.DatasetIds)
		if err != nil {
			zap.L().Error("json marshal", zap.Error(err))
		}
	}

	return &SysAlertAiReason{
		ID:         reason.ID,
		Name:       reason.Name,
		Prompt:     reason.Prompt,
		Message:    reason.Message,
		DatasetIds: string(datasetIds),
	}
}

type SetAlertAiReason struct {
	ID         int64    `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id,string"` // 主键
	Name       string   `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                     // 告警策略名称
	Prompt     string   `gorm:"column:prompt;comment:提示词设置" json:"prompt"`                           // 提示词设置
	Message    string   `gorm:"column:message;comment:提问消息模板" json:"message"`                        // 提问消息模板
	DatasetIds []string `gorm:"column:dataset_ids;comment:知识库id列表" json:"dataset_ids"`               // 知识库id列表
}

type SysAlertAiReasonReq struct {
	Name    string `form:"name"`
	Prompt  string `form:"prompt"`
	Message string `form:"message"`
	baize.BaseEntityDQL
}

type SysAlertReasonList struct {
	Rows  []*SysAlertAiReason `json:"rows"`
	Total uint64              `json:"total"`
}
