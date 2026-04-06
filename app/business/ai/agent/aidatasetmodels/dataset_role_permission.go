package aidatasetmodels

import "nova-factory-server/app/baize"

// DatasetRolePermission 知识库/文档-角色权限记录。
type DatasetRolePermission struct {
	ID         int64  `json:"id,string" gorm:"column:id"`
	RoleID     int64  `json:"roleId,string" gorm:"column:role_id"`
	DatasetID  int64  `json:"datasetId,string" gorm:"column:dataset_id"`
	DocumentID int64  `json:"documentId,string" gorm:"column:document_id"`
	Permission string `json:"permission" gorm:"column:permission"`
	DeptID     int64  `json:"deptId,string" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// SetDatasetRolePermission 知识库/文档-角色权限设置参数。
type SetDatasetRolePermission struct {
	ID         int64  `json:"id,string"`
	RoleID     int64  `json:"roleId,string" binding:"required"`
	DatasetID  int64  `json:"datasetId,string" binding:"required"`
	DocumentID int64  `json:"documentId,string"`
	Permission string `json:"permission"`
}

// DatasetRolePermissionQuery 知识库/文档-角色权限列表查询参数。
type DatasetRolePermissionQuery struct {
	RoleID     int64  `form:"roleId,string"`
	DatasetID  int64  `form:"datasetId,string"`
	DocumentID int64  `form:"documentId,string"`
	Permission string `form:"permission"`
	baize.BaseEntityDQL
}

// DatasetRolePermissionListData 知识库/文档-角色权限列表数据。
type DatasetRolePermissionListData struct {
	Rows  []*DatasetRolePermission `json:"rows"`
	Total int64                    `json:"total"`
}
