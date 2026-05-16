package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockMoveServiceImpl 提供 ERP 库存调拨单业务实现。
type StockMoveServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockMove, stockmodels.StockMoveUpsert, stockmodels.StockMoveQuery]
}

// NewStockMoveService 创建 ERP 库存调拨单服务。
func NewStockMoveService(dao stockdao.IStockMoveDao) stockservice.IStockMoveService {
	return &StockMoveServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockMove, stockmodels.StockMoveUpsert, stockmodels.StockMoveQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_move",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "调拨单号"}},
		}),
	}
}

// List 查询 ERP 库存调拨单列表。
func (s *StockMoveServiceImpl) List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveListData{Rows: result.Rows, Total: result.Total}, nil
}
