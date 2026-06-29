package impl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderDetailDaoImpl shop 订单明细 DAO。
type OrderDetailDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewOrderDetailDao 创建 shop 订单明细 DAO。
func NewOrderDetailDao(db *gorm.DB) dao.IOrderDetailDao {
	return newOrderDetailDaoImpl(db)
}

// newOrderDetailDaoImpl 创建 shop 订单明细 DAO 具体实现。
func newOrderDetailDaoImpl(db *gorm.DB) *OrderDetailDaoImpl {
	return &OrderDetailDaoImpl{
		db:    db,
		table: "shop_order_detail",
	}
}

// BatchCreate 批量创建订单明细记录。
func (d *OrderDetailDaoImpl) BatchCreate(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, details []*models.OrderDetailSet) error {
	if len(details) == 0 {
		return nil
	}
	if err := d.validateOIDs(tx, orderID, details); err != nil {
		return err
	}
	rows := make([]*models.OrderDetail, 0, len(details))
	for _, item := range details {
		if item == nil {
			continue
		}
		row := &models.OrderDetail{
			OrderID:        orderID,
			Tid:            tid,
			OID:            item.OID,
			Barcode:        item.Barcode,
			EShopGoodsID:   item.EShopGoodsID,
			OuterIID:       item.OuterIID,
			EShopGoodsName: item.EShopGoodsName,
			EShopSkuID:     item.EShopSkuID,
			EShopSkuName:   item.EShopSkuName,
			NumIID:         item.NumIID,
			SkuID:          item.SkuID,
			Num:            item.Num,
			Payment:        item.Payment,
			PicPath:        item.PicPath,
			Weight:         item.Weight,
			Size:           item.Size,
			UnitID:         item.UnitID,
			UnitQty:        item.UnitQty,
			State:          commonStatus.NORMAL,
		}
		now := time.Now()
		row.CreateBy = baizeContext.GetUserId(c)
		row.CreateTime = &now
		row.UpdateBy = baizeContext.GetUserId(c)
		row.UpdateTime = &now
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.table).Create(&rows).Error
}

// BatchCreateByOrder 批量创建事件同步转换后的订单明细。
func (d *OrderDetailDaoImpl) BatchCreateByOrder(tx *gorm.DB, orderID uint64, order *models.Order, now *time.Time) error {
	if order == nil || len(order.Details) == 0 {
		return nil
	}
	rows := make([]*models.OrderDetail, 0, len(order.Details))
	for _, item := range order.Details {
		if item == nil {
			continue
		}
		row := &models.OrderDetail{
			OrderID:        orderID,
			Tid:            order.Tid,
			OID:            item.OID,
			Barcode:        item.Barcode,
			EShopGoodsID:   item.EShopGoodsID,
			OuterIID:       item.OuterIID,
			EShopGoodsName: item.EShopGoodsName,
			EShopSkuID:     item.EShopSkuID,
			EShopSkuName:   item.EShopSkuName,
			NumIID:         item.NumIID,
			SkuID:          item.SkuID,
			Num:            item.Num,
			Payment:        item.Payment,
			PicPath:        item.PicPath,
			Weight:         item.Weight,
			Size:           item.Size,
			UnitID:         item.UnitID,
			UnitQty:        item.UnitQty,
			DeptID:         order.DeptID,
			BaseEntity: baize.BaseEntity{
				CreateBy:   order.CreateBy,
				CreateTime: now,
				UpdateBy:   order.UpdateBy,
				UpdateTime: now,
			},
			State: commonStatus.NORMAL,
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.table).Create(&rows).Error
}

// DeleteByOrderID 按订单ID删除明细记录。
func (d *OrderDetailDaoImpl) DeleteByOrderID(tx *gorm.DB, orderID uint64) error {
	return tx.Table(d.table).
		Where("order_id = ?", orderID).
		Delete(&models.OrderDetail{}).Error
}

// DeleteByOrderIDs 按订单ID集合删除明细记录。
func (d *OrderDetailDaoImpl) DeleteByOrderIDs(tx *gorm.DB, orderIDs []uint64) error {
	if len(orderIDs) == 0 {
		return nil
	}
	return tx.Table(d.table).
		Where("order_id IN ?", orderIDs).
		Delete(nil).Error
}

// DeleteByTidAndOIDs 按订单编号和明细 OID 删除旧明细记录。
func (d *OrderDetailDaoImpl) DeleteByTidAndOIDs(tx *gorm.DB, tid string, details []*models.OrderDetailSet) error {
	tid = strings.TrimSpace(tid)
	if tid == "" || len(details) == 0 {
		return nil
	}
	oidSet := make(map[string]struct{}, len(details))
	oids := make([]string, 0, len(details))
	for _, item := range details {
		if item == nil {
			continue
		}
		oid := strings.TrimSpace(item.OID)
		if oid == "" {
			continue
		}
		if _, ok := oidSet[oid]; ok {
			continue
		}
		oidSet[oid] = struct{}{}
		oids = append(oids, oid)
	}
	if len(oids) == 0 {
		return nil
	}
	return tx.Table(d.table).
		Where("tid = ?", tid).
		Where("oid IN ?", oids).
		Delete(&models.OrderDetail{}).Error
}

// ListByOrderIDs 按订单ID集合查询明细记录。
func (d *OrderDetailDaoImpl) ListByOrderIDs(c *gin.Context, orderIDs []uint64) ([]*models.OrderDetail, error) {
	return d.listByOrderIDsWithDB(c, d.db.WithContext(c), orderIDs)
}

func (d *OrderDetailDaoImpl) ListByOrderID(c *gin.Context, orderIDs uint64) ([]*models.OrderDetail, error) {
	return d.listByOrderIDsWithDB(c, d.db.WithContext(c), []uint64{orderIDs})
}

// ListByTidTx 在事务内按订单编号查询该订单所有有效明细。
func (d *OrderDetailDaoImpl) ListByTidTx(tx *gorm.DB, tid string) ([]*models.OrderDetail, error) {
	tid = strings.TrimSpace(tid)
	if tid == "" {
		return []*models.OrderDetail{}, nil
	}
	rowList := make([]*models.OrderDetail, 0)
	if err := tx.Table(d.table).
		Where("tid = ?", tid).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rowList).Error; err != nil {
		return nil, err
	}
	return rowList, nil
}

