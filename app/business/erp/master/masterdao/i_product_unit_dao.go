package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IProductUnitDao ERP 产品单位数据访问接口
type IProductUnitDao interface {
	Create(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	Update(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductUnit, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductUnit, error)
	GetByIDs(c *gin.Context, ids []int64) ([]*mastermodels.ProductUnit, error)
	ListPage(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error)
	List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error)
}
