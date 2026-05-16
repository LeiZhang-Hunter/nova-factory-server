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

// PurchaseOrderDaoImpl 提供数据访问能力。
type PurchaseOrderDaoImpl struct {
	db *gorm.DB
}

// NewPurchaseOrderDao 创建 DAO。
func NewPurchaseOrderDao(db *gorm.DB) purchasedao.IPurchaseOrderDao {
	return &PurchaseOrderDaoImpl{db: db}
}

func (d *PurchaseOrderDaoImpl) Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error) {
	model := new(purchasemodels.PurchaseOrder)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_purchase_order").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *PurchaseOrderDaoImpl) Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(purchasemodels.PurchaseOrder)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_purchase_order").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *PurchaseOrderDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_purchase_order").Where("id IN ?", ids)
	if erpbiz.HasField(new(purchasemodels.PurchaseOrder), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *PurchaseOrderDaoImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error) {
	item := new(purchasemodels.PurchaseOrder)
	db := d.db.WithContext(c).Table("erp_purchase_order").Where("id = ?", id)
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

func (d *PurchaseOrderDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseOrder, error) {
	if column == "" {
		return nil, nil
	}
	item := new(purchasemodels.PurchaseOrder)
	db := d.db.WithContext(c).Table("erp_purchase_order").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *PurchaseOrderDaoImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*erpbiz.PageResult[purchasemodels.PurchaseOrder], error) {
	if req == nil {
		req = new(purchasemodels.PurchaseOrderQuery)
	}
	db := d.db.WithContext(c).Table("erp_purchase_order")
	if erpbiz.HasField(new(purchasemodels.PurchaseOrder), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]purchasemodels.PurchaseOrder, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*purchasemodels.PurchaseOrder, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[purchasemodels.PurchaseOrder]{Rows: result, Total: total}, nil
}

func (d *PurchaseOrderDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderListData{Rows: result.Rows, Total: result.Total}, nil
}
