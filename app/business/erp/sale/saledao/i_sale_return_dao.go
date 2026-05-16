package saledao

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// ISaleReturnDao ERP 销售退货数据访问接口
type ISaleReturnDao interface {
	Create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	Update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturn, error)
	GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleReturn, error)
	ListPage(c *gin.Context, req *salemodels.SaleReturnQuery) (*erpbiz.PageResult[salemodels.SaleReturn], error)
	List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error)
}
