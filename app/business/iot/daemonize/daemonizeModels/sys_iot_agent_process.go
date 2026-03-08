package daemonizeModels

import (
	"nova-factory-server/app/baize"
	"time"
)

// SysIotAgentProcess agent进程信息
type SysIotAgentProcess struct {
	ID            int32     `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id,string"`            // 自增标识
	AgentObjectID int64     `gorm:"column:agent_object_id;not null;comment:agent uuid" json:"agent_object_id,string"` // agent uuid
	Status        int32     `gorm:"column:status;not null;comment:状态 1-运行 2-停止" json:"status"`                        // 状态 1-运行 2-停止
	Name          string    `gorm:"column:name;not null;comment:名字" json:"name"`                                      // 名字
	Version       string    `gorm:"column:version;not null;comment:版本" json:"version"`                                // 版本
	StartTime     time.Time `gorm:"column:start_time;comment:启动时间" json:"start_time"`                                 // 启动时间
	DeptID        int64     `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                       // 部门ID
	State         bool      `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                 // 操作状态（0正常 -1删除）
	baize.BaseEntity
}
