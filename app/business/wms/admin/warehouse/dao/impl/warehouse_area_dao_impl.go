package impl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/wms/admin/warehouse/dao"
	"nova-factory-server/app/business/wms/admin/warehouse/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WarehouseAreaDaoImpl 提供 WMS 库区的数据访问能力。
type WarehouseAreaDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewWarehouseAreaDao 创建 WMS 库区 DAO。
func NewWarehouseAreaDao(db *gorm.DB) dao.IWarehouseAreaDao {
	return &WarehouseAreaDaoImpl{
		db:    db,
		table: "wms_warehouse_area",
	}
}

// Set 新增或修改 WMS 库区。
func (d *WarehouseAreaDaoImpl) Set(c *gin.Context, req *models.WarehouseAreaSet) (*models.WarehouseArea, error) {
	if req.ID > 0 {
		if err := d.update(c, req); err != nil {
			return nil, err
		}
		return d.GetByID(c, req.ID)
	}
	return d.create(c, req)
}

// DeleteByIDs 软删除 WMS 库区。
func (d *WarehouseAreaDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return d.db.WithContext(c).Table(d.table).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":            commonStatus.DELETE,
			"update_by":        baizeContext.GetUserId(c),
			"update_time":      now,
			"last_update_time": now,
		}).Error
}

// GetByID 根据主键查询 WMS 库区。
func (d *WarehouseAreaDaoImpl) GetByID(c *gin.Context, id int64) (*models.WarehouseArea, error) {
	item := new(models.WarehouseArea)
	if err := d.baseQuery(c).Where("id = ?", id).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

// GetByIDs 根据主键批量查询 WMS 库区。
func (d *WarehouseAreaDaoImpl) GetByIDs(c *gin.Context, ids []int64) ([]*models.WarehouseArea, error) {
	if len(ids) == 0 {
		return []*models.WarehouseArea{}, nil
	}
	rows := make([]*models.WarehouseArea, 0, len(ids))
	if err := d.baseQuery(c).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetByName 按仓库、父级和名称查询 WMS 库区。
func (d *WarehouseAreaDaoImpl) GetByName(c *gin.Context, warehouseID, parentID int64, areaName string) (*models.WarehouseArea, error) {
	item := new(models.WarehouseArea)
	if err := d.baseQuery(c).
		Where("warehouse_id = ?", warehouseID).
		Where("parent_id = ?", parentID).
		Where("area_name = ?", strings.TrimSpace(areaName)).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

// List 分页查询 WMS 库区列表。
func (d *WarehouseAreaDaoImpl) List(c *gin.Context, req *models.WarehouseAreaQuery) (*models.WarehouseAreaListData, error) {
	if req == nil {
		req = new(models.WarehouseAreaQuery)
	}
	db := d.baseQuery(c)
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.AreaName != "" {
		db = db.Where("area_name LIKE ?", "%"+strings.TrimSpace(req.AreaName)+"%")
	}
	if req.ParentID > 0 {
		db = db.Where("parent_id = ?", req.ParentID)
	}
	if req.WarehouseAreaProperty != nil {
		db = db.Where("warehouse_area_property = ?", *req.WarehouseAreaProperty)
	}
	if req.IsValid != nil {
		db = db.Where("is_valid = ?", *req.IsValid)
	}
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.WarehouseArea, 0)
	if err := db.Order("warehouse_id ASC").Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.WarehouseAreaListData{Rows: rows, Total: total}, nil
}

func (d *WarehouseAreaDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return d.db.WithContext(c).Table(d.table).Where("state = ?", commonStatus.NORMAL)
}

func (d *WarehouseAreaDaoImpl) create(c *gin.Context, req *models.WarehouseAreaSet) (*models.WarehouseArea, error) {
	now := time.Now()
	item := &models.WarehouseArea{
		ID:                    snowflake.GenID(),
		WarehouseID:           req.WarehouseID,
		AreaName:              req.AreaName,
		ParentID:              req.ParentID,
		WarehouseAreaProperty: req.WarehouseAreaProperty,
		LastUpdateTime:        &now,
		IsValid:               req.IsValid,
		DeptID:                baizeContext.GetDeptId(c),
		State:                 commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table(d.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (d *WarehouseAreaDaoImpl) update(c *gin.Context, req *models.WarehouseAreaSet) error {
	existing, err := d.GetByID(c, req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("库区不存在")
	}
	now := time.Now()
	return d.db.WithContext(c).Table(d.table).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"warehouse_id":            req.WarehouseID,
			"area_name":               req.AreaName,
			"parent_id":               req.ParentID,
			"warehouse_area_property": req.WarehouseAreaProperty,
			"is_valid":                req.IsValid,
			"update_by":               baizeContext.GetUserId(c),
			"update_time":             now,
			"last_update_time":        now,
		}).Error
}
