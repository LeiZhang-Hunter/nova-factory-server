package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductUnitDaoImpl 提供 ERP 产品单位数据访问能力。
type ProductUnitDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.ProductUnit, mastermodels.ProductUnitUpsert, mastermodels.ProductUnitQuery]
}

// NewProductUnitDao 创建 ERP 产品单位 DAO。
func NewProductUnitDao(db *gorm.DB) masterdao.IProductUnitDao {
	return &ProductUnitDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.ProductUnit, mastermodels.ProductUnitUpsert, mastermodels.ProductUnitQuery](db, erpcrud.EntityConfig{
			Table:        "erp_product_unit",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品单位列表。
func (d *ProductUnitDaoImpl) List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductUnitListData{Rows: result.Rows, Total: result.Total}, nil
}
