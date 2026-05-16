package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpcrud"
)

// IStockDao ERP 产品库存数据访问接口
type IStockDao interface {
	Create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	Update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.Stock, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.Stock, error)
	ListPage(c *gin.Context, req *stockmodels.StockQuery) (*erpcrud.PageResult[stockmodels.Stock], error)
	List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error)
}

// IStockCheckDao ERP 库存盘点单数据访问接口
type IStockCheckDao interface {
	Create(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	Update(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheck, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockCheck, error)
	ListPage(c *gin.Context, req *stockmodels.StockCheckQuery) (*erpcrud.PageResult[stockmodels.StockCheck], error)
	List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error)
}

// IStockCheckItemDao ERP 库存盘点单项数据访问接口
type IStockCheckItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	Update(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheckItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockCheckItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*erpcrud.PageResult[stockmodels.StockCheckItem], error)
	List(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*stockmodels.StockCheckItemListData, error)
}

// IStockInDao ERP 其它入库单数据访问接口
type IStockInDao interface {
	Create(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	Update(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockIn, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockIn, error)
	ListPage(c *gin.Context, req *stockmodels.StockInQuery) (*erpcrud.PageResult[stockmodels.StockIn], error)
	List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error)
}

// IStockInItemDao ERP 其它入库单项数据访问接口
type IStockInItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	Update(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockInItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockInItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockInItemQuery) (*erpcrud.PageResult[stockmodels.StockInItem], error)
	List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error)
}

// IStockMoveDao ERP 库存调拨单数据访问接口
type IStockMoveDao interface {
	Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockMove, error)
	ListPage(c *gin.Context, req *stockmodels.StockMoveQuery) (*erpcrud.PageResult[stockmodels.StockMove], error)
	List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error)
}

// IStockMoveItemDao ERP 库存调拨单项数据访问接口
type IStockMoveItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error)
	Update(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMoveItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockMoveItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*erpcrud.PageResult[stockmodels.StockMoveItem], error)
	List(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*stockmodels.StockMoveItemListData, error)
}

// IStockOutDao ERP 其它出库单数据访问接口
type IStockOutDao interface {
	Create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	Update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOut, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockOut, error)
	ListPage(c *gin.Context, req *stockmodels.StockOutQuery) (*erpcrud.PageResult[stockmodels.StockOut], error)
	List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error)
}

// IStockOutItemDao ERP 其它出库单项数据访问接口
type IStockOutItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	Update(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOutItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockOutItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockOutItemQuery) (*erpcrud.PageResult[stockmodels.StockOutItem], error)
	List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error)
}

// IStockRecordDao ERP 产品库存明细数据访问接口
type IStockRecordDao interface {
	Create(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	Update(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockRecord, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockRecord, error)
	ListPage(c *gin.Context, req *stockmodels.StockRecordQuery) (*erpcrud.PageResult[stockmodels.StockRecord], error)
	List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error)
}
