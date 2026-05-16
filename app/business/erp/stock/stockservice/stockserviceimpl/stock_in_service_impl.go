package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockInServiceImpl 提供 ERP 其它入库单业务实现。
type StockInServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockIn, stockmodels.StockInUpsert, stockmodels.StockInQuery]
}

// NewStockInService 创建 ERP 其它入库单服务。
func NewStockInService(dao stockdao.IStockInDao) stockservice.IStockInService {
	return &StockInServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockIn, stockmodels.StockInUpsert, stockmodels.StockInQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_in",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "入库单号"}},
		}),
	}
}

// List 查询 ERP 其它入库单列表。
func (s *StockInServiceImpl) List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockInListData{Rows: result.Rows, Total: result.Total}, nil
}
