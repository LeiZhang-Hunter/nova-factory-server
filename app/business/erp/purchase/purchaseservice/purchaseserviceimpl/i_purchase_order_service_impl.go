package purchaseserviceimpl

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	purchaseOrderTable     = "erp_purchase_order"
	purchaseOrderItemTable = "erp_purchase_order_items"
)

// PurchaseOrderServiceImpl 提供 ERP 采购订单业务实现。
type PurchaseOrderServiceImpl struct {
	db          *gorm.DB
	dao         purchasedao.IPurchaseOrderDao
	productDao  masterdao.IProductDao
	supplierDao masterdao.ISupplierDao
	accountDao  masterdao.IAccountDao
}

// NewPurchaseOrderService 创建 ERP 采购订单服务。
func NewPurchaseOrderService(
	db *gorm.DB,
	dao purchasedao.IPurchaseOrderDao,
	productDao masterdao.IProductDao,
	supplierDao masterdao.ISupplierDao,
	accountDao masterdao.IAccountDao,
) purchaseservice.IPurchaseOrderService {
	return &PurchaseOrderServiceImpl{
		db:          db,
		dao:         dao,
		productDao:  productDao,
		supplierDao: supplierDao,
		accountDao:  accountDao,
	}
}

// Create 创建 ERP 采购订单。
func (s *PurchaseOrderServiceImpl) Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error) {
	if err := s.validateSaveReq(c, req); err != nil {
		return nil, err
	}
	orderTime, err := parseTimeValue(req.OrderTime)
	if err != nil {
		return nil, err
	}
	items, err := s.buildOrderItems(c, req.Items)
	if err != nil {
		return nil, err
	}
	no, err := s.generateOrderNo(c)
	if err != nil {
		return nil, err
	}
	order := &purchasemodels.PurchaseOrder{
		No:              no,
		Status:          auditStatusProcess,
		SupplierID:      req.SupplierID,
		AccountID:       req.AccountID,
		OrderTime:       orderTime,
		DiscountPercent: roundAmount(req.DiscountPercent),
		DepositPrice:    roundAmount(req.DepositPrice),
		FileURL:         strings.TrimSpace(req.FileURL),
		Remark:          strings.TrimSpace(req.Remark),
	}
	s.calculateTotals(order, items)
	order.ID = snowflake.GenID()
	order.DeptID = baizeContext.GetDeptId(c)
	order.State = commonStatus.NORMAL
	order.SetCreateBy(baizeContext.GetUserId(c))

	tx := s.db.WithContext(c).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := tx.Table(purchaseOrderTable).Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, item := range items {
		item.OrderID = order.ID
	}
	if len(items) > 0 {
		if err := tx.Table(purchaseOrderItemTable).Create(&items).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, order.ID)
}

// Update 更新 ERP 采购订单。
func (s *PurchaseOrderServiceImpl) Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	exists, err := s.mustGetOrder(c, req.ID)
	if err != nil {
		return nil, err
	}
	if exists.Status == auditStatusApprove {
		return nil, fmt.Errorf("采购订单[%s]已审核，不能修改", exists.No)
	}
	if err := s.validateSaveReq(c, req); err != nil {
		return nil, err
	}
	orderTime, err := parseTimeValue(req.OrderTime)
	if err != nil {
		return nil, err
	}
	items, err := s.buildOrderItems(c, req.Items)
	if err != nil {
		return nil, err
	}
	order := &purchasemodels.PurchaseOrder{
		ID:              exists.ID,
		No:              exists.No,
		Status:          exists.Status,
		SupplierID:      req.SupplierID,
		AccountID:       req.AccountID,
		OrderTime:       orderTime,
		DiscountPercent: roundAmount(req.DiscountPercent),
		DepositPrice:    roundAmount(req.DepositPrice),
		FileURL:         strings.TrimSpace(req.FileURL),
		Remark:          strings.TrimSpace(req.Remark),
		InCount:         exists.InCount,
		ReturnCount:     exists.ReturnCount,
		State:           exists.State,
	}
	s.calculateTotals(order, items)
	if err := prepareOrderUpdate(order, c); err != nil {
		return nil, err
	}

	tx := s.db.WithContext(c).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := tx.Table(purchaseOrderTable).
		Where("id = ? AND state = ?", order.ID, commonStatus.NORMAL).
		Updates(map[string]any{
			"supplier_id":         order.SupplierID,
			"account_id":          order.AccountID,
			"order_time":          order.OrderTime,
			"total_count":         order.TotalCount,
			"total_price":         order.TotalPrice,
			"total_product_price": order.TotalProductPrice,
			"total_tax_price":     order.TotalTaxPrice,
			"discount_percent":    order.DiscountPercent,
			"discount_price":      order.DiscountPrice,
			"deposit_price":       order.DepositPrice,
			"file_url":            order.FileURL,
			"remark":              order.Remark,
			"update_by":           order.UpdateBy,
			"update_time":         order.UpdateTime,
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, item := range items {
		item.OrderID = order.ID
	}
	if err := replaceChildren(tx, c, purchaseOrderItemTable, "order_id", order.ID, items); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, order.ID)
}

