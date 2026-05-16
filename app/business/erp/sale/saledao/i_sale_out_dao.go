package saledao

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// ISaleOutDao ERP 销售出库数据访问接口
type ISaleOutDao interface {
	Create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	Update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOut, error)
	GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleOut, error)
	ListPage(c *gin.Context, req *salemodels.SaleOutQuery) (*erpbiz.PageResult[salemodels.SaleOut], error)
	List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error)
}
