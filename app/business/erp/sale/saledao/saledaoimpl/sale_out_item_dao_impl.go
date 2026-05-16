package saledaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaleOutItemDaoImpl 提供 ERP 销售出库项数据访问能力。
type SaleOutItemDaoImpl struct {
	*erpcrud.CRUDDao[salemodels.SaleOutItem, salemodels.SaleOutItemUpsert, salemodels.SaleOutItemQuery]
}

// NewSaleOutItemDao 创建 ERP 销售出库项 DAO。
func NewSaleOutItemDao(db *gorm.DB) saledao.ISaleOutItemDao {
	return &SaleOutItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[salemodels.SaleOutItem, salemodels.SaleOutItemUpsert, salemodels.SaleOutItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_sale_out_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 销售出库项列表。
func (d *SaleOutItemDaoImpl) List(c *gin.Context, req *salemodels.SaleOutItemQuery) (*salemodels.SaleOutItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutItemListData{Rows: result.Rows, Total: result.Total}, nil
}
