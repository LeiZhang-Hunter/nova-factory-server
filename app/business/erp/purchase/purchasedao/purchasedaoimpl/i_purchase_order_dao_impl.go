package purchasedaoimpl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PurchaseOrderDaoImpl struct {
	db *gorm.DB
}

func NewPurchaseOrderDao(db *gorm.DB) purchasedao.IPurchaseOrderDao {
	return &PurchaseOrderDaoImpl{db: db}
}

func (d *PurchaseOrderDaoImpl) Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := purchasemodels.PurchaseOrderUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_purchase_order").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *PurchaseOrderDaoImpl) Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.OrderTime != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", req.OrderTime, time.Local); err == nil {
			updates["order_time"] = &parsed
		}
	}
	if req.No != "" {
		updates["no"] = req.No
	}
	if req.SupplierID > 0 {
		updates["supplier_id"] = req.SupplierID
	}
	if req.AccountID > 0 {
		updates["account_id"] = req.AccountID
	}
	if req.DiscountPercent != 0 {
		updates["discount_percent"] = req.DiscountPercent
	}
	if req.DepositPrice != 0 {
		updates["deposit_price"] = req.DepositPrice
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_purchase_order").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *PurchaseOrderDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_purchase_order").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *PurchaseOrderDaoImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error) {
	item := new(purchasemodels.PurchaseOrder)
	if err := d.db.WithContext(c).Table("erp_purchase_order").
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
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
	if err := d.db.WithContext(c).Table("erp_purchase_order").
		Where(fmt.Sprintf("%s = ?", column), value).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *PurchaseOrderDaoImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error) {
	if req == nil {
		req = new(purchasemodels.PurchaseOrderQuery)
	}
	db := d.db.WithContext(c).Table("erp_purchase_order").Where("state = ?", commonStatus.NORMAL)
	db = applyPurchaseOrderFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]purchasemodels.PurchaseOrder, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*purchasemodels.PurchaseOrder, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &purchasemodels.PurchaseOrderListData{Rows: result, Total: total}, nil
}

func (d *PurchaseOrderDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderListData{Rows: result.Rows, Total: result.Total}, nil
}
