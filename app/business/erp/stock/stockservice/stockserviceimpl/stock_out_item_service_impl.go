package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockOutItemServiceImpl 提供 ERP 其它出库单项业务实现。
type StockOutItemServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockOutItem, stockmodels.StockOutItemUpsert, stockmodels.StockOutItemQuery]
}

// NewStockOutItemService 创建 ERP 其它出库单项服务。
func NewStockOutItemService(dao stockdao.IStockOutItemDao) stockservice.IStockOutItemService {
	return &StockOutItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockOutItem, stockmodels.StockOutItemUpsert, stockmodels.StockOutItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_out_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 其它出库单项列表。
func (s *StockOutItemServiceImpl) List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutItemListData{Rows: result.Rows, Total: result.Total}, nil
}
