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

type PurchaseOrderItemDaoImpl struct {
	db *gorm.DB
}

func NewPurchaseOrderItemDao(db *gorm.DB) purchasedao.IPurchaseOrderItemDao {
	return &PurchaseOrderItemDaoImpl{db: db}
}

func (d *PurchaseOrderItemDaoImpl) Create(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := purchasemodels.PurchaseOrderItemUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_purchase_order_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *PurchaseOrderItemDaoImpl) Update(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.OrderID > 0 {
		updates["order_id"] = req.OrderID
	}
	if req.ProductID > 0 {
		updates["product_id"] = req.ProductID
	}
	if req.ProductUnitID > 0 {
		updates["product_unit_id"] = req.ProductUnitID
	}
	if req.ProductPrice != 0 {
		updates["product_price"] = req.ProductPrice
	}
	if req.Count != 0 {
		updates["count"] = req.Count
	}
	if req.TotalPrice != 0 {
		updates["total_price"] = req.TotalPrice
	}
	if req.TaxPercent != 0 {
		updates["tax_percent"] = req.TaxPercent
	}
	if req.TaxPrice != 0 {
		updates["tax_price"] = req.TaxPrice
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.InCount != 0 {
		updates["in_count"] = req.InCount
	}
	if req.ReturnCount != 0 {
		updates["return_count"] = req.ReturnCount
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_purchase_order_item").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *PurchaseOrderItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_purchase_order_item").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *PurchaseOrderItemDaoImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrderItem, error) {
	item := new(purchasemodels.PurchaseOrderItem)
	if err := d.db.WithContext(c).Table("erp_purchase_order_item").
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

func (d *PurchaseOrderItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseOrderItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(purchasemodels.PurchaseOrderItem)
	if err := d.db.WithContext(c).Table("erp_purchase_order_item").
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

func (d *PurchaseOrderItemDaoImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error) {
	if req == nil {
		req = new(purchasemodels.PurchaseOrderItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_purchase_order_item").Where("state = ?", commonStatus.NORMAL)
	db = applyPurchaseOrderItemFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]purchasemodels.PurchaseOrderItem, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*purchasemodels.PurchaseOrderItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &purchasemodels.PurchaseOrderItemListData{Rows: result, Total: total}, nil
}

func (d *PurchaseOrderItemDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderItemListData{Rows: result.Rows, Total: result.Total}, nil
}
