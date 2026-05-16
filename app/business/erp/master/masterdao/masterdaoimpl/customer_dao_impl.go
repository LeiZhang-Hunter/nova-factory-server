package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerDaoImpl 提供 ERP 客户数据访问能力。
type CustomerDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.Customer, mastermodels.CustomerUpsert, mastermodels.CustomerQuery]
}

// NewCustomerDao 创建 ERP 客户 DAO。
func NewCustomerDao(db *gorm.DB) masterdao.ICustomerDao {
	return &CustomerDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.Customer, mastermodels.CustomerUpsert, mastermodels.CustomerQuery](db, erpcrud.EntityConfig{
			Table:        "erp_customer",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "Code", Column: "code", Label: "客户编码，对接 erp_order.b_type_code"}},
		}),
	}
}

// List 查询 ERP 客户列表。
func (d *CustomerDaoImpl) List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.CustomerListData{Rows: result.Rows, Total: result.Total}, nil
}
