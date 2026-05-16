package saleservice

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
)

// ISaleOutItemService ERP 销售出库项服务接口
type ISaleOutItemService interface {
	Create(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error)
	Update(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOutItem, error)
	List(c *gin.Context, req *salemodels.SaleOutItemQuery) (*salemodels.SaleOutItemListData, error)
}
