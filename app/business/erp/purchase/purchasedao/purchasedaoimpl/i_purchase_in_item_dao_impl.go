package purchasedaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseInItemDaoImpl 提供数据访问能力。
type PurchaseInItemDaoImpl struct {
	db *gorm.DB
}

// NewPurchaseInItemDao 创建 DAO。
func NewPurchaseInItemDao(db *gorm.DB) purchasedao.IPurchaseInItemDao {
	return &PurchaseInItemDaoImpl{db: db}
}

func (d *PurchaseInItemDaoImpl) Create(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error) {
	model := new(purchasemodels.PurchaseInItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_purchase_in_items").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *PurchaseInItemDaoImpl) Update(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(purchasemodels.PurchaseInItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_purchase_in_items").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *PurchaseInItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_purchase_in_items").Where("id IN ?", ids)
	if erpbiz.HasField(new(purchasemodels.PurchaseInItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *PurchaseInItemDaoImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseInItem, error) {
	item := new(purchasemodels.PurchaseInItem)
	db := d.db.WithContext(c).Table("erp_purchase_in_items").Where("id = ?", id)
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

func (d *PurchaseInItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseInItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(purchasemodels.PurchaseInItem)
	db := d.db.WithContext(c).Table("erp_purchase_in_items").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *PurchaseInItemDaoImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*erpbiz.PageResult[purchasemodels.PurchaseInItem], error) {
	if req == nil {
		req = new(purchasemodels.PurchaseInItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_purchase_in_items")
	if erpbiz.HasField(new(purchasemodels.PurchaseInItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]purchasemodels.PurchaseInItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*purchasemodels.PurchaseInItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[purchasemodels.PurchaseInItem]{Rows: result, Total: total}, nil
}

func (d *PurchaseInItemDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseInItemListData{Rows: result.Rows, Total: result.Total}, nil
}
