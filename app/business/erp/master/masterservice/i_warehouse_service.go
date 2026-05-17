package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IWarehouseService ERP 仓库服务接口
type IWarehouseService interface {
	Create(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	Update(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Warehouse, error)
	List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error)
}
