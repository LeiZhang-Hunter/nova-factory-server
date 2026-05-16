package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockCheckItemServiceImpl 提供 ERP 库存盘点单项业务实现。
type StockCheckItemServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockCheckItem, stockmodels.StockCheckItemUpsert, stockmodels.StockCheckItemQuery]
}

// NewStockCheckItemService 创建 ERP 库存盘点单项服务。
func NewStockCheckItemService(dao stockdao.IStockCheckItemDao) stockservice.IStockCheckItemService {
	return &StockCheckItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockCheckItem, stockmodels.StockCheckItemUpsert, stockmodels.StockCheckItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_check_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 库存盘点单项列表。
func (s *StockCheckItemServiceImpl) List(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*stockmodels.StockCheckItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockCheckItemListData{Rows: result.Rows, Total: result.Total}, nil
}
