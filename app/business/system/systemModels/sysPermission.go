package systemModels

import (
	"nova-factory-server/app/baize"
)

type SysPermissionVo struct {
	PermissionId   int64  `json:"permissionId,string" db:"permission_id"` //权限ID
	PermissionName string `json:"permissionName" db:"permission_name"`    //权限名称
	ParentId       int64  `json:"parentId,string" db:"parent_id"`         //父权限ID
	Permission     string `json:"permission" db:"permission"`             //权限标识符
	Sort           int    `json:"sort" db:"sort"`                         //排序
	Status         string `json:"status" db:"status"`                     // 状态
	baize.BaseEntity
}
type SysPermissionAdd struct {
	PermissionId   int64  `json:"permissionId,string" db:"permission_id" swaggerignore:"true"` //权限ID
	PermissionName string `json:"permissionName" db:"permission_name"`                         //权限名称
	ParentId       int64  `json:"parentId,string" db:"parent_id"`                              //父权限ID
	Permission     string `json:"permission" db:"permission" binding:"required"`               //权限标识符
	Sort           int    `json:"sort" db:"sort"`                                              //排序
	Status         string `json:"status" db:"status"`                                          // 状态
	baize.BaseEntity
}

type SysPermissionEdit struct {
	PermissionId   int64  `json:"permissionId,string" db:"permission_id" binding:"required"` //权限ID
	PermissionName string `json:"permissionName" db:"permission_name"`                       //权限名称
	ParentId       int64  `json:"parentId,string" db:"parent_id"`                            //父权限ID
	Permission     string `json:"permission" db:"permission"`                                //权限标识符
	Sort           *int   `json:"sort" db:"sort"`                                            //排序
	Status         string `json:"status" db:"status"`                                        // 状态
	baize.BaseEntity
}

type SysPermissionDQL struct {
	Status string `form:"status" db:"status"` // 状态
	baize.BaseEntity
}
