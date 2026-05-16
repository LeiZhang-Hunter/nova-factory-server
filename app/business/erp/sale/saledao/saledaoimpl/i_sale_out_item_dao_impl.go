package saledaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaleOutItemDaoImpl 提供数据访问能力。
type SaleOutItemDaoImpl struct {
	db *gorm.DB
}

// NewSaleOutItemDao 创建 DAO。
func NewSaleOutItemDao(db *gorm.DB) saledao.ISaleOutItemDao {
	return &SaleOutItemDaoImpl{db: db}
}

func (d *SaleOutItemDaoImpl) Create(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error) {
	model := new(salemodels.SaleOutItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_sale_out_items").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *SaleOutItemDaoImpl) Update(c *gin.Context, req *salemodels.SaleOutItemUpsert) (*salemodels.SaleOutItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(salemodels.SaleOutItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_sale_out_items").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *SaleOutItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_sale_out_items").Where("id IN ?", ids)
	if erpbiz.HasField(new(salemodels.SaleOutItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *SaleOutItemDaoImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleOutItem, error) {
	item := new(salemodels.SaleOutItem)
	db := d.db.WithContext(c).Table("erp_sale_out_items").Where("id = ?", id)
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

func (d *SaleOutItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleOutItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(salemodels.SaleOutItem)
	db := d.db.WithContext(c).Table("erp_sale_out_items").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *SaleOutItemDaoImpl) ListPage(c *gin.Context, req *salemodels.SaleOutItemQuery) (*erpbiz.PageResult[salemodels.SaleOutItem], error) {
	if req == nil {
		req = new(salemodels.SaleOutItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_sale_out_items")
	if erpbiz.HasField(new(salemodels.SaleOutItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]salemodels.SaleOutItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*salemodels.SaleOutItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[salemodels.SaleOutItem]{Rows: result, Total: total}, nil
}

func (d *SaleOutItemDaoImpl) List(c *gin.Context, req *salemodels.SaleOutItemQuery) (*salemodels.SaleOutItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutItemListData{Rows: result.Rows, Total: result.Total}, nil
}
