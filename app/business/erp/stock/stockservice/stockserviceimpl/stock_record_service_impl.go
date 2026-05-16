package stockserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockRecordServiceImpl 提供 ERP 产品库存明细业务实现。
type StockRecordServiceImpl struct {
	*erpcrud.CRUDService[stockmodels.StockRecord, stockmodels.StockRecordUpsert, stockmodels.StockRecordQuery]
}

// NewStockRecordService 创建 ERP 产品库存明细服务。
func NewStockRecordService(dao stockdao.IStockRecordDao) stockservice.IStockRecordService {
	return &StockRecordServiceImpl{
		CRUDService: erpcrud.NewCRUDService[stockmodels.StockRecord, stockmodels.StockRecordUpsert, stockmodels.StockRecordQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_stock_record",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品库存明细列表。
func (s *StockRecordServiceImpl) List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockRecordListData{Rows: result.Rows, Total: result.Total}, nil
}
