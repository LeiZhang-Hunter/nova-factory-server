package dao

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"

	"github.com/gin-gonic/gin"
)

// IWarehouseAreaDao WMS 库区数据访问接口。
type IWarehouseAreaDao interface {
	Set(c *gin.Context, req *models.WarehouseAreaSet) (*models.WarehouseArea, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.WarehouseArea, error)
	GetByName(c *gin.Context, warehouseID, parentID int64, areaName string) (*models.WarehouseArea, error)
	List(c *gin.Context, req *models.WarehouseAreaQuery) (*models.WarehouseAreaListData, error)
}
