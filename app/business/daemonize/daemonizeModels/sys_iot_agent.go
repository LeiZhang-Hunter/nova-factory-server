package daemonizeModels

import (
	"nova-factory-server/app/baize"
	"time"
)

// SysIotAgent agent信息
type SysIotAgent struct {
	ObjectID          uint64     `gorm:"column:object_id;primaryKey;comment:agent uuid" json:"object_id"`                           // agent uuid
	Name              string     `gorm:"column:name;not null;comment:agent名字" json:"name"`                                          // agent名字
	Username          string     `gorm:"column:username;not null;comment:username" json:"username"`                                 // username
	Password          string     `gorm:"column:password;not null;comment:password" json:"password"`                                 // password
	OperateState      int32      `gorm:"column:operate_state;not null;comment:操作状态 1-启动中 2-停止中 3-启动失败 4-停止失败" json:"operate_state"` // 操作状态 1-启动中 2-停止中 3-启动失败 4-停止失败
	OperateTime       *time.Time `gorm:"column:operate_time;comment:操作时间" json:"operate_time,omitemptyZ"`                           // 操作时间
	Version           string     `gorm:"column:version;not null;comment:agent版本" json:"version"`                                    // agent版本
	ConfigUUID        string     `gorm:"column:config_uuid;not null;comment:配置版本" json:"config_uuid"`                               // 配置版本
	Ipv4              string     `gorm:"column:ipv4;not null;comment:ipv4地址" json:"ipv4"`                                           // ipv4地址
	Ipv6              string     `gorm:"column:ipv6;not null;comment:ipv6地址" json:"ipv6"`                                           // ipv6地址
	LastHeartbeatTime *time.Time `gorm:"column:last_heartbeat_time;comment:上次心跳时间" json:"last_heartbeat_time"`                      // 上次心跳时间
	UpdateConfigTime  *time.Time `gorm:"column:update_config_time;comment:更新配置时间" json:"update_config_time"`                        // 更新配置时间
	DeptID            int64      `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                                // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func ToSysIotAgent(set *SysIotAgentSetReq) *SysIotAgent {
	return &SysIotAgent{
		ObjectID:   set.ObjectID,
		Name:       set.Name,
		Username:   set.Username,
		Password:   set.Password,
		Version:    set.Version,
		ConfigUUID: set.ConfigUUID,
	}
}

type SysIotAgentSetReq struct {
	ObjectID   uint64 `gorm:"column:object_id;primaryKey;comment:agent uuid" json:"object_id"` // agent uuid
	Name       string `gorm:"column:name;not null;comment:agent名字" json:"name"`                // agent名字
	Username   string `gorm:"column:username;not null;comment:username" json:"username"`       // username
	Password   string `gorm:"column:password;not null;comment:password" json:"password"`       // password
	Version    string `gorm:"column:version;not null;comment:agent版本" json:"version"`          // agent版本
	ConfigUUID string `gorm:"column:config_uuid;not null;comment:配置版本" json:"config_uuid"`     // 配置版本
}

type SysIotAgentListReq struct {
	baize.BaseEntityDQL
}

type SysIotAgentListData struct {
	Rows  []*SysIotAgent `json:"rows"`
	Total int64          `json:"total"`
}
