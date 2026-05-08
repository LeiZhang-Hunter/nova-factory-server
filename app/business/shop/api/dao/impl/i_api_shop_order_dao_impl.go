package impl

import (
	"time"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopOrderDaoImpl 提供订单的数据库访问能力。
type IApiShopOrderDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopOrderDaoImpl 创建订单 DAO 实现。
func NewIApiShopOrderDaoImpl(db *gorm.DB) dao.IApiShopOrderDao {
	return &IApiShopOrderDaoImpl{
		db:        db,
		tableName: "shop_order",
	}
}

// Create 新增订单记录。
func (d *IApiShopOrderDaoImpl) Create(c *gin.Context, order *models.Order) (*models.Order, error) {
	if err := getCurrentDB(c, d.db).WithContext(c).Table(d.tableName).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID 根据ID获取订单。
func (d *IApiShopOrderDaoImpl) GetByID(c *gin.Context, id int64) (*models.Order, error) {
	var order models.Order
	err := d.db.WithContext(c).Table(d.tableName).Where("id = ? AND state = 0", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单。
func (d *IApiShopOrderDaoImpl) GetByOrderNo(c *gin.Context, orderNo string) (*models.Order, error) {
	var order models.Order
	err := d.db.WithContext(c).Table(d.tableName).Where("order_no = ? AND state = 0", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// List 查询订单列表，支持分页和条件筛选。
func (d *IApiShopOrderDaoImpl) List(c *gin.Context, query *models.OrderQuery) (*models.OrderListData, error) {
	var total int64
	var orders []*models.Order

	// 构建查询条件：仅查询未删除记录
	q := d.db.WithContext(c).Table(d.tableName).Where("state = 0")

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

	// 批量获取订单商品明细
	orderIDs := make([]int64, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}
	itemsMap, err := d.getItemsByOrderIDs(c, orderIDs)
	if err != nil {
		return nil, err
	}

	return &models.OrderListData{
		Rows:  d.toOrderVOListWithItems(orders, itemsMap),
		Total: total,
	}, nil
}

// toOrderVOListWithItems 将订单列表转换为带商品明细的视图对象列表。
func (d *IApiShopOrderDaoImpl) toOrderVOListWithItems(orders []*models.Order, itemsMap map[int64][]*models.OrderItem) []*models.OrderVO {
	result := make([]*models.OrderVO, len(orders))
	for i, order := range orders {
		result[i] = &models.OrderVO{
			Order: *order,
			Items: itemsMap[order.ID],
		}
	}
	return result
}

// toOrderVOList 将订单列表转换为视图对象列表。
func (d *IApiShopOrderDaoImpl) toOrderVOList(orders []*models.Order) []*models.OrderVO {
	result := make([]*models.OrderVO, len(orders))
	for i, order := range orders {
		result[i] = &models.OrderVO{Order: *order}
	}
	return result
}

// getItemsByOrderIDs 批量获取订单商品明细，按订单ID分组。
func (d *IApiShopOrderDaoImpl) getItemsByOrderIDs(c *gin.Context, orderIDs []int64) (map[int64][]*models.OrderItem, error) {
	if len(orderIDs) == 0 {
		return make(map[int64][]*models.OrderItem), nil
	}
	var items []*models.OrderItem
	err := d.db.WithContext(c).Table("shop_order_item").
		Where("order_id IN ? AND state = 0", orderIDs).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]*models.OrderItem)
	for _, item := range items {
		result[item.OrderID] = append(result[item.OrderID], item)
	}
	return result, nil
}

// UpdateStatus 更新订单状态，使用乐观锁版本号控制并发。
func (d *IApiShopOrderDaoImpl) UpdateStatus(c *gin.Context, id int64, status int32, version int32) (int64, error) {
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

	result := d.db.WithContext(c).Table(d.tableName).
		Where("id = ? AND version = ? AND state = 0", id, version).
		Updates(updates)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// Cancel 取消订单，仅允许对待支付的订单进行取消。
func (d *IApiShopOrderDaoImpl) Cancel(c *gin.Context, id int64, reason string, version int32) (int64, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	result := d.db.WithContext(c).Table(d.tableName).
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

// GetStatistics 获取用户各状态订单数量统计。
func (d *IApiShopOrderDaoImpl) GetStatistics(c *gin.Context, userID int64) (*models.OrderStatistics, error) {
	stats := &models.OrderStatistics{}

	// 查询待付款
	d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusPending).
		Count(&stats.PendingPay)

	// 查询待发货（已支付）
	d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusPaid).
		Count(&stats.PendingSend)

	// 查询待收货（已发货）
	d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusShipped).
		Count(&stats.PendingReceive)

	// 查询已完成
	d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusCompleted).
		Count(&stats.Completed)

	// 查询已取消
	d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND status = ? AND state = 0", userID, models.OrderStatusCancelled).
		Count(&stats.Cancelled)

	return stats, nil
}