// UpdateStatus 更新 ERP 采购订单状态。
func (s *PurchaseOrderServiceImpl) UpdateStatus(c *gin.Context, req *purchasemodels.PurchaseOrderStatusReq) error {
	if req == nil || req.ID <= 0 {
		return errors.New("id不能为空")
	}
	if req.Status != auditStatusProcess && req.Status != auditStatusApprove {
		return errors.New("状态不正确")
	}
	order, err := s.mustGetOrder(c, req.ID)
	if err != nil {
		return err
	}
	if order.Status == req.Status {
		if req.Status == auditStatusApprove {
			return errors.New("采购订单已审核")
		}
		return errors.New("采购订单已是未审核状态")
	}
	if req.Status == auditStatusProcess {
		if order.InCount > 0 {
			return errors.New("采购订单已存在入库记录，不能反审核")
		}
		if order.ReturnCount > 0 {
			return errors.New("采购订单已存在退货记录，不能反审核")
		}
	}
	now := time.Now()
	return s.db.WithContext(c).
		Table(purchaseOrderTable).
		Where("id = ? AND state = ?", req.ID, commonStatus.NORMAL).
		Updates(map[string]any{
			"status":      req.Status,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": &now,
		}).Error
}

// DeleteByIDs 删除 ERP 采购订单。
func (s *PurchaseOrderServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	orders := make([]*purchasemodels.PurchaseOrder, 0)
	if err := s.db.WithContext(c).
		Table(purchaseOrderTable).
		Where("id IN ? AND state = ?", ids, commonStatus.NORMAL).
		Find(&orders).Error; err != nil {
		return err
	}
	for _, order := range orders {
		if order.Status == auditStatusApprove {
			return fmt.Errorf("采购订单[%s]已审核，不能删除", order.No)
		}
	}
	now := time.Now()
	tx := s.db.WithContext(c).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	updateMap := map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": &now,
	}
	if err := tx.Table(purchaseOrderTable).
		Where("id IN ? AND state = ?", ids, commonStatus.NORMAL).
		Updates(updateMap).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table(purchaseOrderItemTable).
		Where("order_id IN ? AND state = ?", ids, commonStatus.NORMAL).
		Updates(updateMap).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// GetByID 查询 ERP 采购订单详情。
