package saleserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"

	"github.com/gin-gonic/gin"
)

// SaleReturnServiceImpl 提供 ERP 销售退货业务实现。
type SaleReturnServiceImpl struct {
	*erpcrud.CRUDService[salemodels.SaleReturn, salemodels.SaleReturnUpsert, salemodels.SaleReturnQuery]
}

// NewSaleReturnService 创建 ERP 销售退货服务。
func NewSaleReturnService(dao saledao.ISaleReturnDao) saleservice.ISaleReturnService {
	return &SaleReturnServiceImpl{
		CRUDService: erpcrud.NewCRUDService[salemodels.SaleReturn, salemodels.SaleReturnUpsert, salemodels.SaleReturnQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_sale_return",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "销售退货单号"}},
		}),
	}
}

// List 查询 ERP 销售退货列表。
func (s *SaleReturnServiceImpl) List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
