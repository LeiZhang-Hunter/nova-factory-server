package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// WarehouseServiceImpl 提供 ERP 仓库业务实现。
type WarehouseServiceImpl struct {
	*erpcrud.CRUDService[mastermodels.Warehouse, mastermodels.WarehouseUpsert, mastermodels.WarehouseQuery]
}

// NewWarehouseService 创建 ERP 仓库服务。
func NewWarehouseService(dao masterdao.IWarehouseDao) masterservice.IWarehouseService {
	return &WarehouseServiceImpl{
		CRUDService: erpcrud.NewCRUDService[mastermodels.Warehouse, mastermodels.WarehouseUpsert, mastermodels.WarehouseQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_warehouse",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 仓库列表。
func (s *WarehouseServiceImpl) List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.WarehouseListData{Rows: result.Rows, Total: result.Total}, nil
}
