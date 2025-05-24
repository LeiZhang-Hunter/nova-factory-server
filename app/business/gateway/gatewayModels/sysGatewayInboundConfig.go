package gatewayModels

import (
	"nova-factory-server/app/baize"
)

type SysGatewayInboundConfig struct {
	GatewayConfigID int64  `gorm:"column:gateway_config_id;primaryKey;comment:设备主键" json:"gateway_config_id"` // 设备主键
	Name            string `gorm:"column:name;comment:设备名称" json:"name"`                                      // 设备名称
	ProtocolType    string `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"`                    // 协议类型
	Config          string `gorm:"column:config;comment:配置" json:"config"`                                    // 配置
	DeptID          int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                // 部门ID
	Status          int32  `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`                         // 操作状态（0正常 1异常）
	State           bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                          // 操作状态（0正常 -1删除）
	CreateUserName  string `json:"createUserName" gorm:"-"`
	UpdateUserName  string `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

func NewSysGatewayInboundConfig(set *SysSetGatewayInboundConfig) *SysGatewayInboundConfig {
	return &SysGatewayInboundConfig{
		GatewayConfigID: set.GatewayConfigID,
		Name:            set.Name,
		ProtocolType:    set.ProtocolType,
		Config:          set.Config,
		DeptID:          set.GatewayConfigID,
		Status:          set.Status,
	}
}

type SysSetGatewayInboundConfig struct {
	GatewayConfigID int64  `gorm:"column:gateway_config_id;primaryKey;comment:设备主键" json:"gateway_config_id"` // 设备主键
	Name            string `gorm:"column:name;comment:设备名称" json:"name"`                                      // 设备名称
	ProtocolType    string `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"`                    // 协议类型
	Config          string `gorm:"column:config;comment:配置" json:"config"`                                    // 配置
	Status          int32  `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`                         // 操作状态（0正常 1异常）
	State           bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                          // 操作状态（0正常 -1删除）
}

type SysSetGatewayInboundConfigReq struct {
	Name         string `gorm:"column:name;comment:设备名称" json:"name"`                   // 设备名称
	ProtocolType string `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"` // 协议类型
	Status       int32  `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`      // 操作状态（0正常 1异常）
	baize.BaseEntityDQL
}

type SysSetGatewayInboundConfigList struct {
	Rows  []*SysGatewayInboundConfig
	Total int64 `json:"total"`
	baize.BaseEntityDQL
}
