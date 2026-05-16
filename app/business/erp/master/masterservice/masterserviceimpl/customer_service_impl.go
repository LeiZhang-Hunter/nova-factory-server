package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// CustomerServiceImpl 提供 ERP 客户业务实现。
type CustomerServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.Customer, mastermodels.CustomerUpsert, mastermodels.CustomerQuery]
}

// NewCustomerService 创建 ERP 客户服务。
func NewCustomerService(dao masterdao.ICustomerDao) masterservice.ICustomerService {
	return &CustomerServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.Customer, mastermodels.CustomerUpsert, mastermodels.CustomerQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_customer",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "Code", Column: "code", Label: "客户编码，对接 erp_order.b_type_code"}},
		}),
	}
}

// List 查询 ERP 客户列表。
func (s *CustomerServiceImpl) List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.CustomerListData{Rows: result.Rows, Total: result.Total}, nil
}
