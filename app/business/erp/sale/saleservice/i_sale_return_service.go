package saleservice

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
)

// ISaleReturnService ERP 销售退货服务接口
type ISaleReturnService interface {
	Create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	Update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturn, error)
	List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error)
}
