package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockServiceImpl 提供 ERP 产品库存业务实现。
type StockServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.Stock, stockmodels.StockUpsert, stockmodels.StockQuery]
}

// NewStockService 创建 ERP 产品库存服务。
func NewStockService(dao stockdao.IStockDao) stockservice.IStockService {
	return &StockServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.Stock, stockmodels.StockUpsert, stockmodels.StockQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品库存列表。
func (s *StockServiceImpl) List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockListData{Rows: result.Rows, Total: result.Total}, nil
}
