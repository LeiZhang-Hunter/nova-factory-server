package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// ProductServiceImpl 提供 ERP 产品业务实现。
type ProductServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.Product, mastermodels.ProductUpsert, mastermodels.ProductQuery]
}

// NewProductService 创建 ERP 产品服务。
func NewProductService(dao masterdao.IProductDao) masterservice.IProductService {
	return &ProductServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.Product, mastermodels.ProductUpsert, mastermodels.ProductQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_product",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品列表。
func (s *ProductServiceImpl) List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductListData{Rows: result.Rows, Total: result.Total}, nil
}
