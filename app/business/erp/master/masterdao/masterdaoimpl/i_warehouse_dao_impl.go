package masterdaoimpl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseDaoImpl struct {
	db *gorm.DB
}

func NewWarehouseDao(db *gorm.DB) masterdao.IWarehouseDao {
	return &WarehouseDaoImpl{db: db}
}

func (d *WarehouseDaoImpl) Create(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.WarehouseUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_warehouse").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *WarehouseDaoImpl) Update(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.Principal != "" {
		updates["principal"] = req.Principal
	}
	if req.WarehousePrice != 0 {
		updates["warehouse_price"] = req.WarehousePrice
	}
	if req.TruckagePrice != 0 {
		updates["truckage_price"] = req.TruckagePrice
	}
	updates["status"] = req.Status
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_warehouse").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *WarehouseDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_warehouse").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *WarehouseDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Warehouse, error) {
	item := new(mastermodels.Warehouse)
	if err := d.db.WithContext(c).Table("erp_warehouse").
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *WarehouseDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Warehouse, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Warehouse)
	if err := d.db.WithContext(c).Table("erp_warehouse").
		Where(fmt.Sprintf("%s = ?", column), value).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *WarehouseDaoImpl) ListPage(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error) {
	if req == nil {
		req = new(mastermodels.WarehouseQuery)
	}
	db := d.db.WithContext(c).Table("erp_warehouse").Where("state = ?", commonStatus.NORMAL)
	db = applyWarehouseFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Warehouse, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Warehouse, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.WarehouseListData{Rows: result, Total: total}, nil
}

func (d *WarehouseDaoImpl) List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.WarehouseListData{Rows: result.Rows, Total: result.Total}, nil
}
