package impl

import (
	"fmt"
	"time"

	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopOrderDao struct{}

func NewShopOrderDao() dao.IShopOrderDao {
	return &ShopOrderDao{}
}

func (d *ShopOrderDao) Create(c *gin.Context, order *models.Order) (*models.Order, error) {
	db := c.Value("db").(*gorm.DB)
	if err := db.WithContext(c).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (d *ShopOrderDao) GetByID(c *gin.Context, id int64) (*models.Order, error) {
	db := c.Value("db").(*gorm.DB)
	var order models.Order
	err := db.WithContext(c).Where("id = ? AND state = 0", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (d *ShopOrderDao) GetByOrderNo(c *gin.Context, orderNo string) (*models.Order, error) {
	db := c.Value("db").(*gorm.DB)
	var order models.Order
	err := db.WithContext(c).Where("order_no = ? AND state = 0", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (d *ShopOrderDao) List(c *gin.Context, query *models.OrderQuery) (*models.OrderListData, error) {
	db := c.Value("db").(*gorm.DB)

	var total int64
	var orders []*models.Order

	// 构建查询
	q := db.WithContext(c).Model(&models.Order{}).Where("state = 0")

	if query.UserID > 0 {
		q = q.Where("user_id = ?", query.UserID)
	}
	if query.Status != nil {
		q = q.Where("status = ?", *query.Status)
	}
	if query.OrderNo != "" {
		q = q.Where("order_no LIKE ?", "%"+query.OrderNo+"%")
	}

	// 统计总数
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	page := query.Page
	if page < 1 {
		page = 1
	}
	size := query.Size
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	if err := q.Order("create_time DESC").Offset(int(offset)).Limit(int(size)).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &models.OrderListData{
		Rows:  d.toOrderVOList(orders),
		Total: total,
	}, nil
}

func (d *ShopOrderDao) toOrderVOList(orders []*models.Order) []*models.OrderVO {
	result := make([]*models.OrderVO, len(orders))
	for i, order := range orders {
		result[i] = &models.OrderVO{Order: *order}
	}
	return result
}

func (d *ShopOrderDao) UpdateStatus(c *gin.Context, id int64, status int32, version int32) (int64, error) {
	db := c.Value("db").(*gorm.DB)

	var updates map[string]interface{}
	now := time.Now().Format("2006-01-02 15:04:05")

	switch status {
	case models.OrderStatusPaid:
		updates = map[string]interface{}{
			"status":      status,
			"pay_time":    now,
			"version":     version + 1,
			"update_time": now,
		}
	case models.OrderStatusShipped:
		updates = map[string]interface{}{
			"status":      status,
			"ship_time":   now,
			"version":     version + 1,
			"update_time": now,
		}
	case models.OrderStatusCompleted:
		updates = map[string]interface{}{
			"status":        status,
			"complete_time": now,
			"version":       version + 1,
			"update_time":   now,
		}
	default:
		updates = map[string]interface{}{
			"status":      status,
			"version":     version + 1,
			"update_time": now,
		}
	}

	result := db.WithContext(c).Model(&models.Order{}).
		Where("id = ? AND version = ? AND state = 0", id, version).
		Updates(updates)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (d *ShopOrderDao) Cancel(c *gin.Context, id int64, reason string, version int32) (int64, error) {
	db := c.Value("db").(*gorm.DB)
	now := time.Now().Format("2006-01-02 15:04:05")

	result := db.WithContext(c).Model(&models.Order{}).
		Where("id = ? AND version = ? AND status = ? AND state = 0",
			id, version, models.OrderStatusPending).
		Updates(map[string]interface{}{
			"status":        models.OrderStatusCancelled,
			"cancel_time":   now,
			"cancel_reason": reason,
			"version":       version + 1,
			"update_time":   now,
		})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// GetStatistics 获取订单统计
func (d *ShopOrderDao) GetStatistics(c *gin.Context, userID int64) (*models.OrderStatistics, error) {
	db := c.Value("db").(*gorm.DB)

	stats := &models.OrderStatistics{}

	// 查询待付款
	db.WithContext(c).Model(&models.Order{}).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusPending).
		Count(&stats.PendingPay)

	// 查询待发货（已支付）
	db.WithContext(c).Model(&models.Order{}).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusPaid).
		Count(&stats.PendingSend)

	// 查询待收货（已发货）
	db.WithContext(c).Model(&models.Order{}).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusShipped).
		Count(&stats.PendingReceive)

	// 查询已完成
	db.WithContext(c).Model(&models.Order{}).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusCompleted).
		Count(&stats.Completed)

	// 查询已取消
	db.WithContext(c).Model(&models.Order{}).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusCancelled).
		Count(&stats.Cancelled)

	return stats, nil
}

// ShopOrderItemDao 订单商品明细DAO实现
type ShopOrderItemDao struct{}

func NewShopOrderItemDao() dao.IShopOrderItemDao {
	return &ShopOrderItemDao{}
}

func (d *ShopOrderItemDao) BatchCreate(c *gin.Context, items []*models.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	db := c.Value("db").(*gorm.DB)
	return db.WithContext(c).CreateInBatches(items, 100).Error
}

func (d *ShopOrderItemDao) GetByOrderID(c *gin.Context, orderID int64) ([]*models.OrderItem, error) {
	db := c.Value("db").(*gorm.DB)
	var items []*models.OrderItem
	err := db.WithContext(c).Where("order_id = ? AND state = 0", orderID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (d *ShopOrderItemDao) GetByOrderNo(c *gin.Context, orderNo string) ([]*models.OrderItem, error) {
	db := c.Value("db").(*gorm.DB)
	var items []*models.OrderItem
	err := db.WithContext(c).Where("order_no = ? AND state = 0", orderNo).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// 生成订单号
func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%04d%02d%02d%02d%02d%04d",
		now.Format("200601"),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/10000%10000)
}
