package impl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IShopOrderServiceImpl 提供订单相关的业务实现。
type IShopOrderServiceImpl struct {
	orderDao     dao.IShopOrderDao
	orderItemDao dao.IShopOrderItemDao
	userDao      dao.IShopWechatUserDao
}

// NewIShopOrderServiceImpl 创建订单服务实现。
func NewIShopOrderServiceImpl(orderDao dao.IShopOrderDao, orderItemDao dao.IShopOrderItemDao, userDao dao.IShopWechatUserDao) service.IAppShopOrderService {
	return &IShopOrderServiceImpl{
		orderDao:     orderDao,
		orderItemDao: orderItemDao,
		userDao:      userDao,
	}
}

// GenerateOrderNo 生成唯一订单号，格式: ORD+年月日+时分秒+纳秒随机数。
func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%02d%02d%02d%04d",
		now.Format("20060102"),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/10000%10000)
}

// Create 创建订单，包含订单商品明细，支持事务。
func (s *IShopOrderServiceImpl) Create(c *gin.Context, userID int64, req *models.OrderSetReq) (*models.Order, error) {
	if req == nil || len(req.Items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}
	if req.ReceiverName == "" {
		return nil, errors.New("收货人姓名不能为空")
	}
	if req.ReceiverPhone == "" {
		return nil, errors.New("收货人手机号不能为空")
	}
	if req.ReceiverDetailAddress == "" {
		return nil, errors.New("详细地址不能为空")
	}

	// 获取用户信息
	shopUser, err := s.userDao.GetByID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}

	// 计算订单总金额
	var totalAmount float64
	orderItems := make([]*models.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("商品数量必须大于0")
		}
		itemTotal := item.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, &models.OrderItem{
			OrderNo:     "", // 稍后填充
			GoodsID:     item.GoodsID,
			SkuID:       item.SkuID,
			GoodsName:   item.GoodsName,
			SkuName:     item.SkuName,
			ImageURL:    item.ImageURL,
			Price:       item.Price,
			Quantity:    item.Quantity,
			TotalAmount: itemTotal,
		})
	}

	// 生成订单号
	orderNo := GenerateOrderNo()

	// 创建订单
	order := &models.Order{
		OrderNo:               orderNo,
		UserID:                shopUser.ID,
		TotalAmount:           totalAmount,
		PayAmount:             totalAmount,
		Status:                models.OrderStatusPending,
		ReceiverName:          req.ReceiverName,
		ReceiverPhone:         req.ReceiverPhone,
		ReceiverProvince:      req.ReceiverProvince,
		ReceiverCity:          req.ReceiverCity,
		ReceiverDistrict:      req.ReceiverDistrict,
		ReceiverDetailAddress: req.ReceiverDetailAddress,
		Remark:                req.Remark,
		Version:               0,
	}

	// 开启事务
	db := c.Value("db").(*gorm.DB)
	tx := db.WithContext(c).Begin()
	c.Set("db", tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存订单
	createdOrder, err := s.orderDao.Create(c, order)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}

	// 填充订单号到商品明细
	for _, item := range orderItems {
		item.OrderID = createdOrder.ID
		item.OrderNo = orderNo
	}

	// 保存订单商品明细
	if err := s.orderItemDao.BatchCreate(c, orderItems); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("保存订单商品失败: %v", err)
	}

	tx.Commit()
	return createdOrder, nil
}

// GetByID 获取订单详情，包含商品明细。
func (s *IShopOrderServiceImpl) GetByID(c *gin.Context, id int64) (*models.OrderVO, error) {
	if id == 0 {
		return nil, errors.New("订单ID不能为空")
	}

	order, err := s.orderDao.GetByID(c, id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	items, err := s.orderItemDao.GetByOrderID(c, id)
	if err != nil {
		return nil, errors.New("获取订单商品失败")
	}

	return &models.OrderVO{
		Order: *order,
		Items: items,
	}, nil
}

// List 获取当前用户的订单列表。
func (s *IShopOrderServiceImpl) List(c *gin.Context, userID int64, query *models.OrderQuery) (*models.OrderListData, error) {
	if query == nil {
		query = &models.OrderQuery{}
	}

	query.UserID = userID
	return s.orderDao.List(c, query)
}

// UpdateStatus 更新订单状态，验证状态流转合法性。
func (s *IShopOrderServiceImpl) UpdateStatus(c *gin.Context, userID int64, req *models.OrderStatusReq) error {
	if req.ID == 0 {
		return errors.New("订单ID不能为空")
	}

	// 验证用户权限
	order, err := s.orderDao.GetByID(c, req.ID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.UserID != userID {
		return errors.New("无权操作此订单")
	}

	// 验证状态流转
	if !s.isValidStatusTransition(order.Status, req.Status) {
		return errors.New("非法的状态流转")
	}

	rowsAffected, err := s.orderDao.UpdateStatus(c, req.ID, req.Status, order.Version)
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// Cancel 取消订单，仅允许对待支付的订单进行取消。
func (s *IShopOrderServiceImpl) Cancel(c *gin.Context, userID int64, id int64, reason string) error {
	if id == 0 {
		return errors.New("订单ID不能为空")
	}

	order, err := s.orderDao.GetByID(c, id)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.UserID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("只能取消待支付的订单")
	}

	rowsAffected, err := s.orderDao.Cancel(c, id, reason, order.Version)
	if err != nil {
		return fmt.Errorf("取消订单失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// ConfirmReceive 确认收货，将已发货订单标记为已完成。
func (s *IShopOrderServiceImpl) ConfirmReceive(c *gin.Context, userID int64, id int64) error {
	if id == 0 {
		return errors.New("订单ID不能为空")
	}

	order, err := s.orderDao.GetByID(c, id)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.UserID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusShipped {
		return errors.New("只能确认已发货的订单")
	}

	rowsAffected, err := s.orderDao.UpdateStatus(c, id, models.OrderStatusCompleted, order.Version)
	if err != nil {
		return fmt.Errorf("确认收货失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// GetStatistics 获取当前用户各状态订单数量统计。
func (s *IShopOrderServiceImpl) GetStatistics(c *gin.Context, userID int64) (*models.OrderStatistics, error) {
	return s.orderDao.GetStatistics(c, userID)
}

// isValidStatusTransition 验证订单状态流转是否合法。
func (s *IShopOrderServiceImpl) isValidStatusTransition(from, to int32) bool {
	validTransitions := map[int32][]int32{
		models.OrderStatusPending:   {models.OrderStatusPaid, models.OrderStatusCancelled},
		models.OrderStatusPaid:      {models.OrderStatusShipped, models.OrderStatusCancelled},
		models.OrderStatusShipped:   {models.OrderStatusCompleted},
		models.OrderStatusCompleted: {},
		models.OrderStatusCancelled: {},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
