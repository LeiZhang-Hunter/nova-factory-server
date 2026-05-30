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

// WarehouseLocationDaoImpl 提供 WMS 库位的数据访问能力。
type WarehouseLocationDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewWarehouseLocationDao 创建 WMS 库位 DAO。
func NewWarehouseLocationDao(db *gorm.DB) dao.IWarehouseLocationDao {
	return &WarehouseLocationDaoImpl{
		db:    db,
		table: "wms_warehouse_location",
	}
}

// Set 新增或修改 WMS 库位。
func (d *WarehouseLocationDaoImpl) Set(c *gin.Context, req *models.WarehouseLocationSet) (*models.WarehouseLocation, error) {
	if req.ID > 0 {
		if err := d.update(c, req); err != nil {
			return nil, err
		}
		return d.GetByID(c, req.ID)
	}
	return d.create(c, req)
}

// DeleteByIDs 软删除 WMS 库位。
func (d *WarehouseLocationDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
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

// GetByID 根据主键查询 WMS 库位。
func (d *WarehouseLocationDaoImpl) GetByID(c *gin.Context, id int64) (*models.WarehouseLocation, error) {
	item := new(models.WarehouseLocation)
	if err := d.baseQuery(c).Where("id = ?", id).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

// GetByName 按库区和名称查询 WMS 库位。
func (d *WarehouseLocationDaoImpl) GetByName(c *gin.Context, warehouseAreaID int64, locationName string) (*models.WarehouseLocation, error) {
	item := new(models.WarehouseLocation)
	if err := d.baseQuery(c).
		Where("warehouse_area_id = ?", warehouseAreaID).
		Where("location_name = ?", strings.TrimSpace(locationName)).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

// GetByTagNumber 按标签号查询 WMS 库位。
func (d *WarehouseLocationDaoImpl) GetByTagNumber(c *gin.Context, tagNumber string) (*models.WarehouseLocation, error) {
	if strings.TrimSpace(tagNumber) == "" {
		return nil, nil
	}
	item := new(models.WarehouseLocation)
	if err := d.baseQuery(c).
		Where("tag_number = ?", strings.TrimSpace(tagNumber)).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

// List 分页查询 WMS 库位列表。
func (d *WarehouseLocationDaoImpl) List(c *gin.Context, req *models.WarehouseLocationQuery) (*models.WarehouseLocationListData, error) {
	if req == nil {
		req = new(models.WarehouseLocationQuery)
	}
	db := d.baseQuery(c)
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.WarehouseAreaID > 0 {
		db = db.Where("warehouse_area_id = ?", req.WarehouseAreaID)
	}
	if req.LocationName != "" {
		db = db.Where("location_name LIKE ?", "%"+strings.TrimSpace(req.LocationName)+"%")
	}
	if req.RoadwayNumber != "" {
		db = db.Where("roadway_number = ?", strings.TrimSpace(req.RoadwayNumber))
	}
	if req.ShelfNumber != "" {
		db = db.Where("shelf_number = ?", strings.TrimSpace(req.ShelfNumber))
	}
	if req.LayerNumber != "" {
		db = db.Where("layer_number = ?", strings.TrimSpace(req.LayerNumber))
	}
	if req.TagNumber != "" {
		db = db.Where("tag_number = ?", strings.TrimSpace(req.TagNumber))
	}
	if req.IsValid != nil {
		db = db.Where("is_valid = ?", *req.IsValid)
	}
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.WarehouseLocation, 0)
	if err := db.Order("warehouse_id ASC").Order("warehouse_area_id ASC").Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.WarehouseLocationListData{Rows: rows, Total: total}, nil
}

func (d *WarehouseLocationDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return d.db.WithContext(c).Table(d.table).Where("state = ?", commonStatus.NORMAL)
}

func (d *WarehouseLocationDaoImpl) create(c *gin.Context, req *models.WarehouseLocationSet) (*models.WarehouseLocation, error) {
	now := time.Now()
	item := &models.WarehouseLocation{
		ID:                    snowflake.GenID(),
		WarehouseID:           req.WarehouseID,
		WarehouseName:         req.WarehouseName,
		WarehouseAreaID:       req.WarehouseAreaID,
		WarehouseAreaName:     req.WarehouseAreaName,
		WarehouseAreaProperty: req.WarehouseAreaProperty,
		LocationName:          req.LocationName,
		LocationLength:        req.LocationLength,
		LocationWidth:         req.LocationWidth,
		LocationHeight:        req.LocationHeight,
		LocationVolume:        req.LocationVolume,
		LocationLoad:          req.LocationLoad,
		RoadwayNumber:         req.RoadwayNumber,
		ShelfNumber:           req.ShelfNumber,
		LayerNumber:           req.LayerNumber,
		TagNumber:             req.TagNumber,
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

func (d *WarehouseLocationDaoImpl) update(c *gin.Context, req *models.WarehouseLocationSet) error {
	existing, err := d.GetByID(c, req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("库位不存在")
	}
	now := time.Now()
	return d.db.WithContext(c).Table(d.table).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"warehouse_id":            req.WarehouseID,
			"warehouse_name":          req.WarehouseName,
			"warehouse_area_id":       req.WarehouseAreaID,
			"warehouse_area_name":     req.WarehouseAreaName,
			"warehouse_area_property": req.WarehouseAreaProperty,
			"location_name":           req.LocationName,
			"location_length":         req.LocationLength,
			"location_width":          req.LocationWidth,
			"location_height":         req.LocationHeight,
			"location_volume":         req.LocationVolume,
			"location_load":           req.LocationLoad,
			"roadway_number":          req.RoadwayNumber,
			"shelf_number":            req.ShelfNumber,
			"layer_number":            req.LayerNumber,
			"tag_number":              req.TagNumber,
			"is_valid":                req.IsValid,
			"update_by":               baizeContext.GetUserId(c),
			"update_time":             now,
			"last_update_time":        now,
		}).Error
}
