package masterdaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductUnitDaoImpl 提供数据访问能力。
type ProductUnitDaoImpl struct {
	db *gorm.DB
}

// NewProductUnitDao 创建 DAO。
func NewProductUnitDao(db *gorm.DB) masterdao.IProductUnitDao {
	return &ProductUnitDaoImpl{db: db}
}

func (d *ProductUnitDaoImpl) Create(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error) {
	model := new(mastermodels.ProductUnit)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_product_unit").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *ProductUnitDaoImpl) Update(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(mastermodels.ProductUnit)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_product_unit").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *ProductUnitDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_product_unit").Where("id IN ?", ids)
	if erpbiz.HasField(new(mastermodels.ProductUnit), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *ProductUnitDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.ProductUnit, error) {
	item := new(mastermodels.ProductUnit)
	db := d.db.WithContext(c).Table("erp_product_unit").Where("id = ?", id)
	if erpbiz.HasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *ProductUnitDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductUnit, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.ProductUnit)
	db := d.db.WithContext(c).Table("erp_product_unit").Where(fmt.Sprintf("%s = ?", column), value)
	if erpbiz.HasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *ProductUnitDaoImpl) ListPage(c *gin.Context, req *mastermodels.ProductUnitQuery) (*erpbiz.PageResult[mastermodels.ProductUnit], error) {
	if req == nil {
		req = new(mastermodels.ProductUnitQuery)
	}
	db := d.db.WithContext(c).Table("erp_product_unit")
	if erpbiz.HasField(new(mastermodels.ProductUnit), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.ProductUnit, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.ProductUnit, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[mastermodels.ProductUnit]{Rows: result, Total: total}, nil
}

func (d *ProductUnitDaoImpl) List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductUnitListData{Rows: result.Rows, Total: result.Total}, nil
}
