package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductCategoryDaoImpl 提供 ERP 产品分类数据访问能力。
type ProductCategoryDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.ProductCategory, mastermodels.ProductCategoryUpsert, mastermodels.ProductCategoryQuery]
}

// NewProductCategoryDao 创建 ERP 产品分类 DAO。
func NewProductCategoryDao(db *gorm.DB) masterdao.IProductCategoryDao {
	return &ProductCategoryDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.ProductCategory, mastermodels.ProductCategoryUpsert, mastermodels.ProductCategoryQuery](db, erpcrud.EntityConfig{
			Table:        "erp_product_category",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品分类列表。
func (d *ProductCategoryDaoImpl) List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductCategoryListData{Rows: result.Rows, Total: result.Total}, nil
}
