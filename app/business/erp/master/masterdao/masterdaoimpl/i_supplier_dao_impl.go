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

// SupplierDaoImpl 提供数据访问能力。
type SupplierDaoImpl struct {
	db *gorm.DB
}

// NewSupplierDao 创建 DAO。
func NewSupplierDao(db *gorm.DB) masterdao.ISupplierDao {
	return &SupplierDaoImpl{db: db}
}

func (d *SupplierDaoImpl) Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	model := new(mastermodels.Supplier)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_supplier").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *SupplierDaoImpl) Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(mastermodels.Supplier)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_supplier").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *SupplierDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_supplier").Where("id IN ?", ids)
	if erpbiz.HasField(new(mastermodels.Supplier), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *SupplierDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error) {
	item := new(mastermodels.Supplier)
	db := d.db.WithContext(c).Table("erp_supplier").Where("id = ?", id)
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

func (d *SupplierDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Supplier, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Supplier)
	db := d.db.WithContext(c).Table("erp_supplier").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *SupplierDaoImpl) ListPage(c *gin.Context, req *mastermodels.SupplierQuery) (*erpbiz.PageResult[mastermodels.Supplier], error) {
	if req == nil {
		req = new(mastermodels.SupplierQuery)
	}
	db := d.db.WithContext(c).Table("erp_supplier")
	if erpbiz.HasField(new(mastermodels.Supplier), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Supplier, 0)
	orderBy := strings.TrimSpace("sort ASC, id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Supplier, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[mastermodels.Supplier]{Rows: result, Total: total}, nil
}

func (d *SupplierDaoImpl) List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.SupplierListData{Rows: result.Rows, Total: result.Total}, nil
}
