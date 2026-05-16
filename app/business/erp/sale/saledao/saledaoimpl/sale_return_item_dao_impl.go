package saledaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaleReturnItemDaoImpl 提供 ERP 销售退货项数据访问能力。
type SaleReturnItemDaoImpl struct {
	*erpcrud.CRUDDao[salemodels.SaleReturnItem, salemodels.SaleReturnItemUpsert, salemodels.SaleReturnItemQuery]
}

// NewSaleReturnItemDao 创建 ERP 销售退货项 DAO。
func NewSaleReturnItemDao(db *gorm.DB) saledao.ISaleReturnItemDao {
	return &SaleReturnItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[salemodels.SaleReturnItem, salemodels.SaleReturnItemUpsert, salemodels.SaleReturnItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_sale_return_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 销售退货项列表。
func (d *SaleReturnItemDaoImpl) List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
