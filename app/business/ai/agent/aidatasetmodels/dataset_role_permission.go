package aidatasetmodels

import "nova-factory-server/app/baize"

// DatasetRolePermission 知识库/文档-角色权限记录。
type DatasetRolePermission struct {
	ID               int64    `json:"id,string" gorm:"column:id"`
	RoleID           int64    `json:"roleId,string" gorm:"column:role_id"`
	DatasetIDs       string   `json:"-" gorm:"column:dataset_ids"`
	DatasetUuIDs     string   `json:"-" gorm:"column:dataset_uuids"`
	DatasetIDArray   []string `json:"datasetIds" gorm:"-"`
	DocumentIDs      string   `json:"-" gorm:"column:document_ids"`
	DocumentUuIDs    string   `json:"-" gorm:"column:document_uuids"`
	DocumentIDsArray []string `json:"documentIds" gorm:"-"`
	Permission       string   `json:"permission" gorm:"column:permission"`
	Status           bool     `json:"status,string" gorm:"column:status"`
	DeptID           int64    `json:"deptId,string" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// SetDatasetRolePermission 知识库/文档-角色权限设置参数。
type SetDatasetRolePermission struct {
	ID         int64    `json:"id,string"`
	RoleID     int64    `json:"roleId,string" binding:"required"`
	DatasetID  []string `json:"datasetIds,string" binding:"required"`
	DocumentID []string `json:"documentIds,string"`
	Status     bool     `json:"status"`
	Permission string   `json:"permission"`
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
