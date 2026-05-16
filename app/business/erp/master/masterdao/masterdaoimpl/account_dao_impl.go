package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountDaoImpl 提供 ERP 结算账户数据访问能力。
type AccountDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.Account, mastermodels.AccountUpsert, mastermodels.AccountQuery]
}

// NewAccountDao 创建 ERP 结算账户 DAO。
func NewAccountDao(db *gorm.DB) masterdao.IAccountDao {
	return &AccountDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.Account, mastermodels.AccountUpsert, mastermodels.AccountQuery](db, erpcrud.EntityConfig{
			Table:        "erp_account",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "账户编码，对接 erp_order_account.finance_code"}},
		}),
	}
}

// List 查询 ERP 结算账户列表。
func (d *AccountDaoImpl) List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.AccountListData{Rows: result.Rows, Total: result.Total}, nil
}
