package aidatasetserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"

	"github.com/gin-gonic/gin"
)

type IDatasetRolePermissionServiceImpl struct {
	dao aidatasetdao.IDatasetRolePermissionDao
}

// NewIDatasetRolePermissionServiceImpl 知识库/文档-角色权限服务构造函数。
func NewIDatasetRolePermissionServiceImpl(dao aidatasetdao.IDatasetRolePermissionDao) aidatasetservice.IDatasetRolePermissionService {
	return &IDatasetRolePermissionServiceImpl{dao: dao}
}

// List 查询知识库/文档-角色权限列表。
func (i *IDatasetRolePermissionServiceImpl) List(c *gin.Context, req *aidatasetmodels.DatasetRolePermissionQuery) (*aidatasetmodels.DatasetRolePermissionListData, error) {
	return i.dao.List(c, req)
}

// Set 新增或修改知识库/文档-角色权限。
func (i *IDatasetRolePermissionServiceImpl) Set(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.RoleID == 0 {
		return nil, errors.New("roleId不能为空")
	}
	if len(req.DatasetID) == 0 {
		return nil, errors.New("datasetId不能为空")
	}
	perm := strings.TrimSpace(req.Permission)
	if perm == "" {
		perm = "read"
	}
	switch perm {
	case "read", "write", "admin":
	default:
		return nil, errors.New("permission不合法")
	}
	req.Permission = perm
	if req.ID > 0 {
		return i.dao.Update(c, req)
	}
	return i.dao.Create(c, req)
}

// Remove 删除知识库/文档-角色权限记录。
func (i *IDatasetRolePermissionServiceImpl) Remove(c *gin.Context, ids []int64) error {
	return i.dao.Remove(c, ids)
}
