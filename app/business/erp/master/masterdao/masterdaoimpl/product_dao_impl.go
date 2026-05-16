package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductDaoImpl 提供 ERP 产品数据访问能力。
type ProductDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.Product, mastermodels.ProductUpsert, mastermodels.ProductQuery]
}

// NewProductDao 创建 ERP 产品 DAO。
func NewProductDao(db *gorm.DB) masterdao.IProductDao {
	return &ProductDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.Product, mastermodels.ProductUpsert, mastermodels.ProductQuery](db, erpcrud.EntityConfig{
			Table:        "erp_product",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品列表。
func (d *ProductDaoImpl) List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductListData{Rows: result.Rows, Total: result.Total}, nil
}
