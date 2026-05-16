package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockOutServiceImpl 提供 ERP 其它出库单业务实现。
type StockOutServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockOut, stockmodels.StockOutUpsert, stockmodels.StockOutQuery]
}

// NewStockOutService 创建 ERP 其它出库单服务。
func NewStockOutService(dao stockdao.IStockOutDao) stockservice.IStockOutService {
	return &StockOutServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockOut, stockmodels.StockOutUpsert, stockmodels.StockOutQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_out",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "出库单号"}},
		}),
	}
}

// List 查询 ERP 其它出库单列表。
func (s *StockOutServiceImpl) List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutListData{Rows: result.Rows, Total: result.Total}, nil
}
