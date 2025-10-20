package systemModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/constant/device"
)

type SysDeviceElectricSetting struct {
	ID             int64       `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                    // id
	DeviceID       int64       `gorm:"column:device_id;not null;comment:device_id" json:"device_id,string"` // device_id
	Name           string      `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                     // 智能预警名称
	DevName        string      `gorm:"-" json:"dev_name"`                                                   // 设备名字
	Expression     string      `gorm:"column:expression;comment:表达式" json:"-"`                              // 表达式
	ExpressionData *Expression `gorm:"-" json:"expression"`                                                 // 表达式
	DeptID         int64       `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                          // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SysDeviceElectricSettingVO struct {
	ID         int64       `gorm:"column:id;primaryKey;comment:id" json:"id,string"`                    // id
	DeviceID   int64       `gorm:"column:device_id;not null;comment:device_id" json:"device_id,string"` // device_id
	Name       string      `gorm:"column:name;not null;comment:智能预警名称" json:"name"`                     // 智能预警名称
	Expression *Expression `gorm:"column:expression;comment:表达式" json:"expression"`                     // 表达式
}

func ToSysDeviceElectricSetting(vo *SysDeviceElectricSettingVO) *SysDeviceElectricSetting {
	var expression []byte
	var err error
	if vo.Expression != nil {
		expression, err = json.Marshal(vo.Expression)
		if err != nil {
			zap.L().Error("json encode error", zap.Error(err))
		}
	} else {
		expression = make([]byte, 0)
	}
	return &SysDeviceElectricSetting{
		ID:         vo.ID,
		DeviceID:   vo.DeviceID,
		Name:       vo.Name,
		Expression: string(expression),
	}
}

type SysDeviceElectricSettingDQL struct {
	DeviceId []string `form:"device_id"`
	baize.BaseEntityDQL
}

type SysDeviceElectricSettingData struct {
	Rows  []*SysDeviceElectricSetting `json:"rows"`
	Total int64                       `json:"total"`
}

type Expression struct {
	Rules []Rule `yaml:"rules,omitempty" json:"rules"`
}

type Rule struct {
	MatchType string            `yaml:"matchType,omitempty" json:"matchType"`
	Groups    []Group           `yaml:"groups,omitempty" json:"groups"`
	RunStatus device.RUN_STATUS `yaml:"run_status,omitempty" json:"run_status,string"`
}

type Group struct {
	Key          string `yaml:"key,omitempty" json:"key"`
	Name         string `yaml:"name,omitempty" json:"name"`
	Operator     string `yaml:"operator,omitempty" json:"operator"`
	OperatorName string `yaml:"operatorName,omitempty" json:"operatorName"`
	Value        string `yaml:"value,omitempty" json:"value"`
}
