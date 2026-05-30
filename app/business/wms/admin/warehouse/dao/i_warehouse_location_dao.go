package dao

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"

	"github.com/gin-gonic/gin"
)

// IWarehouseLocationDao WMS 库位数据访问接口。
type IWarehouseLocationDao interface {
	Set(c *gin.Context, req *models.WarehouseLocationSet) (*models.WarehouseLocation, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.WarehouseLocation, error)
	GetByName(c *gin.Context, warehouseAreaID int64, locationName string) (*models.WarehouseLocation, error)
	GetByTagNumber(c *gin.Context, tagNumber string) (*models.WarehouseLocation, error)
	List(c *gin.Context, req *models.WarehouseLocationQuery) (*models.WarehouseLocationListData, error)
}
