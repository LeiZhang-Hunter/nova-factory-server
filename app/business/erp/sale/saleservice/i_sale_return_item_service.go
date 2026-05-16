package saleservice

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
)

// ISaleReturnItemService ERP 销售退货项服务接口
type ISaleReturnItemService interface {
	Create(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error)
	Update(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturnItem, error)
	List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error)
}