func (s *PurchaseOrderServiceImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	order, err := s.dao.GetByID(c, id)
	if err != nil || order == nil {
		return order, err
	}
	items := make([]*purchasemodels.PurchaseOrderItem, 0)
	if err := s.db.WithContext(c).
		Table(purchaseOrderItemTable).
		Where("order_id = ? AND state = ?", id, commonStatus.NORMAL).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

// List 查询 ERP 采购订单列表。
func (s *PurchaseOrderServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error) {
	result, err := s.dao.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderListData{Rows: result.Rows, Total: result.Total}, nil
}

func (s *PurchaseOrderServiceImpl) validateSaveReq(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if req.SupplierID <= 0 {
		return errors.New("供应商不能为空")
	}
	if strings.TrimSpace(req.OrderTime) == "" {
		return errors.New("下单时间不能为空")
	}
	if len(req.Items) == 0 {
		return errors.New("采购订单明细不能为空")
	}
	supplier, err := s.supplierDao.GetByID(c, req.SupplierID)
	if err != nil {
		return err
	}
	if supplier == nil {
		return errors.New("供应商不存在")
	}
	if req.AccountID > 0 {
		account, err := s.accountDao.GetByID(c, req.AccountID)
		if err != nil {
			return err
		}
		if account == nil {
			return errors.New("结算账户不存在")
		}
	}
	return nil
}

func (s *PurchaseOrderServiceImpl) buildOrderItems(c *gin.Context, reqItems []*purchasemodels.PurchaseOrderItemUpsert) ([]*purchasemodels.PurchaseOrderItem, error) {
	items := make([]*purchasemodels.PurchaseOrderItem, 0, len(reqItems))
	for idx, reqItem := range reqItems {
		if reqItem == nil {
			return nil, fmt.Errorf("第%d条明细不能为空", idx+1)
		}
		if reqItem.ProductID <= 0 {
			return nil, fmt.Errorf("第%d条明细产品不能为空", idx+1)
		}
		if reqItem.Count <= 0 {
			return nil, fmt.Errorf("第%d条明细采购数量必须大于0", idx+1)
		}
		if reqItem.ProductPrice < 0 {
			return nil, fmt.Errorf("第%d条明细采购单价不能小于0", idx+1)
		}
		product, err := s.productDao.GetByID(c, reqItem.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, fmt.Errorf("第%d条明细产品不存在", idx+1)
		}
		item := &purchasemodels.PurchaseOrderItem{
			ProductID:     reqItem.ProductID,
			ProductUnitID: product.UnitId,
			ProductPrice:  roundAmount(reqItem.ProductPrice),
			Count:         reqItem.Count,
			TaxPercent:    roundAmount(reqItem.TaxPercent),
			Remark:        strings.TrimSpace(reqItem.Remark),
		}
		item.TotalPrice = roundAmount(item.ProductPrice * item.Count)
		item.TaxPrice = calculatePercentAmount(item.TotalPrice, item.TaxPercent)
		item.ID = snowflake.GenID()
		item.DeptID = baizeContext.GetDeptId(c)
		item.State = commonStatus.NORMAL
		item.SetCreateBy(baizeContext.GetUserId(c))
		items = append(items, item)
	}
	return items, nil
}

func (s *PurchaseOrderServiceImpl) calculateTotals(order *purchasemodels.PurchaseOrder, items []*purchasemodels.PurchaseOrderItem) {
	order.TotalCount = 0
	order.TotalProductPrice = 0
	order.TotalTaxPrice = 0
	for _, item := range items {
		order.TotalCount += item.Count
		order.TotalProductPrice += item.TotalPrice
		order.TotalTaxPrice += item.TaxPrice
	}
	order.TotalCount = roundAmount(order.TotalCount)
	order.TotalProductPrice = roundAmount(order.TotalProductPrice)
	order.TotalTaxPrice = roundAmount(order.TotalTaxPrice)
	totalBeforeDiscount := roundAmount(order.TotalProductPrice + order.TotalTaxPrice)
	order.DiscountPrice = calculatePercentAmount(totalBeforeDiscount, order.DiscountPercent)
	order.TotalPrice = roundAmount(totalBeforeDiscount - order.DiscountPrice)
}

func (s *PurchaseOrderServiceImpl) generateOrderNo(c *gin.Context) (string, error) {
	for i := 0; i < 5; i++ {
		no := generateNo(prefixPurchaseOrder)
		exists, err := s.dao.GetByColumn(c, "no", no)
		if err != nil {
			return "", err
		}
		if exists == nil {
			return no, nil
		}
	}
	return "", errors.New("生成采购订单号失败，请重试")
}

func (s *PurchaseOrderServiceImpl) mustGetOrder(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error) {
	order, err := s.dao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("采购订单不存在")
	}
	return order, nil
}

const (
	auditStatusProcess  int32 = 10
	auditStatusApprove  int32 = 20
	prefixPurchaseOrder       = "CGDD"
)

func generateNo(prefix string) string {
	suffix := snowflake.GenID()
	if suffix < 0 {
		suffix = -suffix
	}
	return strings.TrimSpace(prefix) + time.Now().Format("20060102") + fmt.Sprintf("%06d", suffix%1000000)
}

func parseTimeValue(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	layouts := []string{time.DateTime, "2006-01-02 15:04", "2006-01-02", time.RFC3339}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("时间格式不正确")
}

func roundAmount(value float64) float64 {
	return math.Round(value*100) / 100
}

func calculatePercentAmount(total, percent float64) float64 {
	if total == 0 || percent == 0 {
		return 0
	}
	return roundAmount(total * percent / 100)
}

func prepareOrderUpdate(order *purchasemodels.PurchaseOrder, c *gin.Context) error {
	if order == nil || order.ID <= 0 {
		return errors.New("id不能为空")
	}
	order.SetUpdateBy(baizeContext.GetUserId(c))
	return nil
}

func replaceChildren(tx *gorm.DB, c *gin.Context, table string, parentColumn string, parentID int64, children any) error {
	if parentID <= 0 {
		return errors.New("主表ID不能为空")
	}
	now := time.Now()
	if err := tx.WithContext(c).
		Table(table).
		Where(parentColumn+" = ? AND state = ?", parentID, commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": &now,
		}).Error; err != nil {
		return err
	}
	rv := reflect.ValueOf(children)
	if rv.Kind() != reflect.Slice || rv.Len() == 0 {
		return nil
	}
	return tx.WithContext(c).Table(table).Create(children).Error
}
