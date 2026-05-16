package saleserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"

	"github.com/gin-gonic/gin"
)

// SaleReturnItemServiceImpl 提供 ERP 销售退货项业务实现。
type SaleReturnItemServiceImpl struct {
	*erpcrud.CRUDService[salemodels.SaleReturnItem, salemodels.SaleReturnItemUpsert, salemodels.SaleReturnItemQuery]
}

// NewSaleReturnItemService 创建 ERP 销售退货项服务。
func NewSaleReturnItemService(dao saledao.ISaleReturnItemDao) saleservice.ISaleReturnItemService {
	return &SaleReturnItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[salemodels.SaleReturnItem, salemodels.SaleReturnItemUpsert, salemodels.SaleReturnItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_sale_return_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 销售退货项列表。
func (s *SaleReturnItemServiceImpl) List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
