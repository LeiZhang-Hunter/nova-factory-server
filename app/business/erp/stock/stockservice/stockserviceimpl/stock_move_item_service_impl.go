package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockMoveItemServiceImpl 提供 ERP 库存调拨单项业务实现。
type StockMoveItemServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockMoveItem, stockmodels.StockMoveItemUpsert, stockmodels.StockMoveItemQuery]
}

// NewStockMoveItemService 创建 ERP 库存调拨单项服务。
func NewStockMoveItemService(dao stockdao.IStockMoveItemDao) stockservice.IStockMoveItemService {
	return &StockMoveItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockMoveItem, stockmodels.StockMoveItemUpsert, stockmodels.StockMoveItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_move_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 库存调拨单项列表。
func (s *StockMoveItemServiceImpl) List(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*stockmodels.StockMoveItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveItemListData{Rows: result.Rows, Total: result.Total}, nil
}
