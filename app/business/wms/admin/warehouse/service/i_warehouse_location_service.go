package service

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"

	"github.com/gin-gonic/gin"
)

// IWarehouseLocationService WMS 库位服务接口。
type IWarehouseLocationService interface {
	// Set 新增或修改 WMS 库位。
	Set(c *gin.Context, req *models.WarehouseLocationSet) (*models.WarehouseLocation, error)
	// DeleteByIDs 删除 WMS 库位。
	DeleteByIDs(c *gin.Context, ids []int64) error
	// GetByID 查询 WMS 库位详情。
	GetByID(c *gin.Context, id int64) (*models.WarehouseLocation, error)
	// List 查询 WMS 库位列表。
	List(c *gin.Context, req *models.WarehouseLocationQuery) (*models.WarehouseLocationListData, error)
}
