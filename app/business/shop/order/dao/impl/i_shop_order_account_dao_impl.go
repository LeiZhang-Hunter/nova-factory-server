package impl

import (
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderAccountDaoImpl shop 订单账户 DAO。
type OrderAccountDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewOrderAccountDao 创建 shop 订单账户 DAO。
func NewOrderAccountDao(db *gorm.DB) dao.IOrderAccountDao {
	return newOrderAccountDaoImpl(db)
}

// newOrderAccountDaoImpl 创建 shop 订单账户 DAO 具体实现。
func newOrderAccountDaoImpl(db *gorm.DB) *OrderAccountDaoImpl {
	return &OrderAccountDaoImpl{
		db:    db,
		table: "shop_order_account",
	}
}

// BatchCreate 批量创建订单账户记录。
func (d *OrderAccountDaoImpl) BatchCreate(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, accounts []*models.OrderAccountSet) error {
	if len(accounts) == 0 {
		return nil
	}
	rows := make([]*models.OrderAccount, 0, len(accounts))
	for _, item := range accounts {
		if item == nil {
			continue
		}
		row := &models.OrderAccount{
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
func (d *OrderAccountDaoImpl) BatchCreateByOrder(tx *gorm.DB, orderID uint64, order *models.Order, now *time.Time) error {
	if order == nil || len(order.Accounts) == 0 {
		return nil
	}
	rows := make([]*models.OrderAccount, 0, len(order.Accounts))
	for _, item := range order.Accounts {
		if item == nil {
			continue
		}
		row := &models.OrderAccount{
			OrderID:     orderID,
			Tid:         order.Tid,
			FinanceCode: item.FinanceCode,
			Total:       item.Total,
			DeptID:      order.DeptID,
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

// DeleteByOrderID 按订单ID删除账户记录。
func (d *OrderAccountDaoImpl) DeleteByOrderID(tx *gorm.DB, orderID uint64) error {
	return tx.Table(d.table).
		Where("order_id = ?", orderID).
		Delete(&models.OrderAccount{}).Error
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
func (d *OrderAccountDaoImpl) ListByOrderIDs(c *gin.Context, orderIDs []uint64) ([]*models.OrderAccount, error) {
	return d.listByOrderIDsWithDB(c, d.db.WithContext(c), orderIDs)
}

func (d *OrderAccountDaoImpl) listByOrderIDsWithDB(c *gin.Context, db *gorm.DB, orderIDs []uint64) ([]*models.OrderAccount, error) {
	if len(orderIDs) == 0 {
		return []*models.OrderAccount{}, nil
	}
	rowList := make([]*models.OrderAccount, 0)
	if err := db.Table(d.table).
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rowList).Error; err != nil {
		return nil, err
	}
	return rowList, nil
}
