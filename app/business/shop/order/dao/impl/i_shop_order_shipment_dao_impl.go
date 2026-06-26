package impl

import (
	"fmt"
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

// BatchInsert 批量插入发货物流记录，已存在的 outsid 自动跳过（增量同步）。
func (d *OrderShipmentDaoImpl) BatchInsert(tx *gorm.DB, shipments []*models.OrderShipmentSet) error {
	if len(shipments) == 0 {
		return nil
	}

	// 收集所有非空 outsid，查出现存的
	outsids := make([]string, 0, len(shipments))
	for _, s := range shipments {
		if s == nil {
			continue
		}
		o := strings.TrimSpace(s.Outsid)
		if o != "" {
			outsids = append(outsids, o)
		}
	}
	existSet := make(map[string]struct{}, len(outsids))
	if len(outsids) > 0 {
		var existing []struct {
			Outsid string `gorm:"column:outsid"`
		}
		if err := tx.Table(d.table).
			Select("outsid").
			Where("outsid IN ?", outsids).
			Where("state = ?", commonStatus.NORMAL).
			Find(&existing).Error; err != nil {
			return err
		}
		for _, e := range existing {
			existSet[e.Outsid] = struct{}{}
		}
	}

	now := time.Now()
	rows := make([]*models.OrderShipment, 0, len(shipments))
	for _, s := range shipments {
		if s == nil {
			continue
		}
		o := strings.TrimSpace(s.Outsid)
		if o == "" {
			continue
		}
		if _, exists := existSet[o]; exists {
			continue
		}
		row := &models.OrderShipment{
			OrderID:     s.OrderID,
			Tid:         strings.TrimSpace(s.Tid),
			Issplit:     s.Issplit,
			Outsid:      o,
			Companycode: strings.TrimSpace(s.Companycode),
			SubTid:      strings.TrimSpace(s.SubTid),
			OID:         strings.TrimSpace(s.OID),
			Qty:         s.Qty,
			State:       commonStatus.NORMAL,
		}
		row.CreateTime = &now
		row.UpdateTime = &now
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.table).Create(&rows).Error
}

// ExistsByOutsidTx 按物流单号查重。
func (d *OrderShipmentDaoImpl) ExistsByOutsidTx(tx *gorm.DB, outsid string) (bool, error) {
	outsid = strings.TrimSpace(outsid)
	if outsid == "" {
		return false, nil
	}
	var count int64
	if err := tx.Table(d.table).
		Where("outsid = ?", outsid).
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
