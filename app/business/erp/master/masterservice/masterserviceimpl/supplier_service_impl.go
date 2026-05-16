package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// SupplierServiceImpl 提供 ERP 供应商业务实现。
type SupplierServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.Supplier, mastermodels.SupplierUpsert, mastermodels.SupplierQuery]
}

// NewSupplierService 创建 ERP 供应商服务。
func NewSupplierService(dao masterdao.ISupplierDao) masterservice.ISupplierService {
	return &SupplierServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.Supplier, mastermodels.SupplierUpsert, mastermodels.SupplierQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_supplier",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "Code", Column: "code", Label: "供应商编码"}},
		}),
	}
}

// List 查询 ERP 供应商列表。
func (s *SupplierServiceImpl) List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.SupplierListData{Rows: result.Rows, Total: result.Total}, nil
}
