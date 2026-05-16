package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IProductUnitService ERP 产品单位服务接口
type IProductUnitService interface {
	Create(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	Update(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductUnit, error)
	List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error)
}
