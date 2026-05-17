package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IWarehouseDao ERP 仓库数据访问接口
type IWarehouseDao interface {
	Create(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	Update(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Warehouse, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Warehouse, error)
	ListPage(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error)
	List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error)
}
