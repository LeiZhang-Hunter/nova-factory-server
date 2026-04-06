package aidatasetservice

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"github.com/gin-gonic/gin"
)

type IDatasetRolePermissionService interface {
	List(c *gin.Context, req *aidatasetmodels.DatasetRolePermissionQuery) (*aidatasetmodels.DatasetRolePermissionListData, error)
	Set(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error)
	Remove(c *gin.Context, ids []int64) error
}
