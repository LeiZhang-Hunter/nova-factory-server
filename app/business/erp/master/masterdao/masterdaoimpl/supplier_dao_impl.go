package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SupplierDaoImpl 提供 ERP 供应商数据访问能力。
type SupplierDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.Supplier, mastermodels.SupplierUpsert, mastermodels.SupplierQuery]
}

// NewSupplierDao 创建 ERP 供应商 DAO。
func NewSupplierDao(db *gorm.DB) masterdao.ISupplierDao {
	return &SupplierDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.Supplier, mastermodels.SupplierUpsert, mastermodels.SupplierQuery](db, erpcrud.EntityConfig{
			Table:        "erp_supplier",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "Code", Column: "code", Label: "供应商编码"}},
		}),
	}
}

// List 查询 ERP 供应商列表。
func (d *SupplierDaoImpl) List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.SupplierListData{Rows: result.Rows, Total: result.Total}, nil
}
