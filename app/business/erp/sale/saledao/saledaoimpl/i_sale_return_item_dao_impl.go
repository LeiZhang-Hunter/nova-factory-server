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

// SaleReturnItemDaoImpl 提供数据访问能力。
type SaleReturnItemDaoImpl struct {
	db *gorm.DB
}

// NewSaleReturnItemDao 创建 DAO。
func NewSaleReturnItemDao(db *gorm.DB) saledao.ISaleReturnItemDao {
	return &SaleReturnItemDaoImpl{db: db}
}

func (d *SaleReturnItemDaoImpl) Create(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error) {
	model := new(salemodels.SaleReturnItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_sale_return_items").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *SaleReturnItemDaoImpl) Update(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(salemodels.SaleReturnItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_sale_return_items").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *SaleReturnItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_sale_return_items").Where("id IN ?", ids)
	if erpbiz.HasField(new(salemodels.SaleReturnItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *SaleReturnItemDaoImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleReturnItem, error) {
	item := new(salemodels.SaleReturnItem)
	db := d.db.WithContext(c).Table("erp_sale_return_items").Where("id = ?", id)
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

func (d *SaleReturnItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleReturnItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(salemodels.SaleReturnItem)
	db := d.db.WithContext(c).Table("erp_sale_return_items").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *SaleReturnItemDaoImpl) ListPage(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*erpbiz.PageResult[salemodels.SaleReturnItem], error) {
	if req == nil {
		req = new(salemodels.SaleReturnItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_sale_return_items")
	if erpbiz.HasField(new(salemodels.SaleReturnItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]salemodels.SaleReturnItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*salemodels.SaleReturnItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[salemodels.SaleReturnItem]{Rows: result, Total: total}, nil
}

func (d *SaleReturnItemDaoImpl) List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
