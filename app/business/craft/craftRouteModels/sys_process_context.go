package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

// SysProProcessContent 生产工序内容表
type SysProProcessContent struct {
	ContentID      uint64       `gorm:"column:content_id;primaryKey;autoIncrement:true;comment:内容ID" json:"content_id,string"` // 内容ID
	ProcessID      uint64       `gorm:"column:process_id;not null;comment:工序ID" json:"process_id,string" binding:"required"`   // 工序ID
	ControlName    string       `gorm:"column:control_name;comment:内容说明" json:"control_name" binding:"required"`               // 内容说明
	OrderNum       int32        `gorm:"column:order_num;comment:顺序编号" json:"order_num"`                                        // 顺序编号
	ContentText    string       `gorm:"column:content_text;comment:内容说明" json:"content_text" binding:"required"`               // 内容说明
	Device         string       `gorm:"column:device;comment:辅助设备" json:"device"`                                              // 辅助设备
	Material       string       `gorm:"column:material;comment:辅助材料" json:"material"`                                          // 辅助材料
	DocURL         string       `gorm:"column:doc_url;comment:材料URL" json:"doc_url"`                                           // 材料URL
	ControlType    string       `gorm:"column:control_type;comment:控制方式" json:"control_type"`                                  // 材料URL
	Remark         string       `gorm:"column:remark;comment:备注" json:"remark"`                                                // 备注
	Attr1          string       `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                               // 预留字段1
	Attr2          string       `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                               // 预留字段2
	Attr3          int32        `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                               // 预留字段3
	Attr4          int32        `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                               // 预留字段4
	CreateUserName string       `json:"createUserName" gorm:"-"`
	UpdateUserName string       `json:"updateUserName" gorm:"-"`
	Extension      string       `gorm:"column:extension;comment:触发规则" json:"-"`
	ControlRules   *ControlRule `gorm:"-" json:"control_rules,omitempty"`
	baize.BaseEntity
}

func NewSysProProcessContent(context *SysProSetProcessContent) *SysProProcessContent {
	return &SysProProcessContent{
		ContentID:   context.ContentID,
		ProcessID:   context.ProcessID,
		ControlName: context.ControlName,
		OrderNum:    context.OrderNum,
		ContentText: context.ContentText,
		ControlType: context.ControlType,
		Device:      context.Device,
		Material:    context.Material,
		DocURL:      context.DocURL,
		Remark:      context.Remark,
		Attr1:       context.Attr1,
		Attr2:       context.Attr2,
		Attr3:       context.Attr3,
		Attr4:       context.Attr4,
	}
}

type DeviceRuleInfo struct {
	DeviceId string `json:"deviceId"`
	DataId   string `json:"dataId"`
}

type ControllerAction struct {
	DeviceId    string `json:"device_id"`
	TemplateId  string `json:"template_id"`
	DataId      string `json:"data_id"`
	Value       string `json:"value"`
	ControlMode string `json:"control_mode"`
	Condition   string `json:"condition"`
	Interval    string `json:"interval"`
	DataFormat  string `json:"dataFormat"`
}

type TriggerCase struct {
	NextStep   string `json:"next_step"`
	Connector  string `json:"connector"`
	Conditions []struct {
		DataId     string `json:"data_id"`
		Operator   string `json:"operator"`
		DeviceId   string `json:"device_id"`
		TemplateId string `json:"template_id"`
		Value      string `json:"value"`
		Rule       string `json:"rule"`
		Connector  string `json:"connector"`
	} `json:"conditions"`
}

type TriggerRules struct {
	Name         string             `json:"name"`
	Actions      []ControllerAction `json:"actions"`
	CombinedRule string             `json:"combined_rule"`
	DataIds      []DeviceRuleInfo   `json:"dataIds"`
	Cases        []TriggerCase      `json:"cases"`
}

type PidRules struct {
	Proportional     float64            `json:"proportional"`
	Actions          []ControllerAction `json:"actions"`
	MaxControl       float64            `json:"max_control"`
	MinControl       float64            `json:"min_control"`
	Integral         float64            `json:"integral"`
	Derivative       float64            `json:"derivative"`
	ReferenceSignal  float64            `json:"reference_signal"`
	SamplingInterval uint64             `json:"sampling_interval"`
	DeviceId         string             `json:"device_id"`
	DataId           string             `json:"data_id"`
}

type CaptureData struct {
	DeviceId   string `json:"device_id"`
	DataId     string `json:"data_id"`
	TemplateId string `json:"template_id"`
}

type PredictData struct {
	DeviceId   string `json:"device_id"`
	DataId     string `json:"data_id"`
	TemplateId string `json:"template_id"`
}

type PredictRules struct {
	Actions       []ControllerAction `json:"actions"`
	Cases         []TriggerCase      `json:"cases"`
	Predicts      []*PredictData     `json:"predicts"`
	Threshold     int64              `json:"threshold"`      // threshold
	Model         string             `json:"model"`          // 预测模型
	Interval      int64              `json:"interval"`       // 预测时间段
	PredictLength int64              `json:"predict_length"` // 预测长度
	AggFunction   string             `json:"agg_function"`   // 聚合函数，用来计算图表
	IsContinue    bool               `json:"is_continue"`
	Step          uint64             `json:"step"`        // 聚合窗口的滑动步长（可选，默认与聚合窗口大小相同）
	WindowSize    uint64             `json:"window_size"` // 聚合窗口的大小（必须为正数）
}

// ControlRule 控制规则
type ControlRule struct {
	TriggerRules *TriggerRules  `json:"trigger_rules,omitempty"`
	PidRules     *PidRules      `json:"pid_rules"`
	PredictRules *PredictRules  `json:"predict_rules"`
	CaptureData  []*CaptureData `json:"captures"`
}

// SysProSetProcessContent 设置工序内容请求
type SysProSetProcessContent struct {
	ContentID    uint64       `gorm:"column:content_id;primaryKey;autoIncrement:true;comment:内容ID" json:"content_id,string"` // 内容ID
	ProcessID    uint64       `gorm:"column:process_id;not null;comment:工序ID" json:"process_id,string" binding:"required"`   // 工序ID
	ControlName  string       `gorm:"column:control_name;comment:内容说明" json:"control_name"`                                  // 内容说明
	OrderNum     int32        `gorm:"column:order_num;comment:顺序编号" json:"order_num"`                                        // 顺序编号
	ControlType  string       `json:"control_type"`
	ContentText  string       `gorm:"column:content_text;comment:内容说明" json:"content_text" binding:"required"` // 内容说明
	Device       string       `gorm:"column:device;comment:辅助设备" json:"device"`                                // 辅助设备
	Material     string       `gorm:"column:material;comment:辅助材料" json:"material"`                            // 辅助材料
	DocURL       string       `gorm:"column:doc_url;comment:材料URL" json:"doc_url"`                             // 材料URL
	Remark       string       `gorm:"column:remark;comment:备注" json:"remark"`                                  // 备注
	Attr1        string       `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                 // 预留字段1
	Attr2        string       `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                 // 预留字段2
	Attr3        int32        `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                 // 预留字段3
	Attr4        int32        `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                 // 预留字段4
	ControlRules *ControlRule `gorm:"-" json:"control_rules,omitempty"`
}

type SysProProcessContextListReq struct {
	ProcessID int64 `gorm:"column:process_id;not null;comment:工序ID" form:"process_id" json:"process_id,string"` // 工序ID
	baize.BaseEntityDQL
}

type SysProProcessContextListData struct {
	Rows  []*SysProProcessContent `json:"rows"`
	Total int64                   `json:"total"`
}
