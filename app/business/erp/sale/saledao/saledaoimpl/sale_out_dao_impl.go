package saledaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaleOutDaoImpl 提供 ERP 销售出库数据访问能力。
type SaleOutDaoImpl struct {
	*erpcrud.CRUDDao[salemodels.SaleOut, salemodels.SaleOutUpsert, salemodels.SaleOutQuery]
}

// NewSaleOutDao 创建 ERP 销售出库 DAO。
func NewSaleOutDao(db *gorm.DB) saledao.ISaleOutDao {
	return &SaleOutDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[salemodels.SaleOut, salemodels.SaleOutUpsert, salemodels.SaleOutQuery](db, erpcrud.EntityConfig{
			Table:        "erp_sale_out",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "销售出库单号"}},
		}),
	}
}

// List 查询 ERP 销售出库列表。
func (d *SaleOutDaoImpl) List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutListData{Rows: result.Rows, Total: result.Total}, nil
}
