package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockInItemServiceImpl 提供 ERP 其它入库单项业务实现。
type StockInItemServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockInItem, stockmodels.StockInItemUpsert, stockmodels.StockInItemQuery]
}

// NewStockInItemService 创建 ERP 其它入库单项服务。
func NewStockInItemService(dao stockdao.IStockInItemDao) stockservice.IStockInItemService {
	return &StockInItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockInItem, stockmodels.StockInItemUpsert, stockmodels.StockInItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_in_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 其它入库单项列表。
func (s *StockInItemServiceImpl) List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockInItemListData{Rows: result.Rows, Total: result.Total}, nil
}
