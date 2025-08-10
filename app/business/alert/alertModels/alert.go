package alertModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
)

// SysAlert 告警策略配置
type SysAlert struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"`     // 自增标识
	GatewayID   int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id,string"`   // 网关id
	TemplateID  int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id,string"` // 模板id
	ActionId    int64  `gorm:"column:action_id;not null;comment:处理id" json:"action_id,string"`     // 处理id
	ReasonId    int64  `gorm:"column:reason_id;not null;comment:推理id" json:"reason_id,string"`     // 推理id
	Name        string `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                    // 告警策略名称
	Additions   string `gorm:"column:additions;comment:注解" json:"additions"`                       // 注解
	Advanced    string `gorm:"column:advanced;comment:告警规则" json:"advanced"`                       // 告警规则
	Ignore      string `gorm:"column:ignore;comment:忽略规则" json:"ignore"`                           // 忽略规则
	Matcher     string `gorm:"column:matcher;comment:匹配规则" json:"matcher"`                         // 匹配规则
	Description string `gorm:"column:description;not null;comment:配置版本" json:"description"`        // 配置版本
	DeptID      int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                         // 部门ID
	baize.BaseEntity
	State  bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`  // 操作状态（0正常 -1删除）
	Status bool `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"` // 操作状态（0正常 1异常）
}

func FromSysAlertToSetData(data *SysAlert) *SetSysAlert {
	additions := make(map[string]string)
	if data.Additions != "" {
		err := json.Unmarshal([]byte(data.Additions), &additions)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
		}
	}

	var advanced *logalert.Advanced
	if data.Advanced != "" {
		err := json.Unmarshal([]byte(data.Advanced), &advanced)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
		}
	}

	var ignore []logalert.DeviceMetric = make([]logalert.DeviceMetric, 0)
	if len(data.Ignore) != 0 {
		err := json.Unmarshal([]byte(data.Ignore), &ignore)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
		}
	}

	var matcher *logalert.Matcher
	if data.Matcher != "" {
		err := json.Unmarshal([]byte(data.Matcher), &matcher)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
		}
	}
	return &SetSysAlert{
		ID:          data.ID,
		GatewayID:   data.GatewayID,
		TemplateID:  data.TemplateID,
		ActionId:    data.ActionId,
		ReasonId:    data.ReasonId,
		Name:        data.Name,
		Additions:   additions,
		Advanced:    advanced,
		Ignore:      ignore,
		Matcher:     matcher,
		Description: data.Description,
		Status:      data.Status,
	}
}

type SetSysAlert struct {
	ID          int64                   `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id,string"` // 自增标识
	GatewayID   int64                   `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id,string"`      // 网关id
	TemplateID  int64                   `gorm:"column:template_id;not null;comment:模板id" json:"template_id,string"`    // 模板id
	ActionId    int64                   `gorm:"column:action_id;not null;comment:处理id" json:"action_id,string"`        // 处理id
	ReasonId    int64                   `gorm:"column:reason_id;not null;comment:推理id" json:"reason_id,string"`        // 推理id
	Name        string                  `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                       // 告警策略名称
	Additions   map[string]string       `gorm:"column:additions;comment:注解" json:"additions"`                          // 注解
	Advanced    *logalert.Advanced      `gorm:"column:advanced;comment:告警规则" json:"advanced"`                          // 告警规则
	Ignore      []logalert.DeviceMetric `gorm:"column:ignore;comment:忽略规则" json:"ignore"`
	Matcher     *logalert.Matcher       `gorm:"column:matcher;comment:匹配规则" json:"matcher"`
	Description string                  `gorm:"column:description;not null;comment:描述"  json:"description"` // 配置版本
	Status      bool                    `json:"status"`
}

func ToSysAlert(data *SetSysAlert) *SysAlert {
	var additions []byte
	var err error
	if len(data.Additions) != 0 {
		additions, err = json.Marshal(data.Additions)
		if err != nil {
			zap.L().Error("json marshal error", zap.Error(err))
		}
	}
	var advanced []byte
	if data.Advanced != nil {
		advanced, err = json.Marshal(data.Advanced)
		if zap.L() != nil {
			zap.L().Error("json marshal error", zap.Error(err))
		}
	}

	var ignore []byte
	if data.Ignore != nil {
		ignore, err = json.Marshal(data.Ignore)
		if zap.L() != nil {
			zap.L().Error("json marshal error", zap.Error(err))
		}
	}

	var matcher []byte
	if data.Matcher != nil {
		matcher, err = json.Marshal(data.Matcher)
		if zap.L() != nil {
			zap.L().Error("json marshal error", zap.Error(err))
		}
	}
	return &SysAlert{
		ID:          data.ID,
		GatewayID:   data.GatewayID,
		TemplateID:  data.TemplateID,
		ReasonId:    data.ReasonId,
		ActionId:    data.ActionId,
		Name:        data.Name,
		Additions:   string(additions),
		Advanced:    string(advanced),
		Description: data.Description,
		Ignore:      string(ignore),
		Matcher:     string(matcher),
		Status:      data.Status,
	}
}

type SysAlertListReq struct {
	GatewayID  int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`   // 网关id
	TemplateID int64  `gorm:"column:template_id;not null;comment:模板id" json:"template_id"` // 模板id
	Name       string `gorm:"column:name;not null;comment:告警策略名称" json:"name"`             // 告警策略名称
	Status     *bool  `gorm:"column:status;comment:启用状态" json:"status"`
	baize.BaseEntityDQL
}

type SysAlertList struct {
	Rows  []*SetSysAlert `json:"rows"`
	Total uint64         `json:"total"`
}

type ChangeSysAlert struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"` // 自增标识
	Status bool  `gorm:"column:status;comment:启用状态" json:"status"`
}
