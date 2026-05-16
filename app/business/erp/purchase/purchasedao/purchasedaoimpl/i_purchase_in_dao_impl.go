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

type PurchaseInDaoImpl struct {
	db *gorm.DB
}

func NewPurchaseInDao(db *gorm.DB) purchasedao.IPurchaseInDao {
	return &PurchaseInDaoImpl{db: db}
}

func (d *PurchaseInDaoImpl) Create(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := purchasemodels.PurchaseInUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_purchase_in").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *PurchaseInDaoImpl) Update(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.InTime != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", req.InTime, time.Local); err == nil {
			updates["in_time"] = &parsed
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
	if req.OrderID > 0 {
		updates["order_id"] = req.OrderID
	}
	if req.OrderNo != "" {
		updates["order_no"] = req.OrderNo
	}
	if req.TotalCount != 0 {
		updates["total_count"] = req.TotalCount
	}
	if req.TotalPrice != 0 {
		updates["total_price"] = req.TotalPrice
	}
	if req.PaymentPrice != 0 {
		updates["payment_price"] = req.PaymentPrice
	}
	if req.TotalProductPrice != 0 {
		updates["total_product_price"] = req.TotalProductPrice
	}
	if req.TotalTaxPrice != 0 {
		updates["total_tax_price"] = req.TotalTaxPrice
	}
	if req.DiscountPercent != 0 {
		updates["discount_percent"] = req.DiscountPercent
	}
	if req.DiscountPrice != 0 {
		updates["discount_price"] = req.DiscountPrice
	}
	if req.OtherPrice != 0 {
		updates["other_price"] = req.OtherPrice
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_purchase_in").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *PurchaseInDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_purchase_in").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *PurchaseInDaoImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseIn, error) {
	item := new(purchasemodels.PurchaseIn)
	if err := d.db.WithContext(c).Table("erp_purchase_in").
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

func (d *PurchaseInDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseIn, error) {
	if column == "" {
		return nil, nil
	}
	item := new(purchasemodels.PurchaseIn)
	if err := d.db.WithContext(c).Table("erp_purchase_in").
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

func (d *PurchaseInDaoImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error) {
	if req == nil {
		req = new(purchasemodels.PurchaseInQuery)
	}
	db := d.db.WithContext(c).Table("erp_purchase_in").Where("state = ?", commonStatus.NORMAL)
	db = applyPurchaseInFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]purchasemodels.PurchaseIn, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*purchasemodels.PurchaseIn, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &purchasemodels.PurchaseInListData{Rows: result, Total: total}, nil
}

func (d *PurchaseInDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseInListData{Rows: result.Rows, Total: result.Total}, nil
}
