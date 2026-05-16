package saledao

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpcrud"
)

// ISaleOutDao ERP 销售出库数据访问接口
type ISaleOutDao interface {
	Create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	Update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOut, error)
	GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleOut, error)
	ListPage(c *gin.Context, req *salemodels.SaleOutQuery) (*erpcrud.PageResult[salemodels.SaleOut], error)
	List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error)
}

// ISaleOutItemDao ERP 销售出库项数据访问接口
type ISaleOutItemDao interface {
	Create(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error)
	Update(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOutItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleOutItem, error)
	ListPage(c *gin.Context, req *salemodels.SaleOutItemQuery) (*erpcrud.PageResult[salemodels.SaleOutItem], error)
	List(c *gin.Context, req *salemodels.SaleOutItemQuery) (*salemodels.SaleOutItemListData, error)
}

// ISaleReturnDao ERP 销售退货数据访问接口
type ISaleReturnDao interface {
	Create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	Update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturn, error)
	GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleReturn, error)
	ListPage(c *gin.Context, req *salemodels.SaleReturnQuery) (*erpcrud.PageResult[salemodels.SaleReturn], error)
	List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error)
}

// ISaleReturnItemDao ERP 销售退货项数据访问接口
type ISaleReturnItemDao interface {
	Create(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error)
	Update(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleReturnItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleReturnItem, error)
	ListPage(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*erpcrud.PageResult[salemodels.SaleReturnItem], error)
	List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error)
}
