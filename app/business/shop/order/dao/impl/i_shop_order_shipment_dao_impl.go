package impl

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/commonStatus"

	"gorm.io/gorm"
)

// OrderShipmentDaoImpl 订单发货物流记录 DAO。
type OrderShipmentDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewOrderShipmentDao 创建订单发货物流记录 DAO。
func NewOrderShipmentDao(db *gorm.DB) dao.IOrderShipmentDao {
	return &OrderShipmentDaoImpl{
		db:    db,
		table: "shop_order_shipment",
	}
}

// BatchInsert 批量 upsert 发货物流记录：按 (outsid, oid) 组合唯一，存在则更新，不存在则插入。
func (d *OrderShipmentDaoImpl) BatchInsert(tx *gorm.DB, shipments []*models.OrderShipmentSet) error {
	if len(shipments) == 0 {
		return nil
	}

	now := time.Now()
	for _, s := range shipments {
		if s == nil {
			continue
		}
		outsid := strings.TrimSpace(s.Outsid)
		oid := strings.TrimSpace(s.OID)
		if outsid == "" || oid == "" {
			continue
		}

		// 按 (outsid, oid) 组合查现有记录
		var exist models.OrderShipment
		err := tx.Table(d.table).
			Where("outsid = ?", outsid).
			Where("oid = ?", oid).
			Where("state = ?", commonStatus.NORMAL).
			First(&exist).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询物流记录失败(outsid=%s, oid=%s): %w", outsid, oid, err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在 -> 插入
			row := &models.OrderShipment{
				OrderID:     s.OrderID,
				Tid:         strings.TrimSpace(s.Tid),
				Issplit:     s.Issplit,
				Outsid:      outsid,
				Companycode: strings.TrimSpace(s.Companycode),
				SubTid:      strings.TrimSpace(s.SubTid),
				OID:         oid,
				Qty:         s.Qty,
				State:       commonStatus.NORMAL,
			}
			row.ID = uint64(snowflake.GenID())
			row.CreateTime = &now
			row.UpdateTime = &now
			if err := tx.Table(d.table).Create(row).Error; err != nil {
				return fmt.Errorf("插入物流记录失败(outsid=%s, oid=%s): %w", outsid, oid, err)
			}
		} else {
			// 已存在 -> 更新
			updates := map[string]interface{}{
				"order_id":    s.OrderID,
				"tid":         strings.TrimSpace(s.Tid),
				"issplit":     s.Issplit,
				"companycode": strings.TrimSpace(s.Companycode),
				"subtid":      strings.TrimSpace(s.SubTid),
				"oid":         oid,
				"qty":         s.Qty,
				"update_time": now,
			}
			if err := tx.Table(d.table).
				Where("id = ?", exist.ID).
				Updates(updates).Error; err != nil {
				return fmt.Errorf("更新物流记录失败(outsid=%s, oid=%s): %w", outsid, oid, err)
			}
		}
	}
	return nil
}

// ExistsByOutsidTx 按物流单号查重。
func (d *OrderShipmentDaoImpl) ExistsByOutsidTx(tx *gorm.DB, outsid string, oid string) (bool, error) {
	outsid = strings.TrimSpace(outsid)
	if outsid == "" {
		return false, nil
	}
	var count int64
	if err := tx.Table(d.table).
		Where("outsid = ?", outsid).
		Where("oid = ?", oid).
		Where("state = ?", commonStatus.NORMAL).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListByOrderIDTx 按订单 ID 查询所有物流记录。
func (d *OrderShipmentDaoImpl) ListByOrderIDTx(tx *gorm.DB, orderID uint64) ([]*models.OrderShipment, error) {
	if orderID == 0 {
		return []*models.OrderShipment{}, nil
	}
	rows := make([]*models.OrderShipment, 0)
	if err := tx.Table(d.table).
		Where("order_id = ?", orderID).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询订单物流记录失败(order_id=%d): %w", orderID, err)
	}
	return rows, nil
}

// ListByOrderID 按订单 ID 查询所有物流记录。
func (d *OrderShipmentDaoImpl) ListByOrderID(ctx *gin.Context, orderID uint64) ([]*models.OrderShipment, error) {
	if orderID == 0 {
		return []*models.OrderShipment{}, nil
	}
	rows := make([]*models.OrderShipment, 0)
	if err := d.db.Table(d.table).WithContext(ctx).Where("order_id = ?", orderID).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.OrderShipment{}, nil
		}
		return nil, fmt.Errorf("查询订单物流记录失败(order_id=%d): %w", orderID, err)
	}
	return rows, nil
}
