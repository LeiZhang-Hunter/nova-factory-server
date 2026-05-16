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

// ProductDaoImpl 提供数据访问能力。
type ProductDaoImpl struct {
	db *gorm.DB
}

// NewProductDao 创建 DAO。
func NewProductDao(db *gorm.DB) masterdao.IProductDao {
	return &ProductDaoImpl{db: db}
}

func (d *ProductDaoImpl) Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	model := new(mastermodels.Product)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_product").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *ProductDaoImpl) Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(mastermodels.Product)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_product").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *ProductDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_product").Where("id IN ?", ids)
	if erpbiz.HasField(new(mastermodels.Product), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *ProductDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Product, error) {
	item := new(mastermodels.Product)
	db := d.db.WithContext(c).Table("erp_product").Where("id = ?", id)
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

func (d *ProductDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Product, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Product)
	db := d.db.WithContext(c).Table("erp_product").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *ProductDaoImpl) ListPage(c *gin.Context, req *mastermodels.ProductQuery) (*erpbiz.PageResult[mastermodels.Product], error) {
	if req == nil {
		req = new(mastermodels.ProductQuery)
	}
	db := d.db.WithContext(c).Table("erp_product")
	if erpbiz.HasField(new(mastermodels.Product), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Product, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Product, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[mastermodels.Product]{Rows: result, Total: total}, nil
}

func (d *ProductDaoImpl) List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductListData{Rows: result.Rows, Total: result.Total}, nil
}
