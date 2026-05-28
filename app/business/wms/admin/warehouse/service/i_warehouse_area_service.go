package service

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"

	"github.com/gin-gonic/gin"
)

// IWarehouseAreaService WMS 库区服务接口。
type IWarehouseAreaService interface {
	Set(c *gin.Context, req *models.WarehouseAreaSet) (*models.WarehouseArea, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.WarehouseArea, error)
	List(c *gin.Context, req *models.WarehouseAreaQuery) (*models.WarehouseAreaListData, error)
}
