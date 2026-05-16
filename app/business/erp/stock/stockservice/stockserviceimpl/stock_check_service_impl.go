package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockCheckServiceImpl 提供 ERP 库存盘点单业务实现。
type StockCheckServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockCheck, stockmodels.StockCheckUpsert, stockmodels.StockCheckQuery]
}

// NewStockCheckService 创建 ERP 库存盘点单服务。
func NewStockCheckService(dao stockdao.IStockCheckDao) stockservice.IStockCheckService {
	return &StockCheckServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockCheck, stockmodels.StockCheckUpsert, stockmodels.StockCheckQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_check",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "盘点单号"}},
		}),
	}
}

// List 查询 ERP 库存盘点单列表。
func (s *StockCheckServiceImpl) List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockCheckListData{Rows: result.Rows, Total: result.Total}, nil
}
