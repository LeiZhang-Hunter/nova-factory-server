package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IAccountService ERP 结算账户服务接口
type IAccountService interface {
	Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Account, error)
	List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error)
}

// ICustomerService ERP 客户服务接口
type ICustomerService interface {
	Create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	Update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Customer, error)
	List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error)
}

// IProductService ERP 产品服务接口
type IProductService interface {
	Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Product, error)
	List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error)
}

// IProductCategoryService ERP 产品分类服务接口
type IProductCategoryService interface {
	Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error)
	List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error)
}

// IProductUnitService ERP 产品单位服务接口
type IProductUnitService interface {
	Create(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	Update(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductUnit, error)
	List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error)
}

// ISupplierService ERP 供应商服务接口
type ISupplierService interface {
	Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error)
	List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error)
}

// IWarehouseService ERP 仓库服务接口
type IWarehouseService interface {
	Create(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	Update(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Warehouse, error)
	List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error)
}
