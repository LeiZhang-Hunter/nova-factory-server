package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// ProductCategoryServiceImpl 提供 ERP 产品分类业务实现。
type ProductCategoryServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.ProductCategory, mastermodels.ProductCategoryUpsert, mastermodels.ProductCategoryQuery]
}

// NewProductCategoryService 创建 ERP 产品分类服务。
func NewProductCategoryService(dao masterdao.IProductCategoryDao) masterservice.IProductCategoryService {
	return &ProductCategoryServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.ProductCategory, mastermodels.ProductCategoryUpsert, mastermodels.ProductCategoryQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_product_category",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品分类列表。
func (s *ProductCategoryServiceImpl) List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductCategoryListData{Rows: result.Rows, Total: result.Total}, nil
}
