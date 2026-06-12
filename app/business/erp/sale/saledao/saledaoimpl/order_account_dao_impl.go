package saledaoimpl

import (
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderAccountDaoImpl ERP 订单账户 DAO。
type OrderAccountDaoImpl struct {
	db    *gorm.DB
	table string
}

// erpOrderAccountRow ERP 订单账户表行模型，显式绑定真实表字段名。
type erpOrderAccountRow struct {
	ID          uint64     `gorm:"column:id"`
	OrderID     uint64     `gorm:"column:order_id"`
	Tid         string     `gorm:"column:tid"`
	FinanceCode string     `gorm:"column:finance_code"`
	Total       float64    `gorm:"column:total"`
	DeptID      int64      `gorm:"column:dept_id"`
	CreateBy    int64      `gorm:"column:create_by"`
	CreateTime  *time.Time `gorm:"column:create_time"`
	UpdateBy    int64      `gorm:"column:update_by"`
	UpdateTime  *time.Time `gorm:"column:update_time"`
	State       int32      `gorm:"column:state"`
}

// NewOrderAccountDao 创建 ERP 订单账户 DAO。
func NewOrderAccountDao(db *gorm.DB) saledao.IOrderAccountDao {
	return newOrderAccountDaoImpl(db)
}

// newOrderAccountDaoImpl 创建 ERP 订单账户 DAO 具体实现。
func newOrderAccountDaoImpl(db *gorm.DB) *OrderAccountDaoImpl {
	return &OrderAccountDaoImpl{
		db:    db,
		table: "erp_order_account",
	}
}

// BatchCreate 批量创建订单账户记录。
func (d *OrderAccountDaoImpl) BatchCreate(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, accounts []*salemodels.OrderAccountSet) error {
	if len(accounts) == 0 {
		return nil
	}
	rows := make([]*erpOrderAccountRow, 0, len(accounts))
	for _, item := range accounts {
		if item == nil {
			continue
		}
		row := &erpOrderAccountRow{
			OrderID:     orderID,
			Tid:         tid,
			FinanceCode: item.FinanceCode,
			Total:       item.Total,
			State:       commonStatus.NORMAL,
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

// BatchCreateByOrder 批量创建事件同步转换后的订单账户。
func (d *OrderAccountDaoImpl) BatchCreateByOrder(tx *gorm.DB, orderID uint64, order *salemodels.Order, now *time.Time) error {
	if order == nil || len(order.Accounts) == 0 {
		return nil
	}
	rows := make([]*erpOrderAccountRow, 0, len(order.Accounts))
	for _, item := range order.Accounts {
		if item == nil {
			continue
		}
		row := &erpOrderAccountRow{
			OrderID:     orderID,
			Tid:         order.Tid,
			FinanceCode: item.FinanceCode,
			Total:       item.Total,
			DeptID:      order.DeptID,
			CreateBy:    order.CreateBy,
			CreateTime:  now,
			UpdateBy:    order.UpdateBy,
			UpdateTime:  now,
			State:       commonStatus.NORMAL,
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.table).Create(&rows).Error
}

// DeleteByOrderID 按订单ID删除账户记录。
func (d *OrderAccountDaoImpl) DeleteByOrderID(tx *gorm.DB, orderID uint64) error {
	return tx.Table(d.table).
		Where("order_id = ?", orderID).
		Delete(&erpOrderAccountRow{}).Error
}

// DeleteByOrderIDs 按订单ID集合删除账户记录。
func (d *OrderAccountDaoImpl) DeleteByOrderIDs(tx *gorm.DB, orderIDs []uint64) error {
	if len(orderIDs) == 0 {
		return nil
	}
	return tx.Table(d.table).
		Where("order_id IN ?", orderIDs).
		Delete(nil).Error
}

// ListByOrderIDs 按订单ID集合查询账户记录。
func (d *OrderAccountDaoImpl) ListByOrderIDs(c *gin.Context, orderIDs []uint64) ([]*salemodels.OrderAccount, error) {
	return d.listByOrderIDsWithDB(c, d.db.WithContext(c), orderIDs)
}

func (d *OrderAccountDaoImpl) listByOrderIDsWithDB(c *gin.Context, db *gorm.DB, orderIDs []uint64) ([]*salemodels.OrderAccount, error) {
	if len(orderIDs) == 0 {
		return []*salemodels.OrderAccount{}, nil
	}
	rows := make([]*salemodels.OrderAccount, 0)
	rowList := make([]*erpOrderAccountRow, 0)
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

// toModel 将真实表结构行模型转换为领域模型。
func (r *erpOrderAccountRow) toModel() salemodels.OrderAccount {
	if r == nil {
		return salemodels.OrderAccount{}
	}
	return salemodels.OrderAccount{
		ID:          r.ID,
		OrderID:     r.OrderID,
		Tid:         r.Tid,
		FinanceCode: r.FinanceCode,
		Total:       r.Total,
		DeptID:      r.DeptID,
		BaseEntity: baize.BaseEntity{
			CreateBy:   r.CreateBy,
			CreateTime: r.CreateTime,
			UpdateBy:   r.UpdateBy,
			UpdateTime: r.UpdateTime,
		},
		State: r.State,
	}
}
