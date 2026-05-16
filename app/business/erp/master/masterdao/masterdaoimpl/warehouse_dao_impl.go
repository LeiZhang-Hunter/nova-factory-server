package masterdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WarehouseDaoImpl 提供 ERP 仓库数据访问能力。
type WarehouseDaoImpl struct {
	*erpcrud.CRUDDao[mastermodels.Warehouse, mastermodels.WarehouseUpsert, mastermodels.WarehouseQuery]
}

// NewWarehouseDao 创建 ERP 仓库 DAO。
func NewWarehouseDao(db *gorm.DB) masterdao.IWarehouseDao {
	return &WarehouseDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[mastermodels.Warehouse, mastermodels.WarehouseUpsert, mastermodels.WarehouseQuery](db, erpcrud.EntityConfig{
			Table:        "erp_warehouse",
			OrderBy:      "sort ASC, id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 仓库列表。
func (d *WarehouseDaoImpl) List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.WarehouseListData{Rows: result.Rows, Total: result.Total}, nil
}
