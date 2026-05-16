package saleservice

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
)

// ISaleOutService ERP 销售出库服务接口
type ISaleOutService interface {
	Create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	Update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOut, error)
	List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error)
}

// ISaleOutItemService ERP 销售出库项服务接口
type ISaleOutItemService interface {
	Create(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error)
	Update(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOutItem, error)
	List(c *gin.Context, req *salemodels.SaleOutItemQuery) (*salemodels.SaleOutItemListData, error)
}

// ISaleReturnService ERP 销售退货服务接口
type ISaleReturnService interface {
	Create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	Update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturn, error)
	List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error)
}

// ISaleReturnItemService ERP 销售退货项服务接口
type ISaleReturnItemService interface {
	Create(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error)
	Update(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturnItem, error)
	List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error)
}
