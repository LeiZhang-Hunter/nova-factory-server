package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// AccountServiceImpl 提供 ERP 结算账户业务实现。
type AccountServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.Account, mastermodels.AccountUpsert, mastermodels.AccountQuery]
}

// NewAccountService 创建 ERP 结算账户服务。
func NewAccountService(dao masterdao.IAccountDao) masterservice.IAccountService {
	return &AccountServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.Account, mastermodels.AccountUpsert, mastermodels.AccountQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_account",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "账户编码，对接 erp_order_account.finance_code"}},
		}),
	}
}

// List 查询 ERP 结算账户列表。
func (s *AccountServiceImpl) List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.AccountListData{Rows: result.Rows, Total: result.Total}, nil
}
