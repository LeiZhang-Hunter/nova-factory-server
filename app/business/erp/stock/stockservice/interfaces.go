package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockService ERP 产品库存服务接口
type IStockService interface {
	Create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	Update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.Stock, error)
	List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error)
}

// IStockCheckService ERP 库存盘点单服务接口
type IStockCheckService interface {
	Create(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	Update(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheck, error)
	List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error)
}

// IStockCheckItemService ERP 库存盘点单项服务接口
type IStockCheckItemService interface {
	Create(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	Update(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheckItem, error)
	List(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*stockmodels.StockCheckItemListData, error)
}

// IStockInService ERP 其它入库单服务接口
type IStockInService interface {
	Create(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	Update(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockIn, error)
	List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error)
}

// IStockInItemService ERP 其它入库单项服务接口
type IStockInItemService interface {
	Create(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	Update(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockInItem, error)
	List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error)
}

// IStockMoveService ERP 库存调拨单服务接口
type IStockMoveService interface {
	Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error)
	List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error)
}

// IStockMoveItemService ERP 库存调拨单项服务接口
type IStockMoveItemService interface {
	Create(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error)
	Update(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMoveItem, error)
	List(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*stockmodels.StockMoveItemListData, error)
}

// IStockOutService ERP 其它出库单服务接口
type IStockOutService interface {
	Create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	Update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOut, error)
	List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error)
}

// IStockOutItemService ERP 其它出库单项服务接口
type IStockOutItemService interface {
	Create(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	Update(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOutItem, error)
	List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error)
}

// IStockRecordService ERP 产品库存明细服务接口
type IStockRecordService interface {
	Create(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	Update(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockRecord, error)
	List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error)
}