// IncrementShippedQty 原子累加指定明细行的已发货数量。
// MySQL UPDATE 语句自动持有行锁，并发安全。
func (d *OrderDetailDaoImpl) IncrementShippedQty(tx *gorm.DB, orderID uint64, oid string, qty float64) error {
	oid = strings.TrimSpace(oid)
	if oid == "" {
		return fmt.Errorf("订单明细oid不能为空")
	}
	if qty <= 0 {
		return fmt.Errorf("发货数量必须大于0")
	}
	result := tx.Table(d.table).
		Where("order_id = ?", orderID).
		Where("oid = ?", oid).
		Where("state = ?", commonStatus.NORMAL).
		Update("shipped_qty", gorm.Expr("shipped_qty + ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("订单明细oid不存在: %s (order_id=%d)", oid, orderID)
	}
	return nil
}

func (d *OrderDetailDaoImpl) listByOrderIDsWithDB(c *gin.Context, db *gorm.DB, orderIDs []uint64) ([]*models.OrderDetail, error) {
	if len(orderIDs) == 0 {
		return []*models.OrderDetail{}, nil
	}
	rowList := make([]*models.OrderDetail, 0)
	if err := db.Table(d.table).
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rowList).Error; err != nil {
		return nil, err
	}
	return rowList, nil
}

// validateOIDs 校验订单明细 OID 的唯一性。
func (d *OrderDetailDaoImpl) validateOIDs(tx *gorm.DB, orderID uint64, details []*models.OrderDetailSet) error {
	oidSet := make(map[string]struct{}, len(details))
	oids := make([]string, 0, len(details))
	for _, item := range details {
		if item == nil {
			continue
		}
		oid := strings.TrimSpace(item.OID)
		if oid == "" {
			continue
		}
		if _, exists := oidSet[oid]; exists {
			return fmt.Errorf("订单明细oid重复: %s", oid)
		}
		oidSet[oid] = struct{}{}
		oids = append(oids, oid)
	}
	if len(oids) == 0 {
		return nil
	}

	var exists *models.OrderDetail
	db := tx.Table(d.table).Where("oid IN ?", oids)
	if orderID > 0 {
		db = db.Where("order_id <> ?", orderID)
	}
	if err := db.First(&exists).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return fmt.Errorf("订单明细oid已存在: %s", exists.OID)
}
