package aidatasetdao

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"github.com/gin-gonic/gin"
)

type IDatasetRolePermissionDao interface {
	Create(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error)
	Update(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error)
	List(c *gin.Context, req *aidatasetmodels.DatasetRolePermissionQuery) (*aidatasetmodels.DatasetRolePermissionListData, error)
	Remove(c *gin.Context, ids []int64) error
}
