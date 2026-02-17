package configurationModels

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

// SysConfiguration mapped from table <sys_configuration>
type SysConfiguration struct {
	ID          int64  `gorm:"column:id;primaryKey;comment:组态id" json:"id,string"` // 组态id
	Name        string `gorm:"column:name;not null" json:"name"`                   // 组态名称
	Tag         string `gorm:"column:tag;not null" json:"tag"`                     // 标签
	Annotation  string `gorm:"column:annotation" json:"annotation"`                // 注解
	Description string `gorm:"column:description;not null" json:"description"`     // 描述
	DeptID      int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`         // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

// SetSysConfiguration 创建/编辑组态入参
type SetSysConfiguration struct {
	ID          int64                     `gorm:"column:id;primaryKey;comment:组态id" json:"id,string"`  // 组态id
	Name        string                    `gorm:"column:name;not null" binding:"required" json:"name"` // 组态名称
	Tag         string                    `gorm:"column:tag;not null" binding:"required" json:"tag"`   // 标签
	Annotation  []deviceModels.Annotation `gorm:"column:annotation" json:"annotation"`                 // 注解
	Description string                    `gorm:"column:description;not null" json:"description"`      // 描述
}

// SysConfigurationReq 组态查询参数
type SysConfigurationReq struct {
	Name string `gorm:"column:name;not null" form:"name"`
	Tag  string `gorm:"column:tag;not null" form:"tag"`
	baize.BaseEntityDQL
}

// ToSysConfiguration 转换为持久化实体
func ToSysConfiguration(set *SetSysConfiguration) *SysConfiguration {
	return &SysConfiguration{
		ID:   set.ID,
		Name: set.Name,
		Tag:  set.Tag,
		//Annotation:  set.Annotation,
		Description: set.Description,
	}
}

// SysConfigurationList 组态列表
type SysConfigurationList struct {
	Rows  []*SysConfiguration `json:"rows"`
	Total int64               `json:"total"`
}
