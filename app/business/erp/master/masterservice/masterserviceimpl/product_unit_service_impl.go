package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// ProductUnitServiceImpl 提供 ERP 产品单位业务实现。
type ProductUnitServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.ProductUnit, mastermodels.ProductUnitUpsert, mastermodels.ProductUnitQuery]
}

// NewProductUnitService 创建 ERP 产品单位服务。
func NewProductUnitService(dao masterdao.IProductUnitDao) masterservice.IProductUnitService {
	return &ProductUnitServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.ProductUnit, mastermodels.ProductUnitUpsert, mastermodels.ProductUnitQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_product_unit",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品单位列表。
func (s *ProductUnitServiceImpl) List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductUnitListData{Rows: result.Rows, Total: result.Total}, nil
}
