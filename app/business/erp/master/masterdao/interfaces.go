package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpcrud"
)

// IAccountDao ERP 结算账户数据访问接口
type IAccountDao interface {
	Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Account, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Account, error)
	ListPage(c *gin.Context, req *mastermodels.AccountQuery) (*erpcrud.PageResult[mastermodels.Account], error)
	List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error)
}

// ICustomerDao ERP 客户数据访问接口
type ICustomerDao interface {
	Create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	Update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Customer, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Customer, error)
	ListPage(c *gin.Context, req *mastermodels.CustomerQuery) (*erpcrud.PageResult[mastermodels.Customer], error)
	List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error)
}

// IProductDao ERP 产品数据访问接口
type IProductDao interface {
	Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Product, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Product, error)
	ListPage(c *gin.Context, req *mastermodels.ProductQuery) (*erpcrud.PageResult[mastermodels.Product], error)
	List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error)
}

// IProductCategoryDao ERP 产品分类数据访问接口
type IProductCategoryDao interface {
	Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductCategory, error)
	ListPage(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*erpcrud.PageResult[mastermodels.ProductCategory], error)
	List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error)
}

// IProductUnitDao ERP 产品单位数据访问接口
type IProductUnitDao interface {
	Create(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	Update(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductUnit, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductUnit, error)
	ListPage(c *gin.Context, req *mastermodels.ProductUnitQuery) (*erpcrud.PageResult[mastermodels.ProductUnit], error)
	List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error)
}

// ISupplierDao ERP 供应商数据访问接口
type ISupplierDao interface {
	Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Supplier, error)
	ListPage(c *gin.Context, req *mastermodels.SupplierQuery) (*erpcrud.PageResult[mastermodels.Supplier], error)
	List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error)
}

// IWarehouseDao ERP 仓库数据访问接口
type IWarehouseDao interface {
	Create(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	Update(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Warehouse, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Warehouse, error)
	ListPage(c *gin.Context, req *mastermodels.WarehouseQuery) (*erpcrud.PageResult[mastermodels.Warehouse], error)
	List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error)
}
