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

// PurchaseReturnItemDaoImpl 提供数据访问能力。
type PurchaseReturnItemDaoImpl struct {
	db *gorm.DB
}

// NewPurchaseReturnItemDao 创建 DAO。
func NewPurchaseReturnItemDao(db *gorm.DB) purchasedao.IPurchaseReturnItemDao {
	return &PurchaseReturnItemDaoImpl{db: db}
}

func (d *PurchaseReturnItemDaoImpl) Create(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error) {
	model := new(purchasemodels.PurchaseReturnItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_purchase_return_items").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *PurchaseReturnItemDaoImpl) Update(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(purchasemodels.PurchaseReturnItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_purchase_return_items").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *PurchaseReturnItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_purchase_return_items").Where("id IN ?", ids)
	if erpbiz.HasField(new(purchasemodels.PurchaseReturnItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *PurchaseReturnItemDaoImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturnItem, error) {
	item := new(purchasemodels.PurchaseReturnItem)
	db := d.db.WithContext(c).Table("erp_purchase_return_items").Where("id = ?", id)
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

func (d *PurchaseReturnItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseReturnItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(purchasemodels.PurchaseReturnItem)
	db := d.db.WithContext(c).Table("erp_purchase_return_items").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *PurchaseReturnItemDaoImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*erpbiz.PageResult[purchasemodels.PurchaseReturnItem], error) {
	if req == nil {
		req = new(purchasemodels.PurchaseReturnItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_purchase_return_items")
	if erpbiz.HasField(new(purchasemodels.PurchaseReturnItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]purchasemodels.PurchaseReturnItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*purchasemodels.PurchaseReturnItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[purchasemodels.PurchaseReturnItem]{Rows: result, Total: total}, nil
}

func (d *PurchaseReturnItemDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*purchasemodels.PurchaseReturnItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
