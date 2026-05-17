package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IProductService ERP 产品服务接口
type IProductService interface {
	Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Product, error)
	List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error)
}
