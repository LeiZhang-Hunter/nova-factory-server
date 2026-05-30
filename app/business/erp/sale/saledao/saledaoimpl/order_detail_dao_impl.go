package saledaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderDetailDaoImpl ERP 订单明细 DAO。
type OrderDetailDaoImpl struct {
	db    *gorm.DB
	table string
}

// erpOrderDetailRow ERP 订单明细表行模型，显式绑定真实表字段名。
type erpOrderDetailRow struct {
	ID             uint64     `gorm:"column:id"`
	OrderID        uint64     `gorm:"column:order_id"`
	Tid            string     `gorm:"column:tid"`
	OID            string     `gorm:"column:oid"`
	Barcode        string     `gorm:"column:barcode"`
	EShopGoodsID   string     `gorm:"column:eshop_goods_id"`
	OuterIID       string     `gorm:"column:outer_iid"`
	EShopGoodsName string     `gorm:"column:eshop_goods_name"`
	EShopSkuID     string     `gorm:"column:eshop_sku_id"`
	EShopSkuName   string     `gorm:"column:eshop_sku_name"`
	NumIID         int64      `gorm:"column:num_iid"`
	SkuID          int64      `gorm:"column:sku_id"`
	Num            float64    `gorm:"column:num"`
	Payment        float64    `gorm:"column:payment"`
	PicPath        string     `gorm:"column:pic_path"`
	Weight         float64    `gorm:"column:weight"`
	Size           float64    `gorm:"column:size"`
	UnitID         int64      `gorm:"column:unit_id"`
	UnitQty        float64    `gorm:"column:unit_qty"`
	DeptID         int64      `gorm:"column:dept_id"`
	CreateBy       int64      `gorm:"column:create_by"`
	CreateTime     *time.Time `gorm:"column:create_time"`
	UpdateBy       int64      `gorm:"column:update_by"`
	UpdateTime     *time.Time `gorm:"column:update_time"`
	State          int32      `gorm:"column:state"`
}

// NewOrderDetailDao 创建 ERP 订单明细 DAO。
func NewOrderDetailDao(db *gorm.DB) *OrderDetailDaoImpl {
	return &OrderDetailDaoImpl{
		db:    db,
		table: "erp_order_detail",
	}
}

// BatchCreate 批量创建订单明细记录。
func (d *OrderDetailDaoImpl) BatchCreate(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, details []*salemodels.OrderDetailSet) error {
	if len(details) == 0 {
		return nil
	}
	if err := d.validateOIDs(tx, orderID, details); err != nil {
		return err
	}
	rows := make([]*erpOrderDetailRow, 0, len(details))
	for _, item := range details {
		if item == nil {
			continue
		}
		row := &erpOrderDetailRow{
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

// DeleteByOrderID 按订单ID删除明细记录。
func (d *OrderDetailDaoImpl) DeleteByOrderID(tx *gorm.DB, orderID uint64) error {
	return tx.Table(d.table).
		Where("order_id = ?", orderID).
		Delete(&erpOrderDetailRow{}).Error
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
func (d *OrderDetailDaoImpl) DeleteByTidAndOIDs(tx *gorm.DB, tid string, details []*salemodels.OrderDetailSet) error {
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
		Delete(&erpOrderDetailRow{}).Error
}

// ListByOrderIDs 按订单ID集合查询明细记录。
func (d *OrderDetailDaoImpl) ListByOrderIDs(c *gin.Context, orderIDs []uint64) ([]*salemodels.OrderDetail, error) {
	return d.listByOrderIDsWithDB(c, d.db.WithContext(c), orderIDs)
}

func (d *OrderDetailDaoImpl) listByOrderIDsWithDB(c *gin.Context, db *gorm.DB, orderIDs []uint64) ([]*salemodels.OrderDetail, error) {
	if len(orderIDs) == 0 {
		return []*salemodels.OrderDetail{}, nil
	}
	rows := make([]*salemodels.OrderDetail, 0)
	rowList := make([]*erpOrderDetailRow, 0)
	if err := db.Table(d.table).
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rowList).Error; err != nil {
		return nil, err
	}
	for _, row := range rowList {
		if row == nil {
			continue
		}
		item := row.toModel()
		rows = append(rows, &item)
	}
	return rows, nil
}

// validateOIDs 校验订单明细 OID 的唯一性。
func (d *OrderDetailDaoImpl) validateOIDs(tx *gorm.DB, orderID uint64, details []*salemodels.OrderDetailSet) error {
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

	var exists erpOrderDetailRow
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

// toModel 将真实表结构行模型转换为领域模型。
func (r *erpOrderDetailRow) toModel() salemodels.OrderDetail {
	if r == nil {
		return salemodels.OrderDetail{}
	}
	return salemodels.OrderDetail{
		ID:             r.ID,
		OrderID:        r.OrderID,
		Tid:            r.Tid,
		OID:            r.OID,
		Barcode:        r.Barcode,
		EShopGoodsID:   r.EShopGoodsID,
		OuterIID:       r.OuterIID,
		EShopGoodsName: r.EShopGoodsName,
		EShopSkuID:     r.EShopSkuID,
		EShopSkuName:   r.EShopSkuName,
		NumIID:         r.NumIID,
		SkuID:          r.SkuID,
		Num:            r.Num,
		Payment:        r.Payment,
		PicPath:        r.PicPath,
		Weight:         r.Weight,
		Size:           r.Size,
		UnitID:         r.UnitID,
		UnitQty:        r.UnitQty,
		DeptID:         r.DeptID,
		BaseEntity: baize.BaseEntity{
			CreateBy:   r.CreateBy,
			CreateTime: r.CreateTime,
			UpdateBy:   r.UpdateBy,
			UpdateTime: r.UpdateTime,
		},
		State: r.State,
	}
}
