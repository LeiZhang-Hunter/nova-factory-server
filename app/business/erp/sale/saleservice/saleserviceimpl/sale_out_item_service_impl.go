package saleserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"

	"github.com/gin-gonic/gin"
)

// SaleOutItemServiceImpl 提供 ERP 销售出库项业务实现。
type SaleOutItemServiceImpl struct {
	*erpcrud.CRUDService[salemodels.SaleOutItem, salemodels.SaleOutItemUpsert, salemodels.SaleOutItemQuery]
}

// NewSaleOutItemService 创建 ERP 销售出库项服务。
func NewSaleOutItemService(dao saledao.ISaleOutItemDao) saleservice.ISaleOutItemService {
	return &SaleOutItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[salemodels.SaleOutItem, salemodels.SaleOutItemUpsert, salemodels.SaleOutItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_sale_out_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 销售出库项列表。
func (s *SaleOutItemServiceImpl) List(c *gin.Context, req *salemodels.SaleOutItemQuery) (*salemodels.SaleOutItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutItemListData{Rows: result.Rows, Total: result.Total}, nil
}
