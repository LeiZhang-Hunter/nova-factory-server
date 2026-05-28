package service

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"

	"github.com/gin-gonic/gin"
)

// IWarehouseLocationService WMS 库位服务接口。
type IWarehouseLocationService interface {
	Set(c *gin.Context, req *models.WarehouseLocationSet) (*models.WarehouseLocation, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.WarehouseLocation, error)
	List(c *gin.Context, req *models.WarehouseLocationQuery) (*models.WarehouseLocationListData, error)
}
