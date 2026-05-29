package service

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"

	"github.com/gin-gonic/gin"
)

// IWarehouseAreaService WMS 库区服务接口。
type IWarehouseAreaService interface {
	// Set 新增或修改 WMS 库区。
	Set(c *gin.Context, req *models.WarehouseAreaSet) (*models.WarehouseArea, error)
	// DeleteByIDs 删除 WMS 库区。
	DeleteByIDs(c *gin.Context, ids []int64) error
	// GetByID 查询 WMS 库区详情。
	GetByID(c *gin.Context, id int64) (*models.WarehouseArea, error)
	// List 查询 WMS 库区列表。
	List(c *gin.Context, req *models.WarehouseAreaQuery) (*models.WarehouseAreaListData, error)
}
