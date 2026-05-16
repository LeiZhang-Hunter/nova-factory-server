package saleserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"

	"github.com/gin-gonic/gin"
)

// SaleOutServiceImpl 提供 ERP 销售出库业务实现。
type SaleOutServiceImpl struct {
	*erpcrud.CRUDService[salemodels.SaleOut, salemodels.SaleOutUpsert, salemodels.SaleOutQuery]
}

// NewSaleOutService 创建 ERP 销售出库服务。
func NewSaleOutService(dao saledao.ISaleOutDao) saleservice.ISaleOutService {
	return &SaleOutServiceImpl{
		CRUDService: erpcrud.NewCRUDService[salemodels.SaleOut, salemodels.SaleOutUpsert, salemodels.SaleOutQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_sale_out",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "销售出库单号"}},
		}),
	}
}

// List 查询 ERP 销售出库列表。
func (s *SaleOutServiceImpl) List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutListData{Rows: result.Rows, Total: result.Total}, nil
}
