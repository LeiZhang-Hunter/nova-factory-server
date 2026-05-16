package saledaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaleReturnDaoImpl 提供 ERP 销售退货数据访问能力。
type SaleReturnDaoImpl struct {
	*erpcrud.CRUDDao[salemodels.SaleReturn, salemodels.SaleReturnUpsert, salemodels.SaleReturnQuery]
}

// NewSaleReturnDao 创建 ERP 销售退货 DAO。
func NewSaleReturnDao(db *gorm.DB) saledao.ISaleReturnDao {
	return &SaleReturnDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[salemodels.SaleReturn, salemodels.SaleReturnUpsert, salemodels.SaleReturnQuery](db, erpcrud.EntityConfig{
			Table:        "erp_sale_return",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "销售退货单号"}},
		}),
	}
}

// List 查询 ERP 销售退货列表。
func (d *SaleReturnDaoImpl) List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
