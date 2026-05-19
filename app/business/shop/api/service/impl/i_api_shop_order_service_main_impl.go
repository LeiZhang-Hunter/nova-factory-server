package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	erpordermodels "nova-factory-server/app/business/erp/sale/salemodels"
	models2 "nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	orderConstant "nova-factory-server/app/constant/order"
	"strconv"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/utils/fileUtils"
	order2 "nova-factory-server/app/utils/order"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	orderCachePrefix      = "shop:app:order:cache:"
	orderCreateLockPrefix = "shop:app:order:create:"
	orderCacheTTL         = 10 * time.Minute
	orderCreateLockTTL    = 15 * time.Second
	defaultDeliveryType   = "express"
)

// Cache 固化本次下单商品快照。
func (s *IApiShopOrderServiceImpl) Cache(c *gin.Context, userID int64, req *models.OrderCacheReq) (*models.OrderCacheResp, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	items, err := s.buildOrderCacheItems(c, userID, req)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}

	if _, err := s.userDao.GetByUserID(c, userID); err != nil {
		return nil, errors.New("读取用户信息失败")
	}

	orderKey := order2.GenerateOrderNo()
	cacheData := &models.OrderCacheData{
		OrderKey:       orderKey,
		UserID:         userID,
		Items:          items,
		DeliveryType:   defaultDeliveryType,
		GoodsAmount:    sumCacheItems(items),
		FreightAmount:  0,
		DiscountAmount: 0,
		PayAmount:      sumCacheItems(items),
	}

	if err := s.saveOrderCache(c, cacheData); err != nil {
		return nil, err
	}

	return &models.OrderCacheResp{
		OrderKey:      orderKey,
		Items:         items,
		ExpireSeconds: int64(orderCacheTTL / time.Second),
	}, nil
}

// Confirm 试算当前预订单并返回确认页数据。
func (s *IApiShopOrderServiceImpl) Confirm(c *gin.Context, userID int64, req *models.OrderConfirmReq) (*models.OrderConfirmResp, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if req == nil || strings.TrimSpace(req.OrderKey) == "" {
		return nil, errors.New("orderKey不能为空")
	}

	cacheData, err := s.getOrderCache(c, req.OrderKey)
	if err != nil {
		return nil, err
	}
	if cacheData.UserID != userID {
		return nil, errors.New("无权操作该预订单")
	}

	address, err := s.resolveConfirmAddress(c, userID, req.AddressID)
	if err != nil {
		return nil, err
	}
	if address != nil {
		cacheData.AddressID = address.ID
	}
	if strings.TrimSpace(req.DeliveryType) != "" {
		cacheData.DeliveryType = strings.TrimSpace(req.DeliveryType)
	}
	if err := s.fillOrderItemsStock(c, cacheData.Items); err != nil {
		return nil, err
	}

	s.recalculateOrderAmounts(c, userID, cacheData)
	if err := s.saveOrderCache(c, cacheData); err != nil {
		return nil, err
	}

	return &models.OrderConfirmResp{
		OrderKey:       cacheData.OrderKey,
		Address:        address,
		Items:          cacheData.Items,
		GoodsAmount:    cacheData.GoodsAmount,
		FreightAmount:  cacheData.FreightAmount,
		DiscountAmount: cacheData.DiscountAmount,
		PayAmount:      cacheData.PayAmount,
		DeliveryType:   cacheData.DeliveryType,
		ExpireSeconds:  int64(orderCacheTTL / time.Second),
	}, nil
}

// getItemAvailableStock 读取商品当前可用库存，活动商品按活动库存与 SKU 库存的较小值返回。
func (s *IApiShopOrderServiceImpl) getItemAvailableStock(c *gin.Context, item *models.OrderCacheItem) (int64, error) {
	if item == nil || item.SkuID <= 0 {
		return 0, nil
	}

	sku, err := s.skuDao.GetByID(c, item.SkuID)
	if err != nil {
		return 0, errors.New("读取商品库存失败")
	}
	if sku == nil {
		return 0, nil
	}

	availableStock := sku.Quantity
	if item.SeckillInfo != nil {
		availableStock = minInt64(availableStock, item.SeckillInfo.Stock)
	}
	if item.CombinationInfo != nil {
		availableStock = minInt64(availableStock, item.CombinationInfo.Stock)
	}
	if availableStock < 0 {
		return 0, nil
	}
	return availableStock, nil
}

// fillOrderItemsStock 填充确认单商品的实时库存，并标记是否库存不足。
func (s *IApiShopOrderServiceImpl) fillOrderItemsStock(c *gin.Context, items []*models.OrderCacheItem) error {
	skuMap, err := s.loadOrderItemSkuMap(c, items)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item == nil {
			continue
		}

		availableStock := s.calcItemAvailableStock(item, skuMap[item.SkuID])
		item.AvailableStock = availableStock
		item.StockInsufficient = item.Quantity > availableStock
	}
	return nil
}

// validateCreateItemsStock 在正式扣减库存前统一校验商品库存，并返回明确的不足提示。
func (s *IApiShopOrderServiceImpl) validateCreateItemsStock(c *gin.Context, items []*models.OrderCacheItem) error {
	if err := s.fillOrderItemsStock(c, items); err != nil {
		return err
	}

	insufficientItems := make([]string, 0)
	for _, item := range items {
		if item == nil || !item.StockInsufficient {
			continue
		}
		insufficientItems = append(insufficientItems, s.buildStockInsufficientMessage(item))
	}
	if len(insufficientItems) > 0 {
		return errors.New("下单失败，库存不足: " + strings.Join(insufficientItems, "；"))
	}
	return nil
}

// loadOrderItemSkuMap 批量读取订单商品对应的 SKU 数据。
func (s *IApiShopOrderServiceImpl) loadOrderItemSkuMap(c *gin.Context, items []*models.OrderCacheItem) (map[int64]*shopmodels.GoodsSku, error) {
	skuIDs := collectOrderItemSkuIDs(items)
	if len(skuIDs) == 0 {
		return map[int64]*shopmodels.GoodsSku{}, nil
	}

	skuList, err := s.skuDao.ListByIDs(c, skuIDs)
	if err != nil {
		return nil, errors.New("读取商品库存失败")
	}

	skuMap := make(map[int64]*shopmodels.GoodsSku, len(skuList))
	for _, sku := range skuList {
		if sku == nil {
			continue
		}
		skuMap[int64(sku.ID)] = sku
	}
	return skuMap, nil
}

// calcItemAvailableStock 计算单个商品当前可用库存，活动商品按活动库存与 SKU 库存的较小值返回。
func (s *IApiShopOrderServiceImpl) calcItemAvailableStock(item *models.OrderCacheItem, sku *shopmodels.GoodsSku) int64 {
	if item == nil || sku == nil {
		return 0
	}

	availableStock := sku.Quantity
	if item.SeckillInfo != nil {
		availableStock = minInt64(availableStock, item.SeckillInfo.Stock)
	}
	if item.CombinationInfo != nil {
		availableStock = minInt64(availableStock, item.CombinationInfo.Stock)
	}
	if availableStock < 0 {
		return 0
	}
	return availableStock
}

// buildStockInsufficientMessage 组装单个商品的库存不足提示。
func (s *IApiShopOrderServiceImpl) buildStockInsufficientMessage(item *models.OrderCacheItem) string {
	if item == nil {
		return ""
	}

	name := strings.TrimSpace(item.GoodsName)
	if skuName := strings.TrimSpace(item.SkuName); skuName != "" {
		name = strings.TrimSpace(name + " " + skuName)
	}
	if name == "" {
		name = fmt.Sprintf("SKU:%d", item.SkuID)
	}

	return fmt.Sprintf("%s(购买%d，剩余%d)", name, item.Quantity, item.AvailableStock)
}

// Create 正式创建订单，读取预订单缓存并落库。
func (s *IApiShopOrderServiceImpl) Create(c *gin.Context, userID int64, req *models.OrderCreateReq) (*models.Order, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if req == nil || strings.TrimSpace(req.OrderKey) == "" {
		return nil, errors.New("orderKey不能为空")
	}

	lockKey := orderCreateLockPrefix + req.OrderKey
	if !s.cache.SetNX(context.Background(), lockKey, "1", orderCreateLockTTL) {
		return nil, errors.New("订单正在处理中，请勿重复提交")
	}
	defer s.cache.Del(context.Background(), lockKey)

	cacheData, err := s.getOrderCache(c, req.OrderKey)
	if err != nil {
		return nil, err
	}
	if cacheData.UserID != userID {
		return nil, errors.New("无权操作该预订单")
	}
	if len(cacheData.Items) == 0 {
		return nil, errors.New("预订单商品不能为空")
	}

	address, err := s.resolveConfirmAddress(c, userID, cacheData.AddressID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("请先选择收货地址")
	}

	if err := s.validateCreateItemsStock(c, cacheData.Items); err != nil {
		return nil, err
	}

	s.recalculateOrderAmounts(c, userID, cacheData)

	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}

	orderNo := order2.GenerateOrderNo()
	if s.db == nil {
		return nil, errors.New("数据库连接不存在")
	}
	err = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		c.Set("db", tx)
		for _, item := range cacheData.Items {
			if item == nil {
				continue
			}
			if err := s.skuDao.DeductStock(c, item.SkuID, item.Quantity); err != nil {
				return fmt.Errorf("扣减库存失败: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	createdOrder, err := s.orderDao.Set(c, s.buildERPOrderSet(orderNo, shopUser, address, cacheData, req))
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}

	s.cache.Del(context.Background(), s.buildOrderCacheKey(req.OrderKey))
	return s.toShopOrder(createdOrder), nil
}

// GetByID 获取订单详情，包含商品明细。
func (s *IApiShopOrderServiceImpl) GetByID(c *gin.Context, id int64) (*models.OrderVO, error) {
	if id == 0 {
		return nil, errors.New("订单ID不能为空")
	}

	order, err := s.orderDao.GetByID(c, uint64(id))
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}
	return s.toShopOrderVO(order), nil
}

// List 获取当前用户的订单列表。
func (s *IApiShopOrderServiceImpl) List(c *gin.Context, userID int64, query *models.OrderQuery) (*models.OrderListData, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if query == nil {
		query = &models.OrderQuery{}
	}

	query.UserID = userID

	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	list, err := s.listERPOrders(c, shopUser, query)
	if err != nil {
		return nil, err
	}
	for _, order := range list.Rows {
		for _, item := range order.Items {
			item.ImageURL = fileUtils.BuildAbsoluteURL(c, item.ImageURL)
		}

	}
	return list, nil
}

// UpdateStatus 更新订单状态，验证状态流转合法性。
func (s *IApiShopOrderServiceImpl) UpdateStatus(c *gin.Context, userID int64, req *models.OrderStatusReq) error {
	if userID == 0 {
		return errors.New("用户未登录")
	}
	if req.ID == 0 {
		return errors.New("订单ID不能为空")
	}

	// 验证用户权限
	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return errors.New("商城用户不存在")
	}
	order, err := s.orderDao.GetByID(c, uint64(req.ID))
	if err != nil {
		return errors.New("订单不存在")
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	if !s.isOrderOwnedByUser(order, shopUser) {
		return errors.New("无权操作此订单")
	}

	// 验证状态流转
	currentStatus := s.erpStatusToShopStatus(order.Status)
	if !s.isValidStatusTransition(currentStatus, req.Status) {
		return errors.New("非法的状态流转")
	}

	rowsAffected, err := s.updateERPOrderStatus(c, req.ID, shopUser, req.Status)
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// Cancel 取消订单，仅允许对待支付的订单进行取消。
func (s *IApiShopOrderServiceImpl) Cancel(c *gin.Context, userID int64, id int64, reason string) error {
	if userID == 0 {
		return errors.New("用户未登录")
	}
	if id == 0 {
		return errors.New("订单ID不能为空")
	}

	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return errors.New("商城用户不存在")
	}
	order, err := s.orderDao.GetByID(c, uint64(id))
	if err != nil {
		return errors.New("订单不存在")
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	if !s.isOrderOwnedByUser(order, shopUser) {
		return errors.New("无权操作此订单")
	}

	if s.erpStatusToShopStatus(order.Status) != orderConstant.OrderStatusPending {
		return errors.New("只能取消待支付的订单")
	}

	rowsAffected, err := s.cancelERPOrder(c, id, shopUser, reason)
	if err != nil {
		return fmt.Errorf("取消订单失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// ConfirmReceive 确认收货，将已发货订单标记为已完成。
func (s *IApiShopOrderServiceImpl) ConfirmReceive(c *gin.Context, userID int64, id int64) error {
	if userID == 0 {
		return errors.New("用户未登录")
	}
	if id == 0 {
		return errors.New("订单ID不能为空")
	}

	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return errors.New("商城用户不存在")
	}
	order, err := s.orderDao.GetByID(c, uint64(id))
	if err != nil {
		return errors.New("订单不存在")
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	if !s.isOrderOwnedByUser(order, shopUser) {
		return errors.New("无权操作此订单")
	}

	if s.erpStatusToShopStatus(order.Status) != orderConstant.OrderStatusShipped {
		return errors.New("只能确认已发货的订单")
	}

	rowsAffected, err := s.updateERPOrderStatus(c, id, shopUser, orderConstant.OrderStatusCompleted)
	if err != nil {
		return fmt.Errorf("确认收货失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// GetStatistics 获取当前用户各状态订单数量统计。
func (s *IApiShopOrderServiceImpl) GetStatistics(c *gin.Context, userID int64) (*models.OrderStatistics, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	return s.getERPOrderStatistics(c, shopUser)
}

// isValidStatusTransition 验证订单状态流转是否合法。
func (s *IApiShopOrderServiceImpl) isValidStatusTransition(from, to int32) bool {
	validTransitions := map[int32][]int32{
		orderConstant.OrderStatusPending:   {orderConstant.OrderStatusPaid, orderConstant.OrderStatusCancelled},
		orderConstant.OrderStatusPaid:      {orderConstant.OrderStatusShipped, orderConstant.OrderStatusCancelled},
		orderConstant.OrderStatusShipped:   {orderConstant.OrderStatusCompleted},
		orderConstant.OrderStatusCompleted: {},
		orderConstant.OrderStatusCancelled: {},
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

func (s *IApiShopOrderServiceImpl) buildOrderCacheKey(orderKey string) string {
	return orderCachePrefix + orderKey
}

func (s *IApiShopOrderServiceImpl) saveOrderCache(c *gin.Context, data *models.OrderCacheData) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	s.cache.Set(context.Background(), s.buildOrderCacheKey(data.OrderKey), string(body), orderCacheTTL)
	return nil
}

func (s *IApiShopOrderServiceImpl) getOrderCache(c *gin.Context, orderKey string) (*models.OrderCacheData, error) {
	val, err := s.cache.Get(context.Background(), s.buildOrderCacheKey(orderKey))
	if err != nil {
		return nil, errors.New("预订单已失效，请重新确认商品")
	}
	var data models.OrderCacheData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, errors.New("预订单数据异常")
	}
	return &data, nil
}

// buildOrderCacheItems 构建预下单商品快照，优先处理购物车选中商品。
func (s *IApiShopOrderServiceImpl) buildOrderCacheItems(c *gin.Context, userID int64, req *models.OrderCacheReq) ([]*models.OrderCacheItem, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	if len(req.CartId) > 0 {
		return s.buildCartCacheItems(c, userID, req.CartId)
	}

	item, err := s.buildCacheItems(c, userID, req)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("订单商品不能为空")
	}
	return []*models.OrderCacheItem{item}, nil
}

func (s *IApiShopOrderServiceImpl) buildCacheItems(c *gin.Context, userID int64, req *models.OrderCacheReq) (*models.OrderCacheItem, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("商品数量必须大于0")
	}

	// 秒杀商品下单时通过活动商品ID读取活动价和所属商品信息。
	if req.SecKillId > 0 {
		item, err := s.buildSeckillCacheItem(c, req.SecKillId, req.SkuID, req.Quantity)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	if req.CombinationId > 0 {
		item, err := s.buildCombinationCacheItem(c, req.CombinationId, req.SkuID, req.Quantity)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	if req.SkuID > 0 {
		item, err := s.buildDirectCacheItem(c, userID, req.SkuID, req.Quantity)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	return nil, errors.New("请选择需要下单的商品")
}

// buildCartCacheItems 根据购物车ID列表构建预下单商品快照。
func (s *IApiShopOrderServiceImpl) buildCartCacheItems(c *gin.Context, userID int64, cartIDs []string) ([]*models.OrderCacheItem, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if len(cartIDs) == 0 {
		return nil, errors.New("购物车ID不能为空")
	}

	idList, err := s.parseCartIDs(cartIDs)
	if err != nil {
		return nil, err
	}
	cartList, err := s.cartDao.ListByIDs(c, userID, idList)
	if err != nil {
		return nil, errors.New("读取购物车商品失败")
	}
	if len(cartList) == 0 {
		return nil, errors.New("购物车商品不存在")
	}

	cartMap := make(map[int64]*models.CartDto, len(cartList))
	for _, cart := range cartList {
		if cart == nil {
			continue
		}
		cartMap[cart.ID] = cart
	}

	items := make([]*models.OrderCacheItem, 0, len(idList))
	for _, id := range idList {
		cartInfo, ok := cartMap[id]
		if !ok {
			return nil, errors.New("部分购物车商品不存在")
		}
		item, err := s.buildSingleCacheItem(c, userID, cartInfo.GoodsID, cartInfo.SkuID, cartInfo.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// parseCartIDs 解析购物车ID字符串列表并去重。
func (s *IApiShopOrderServiceImpl) parseCartIDs(cartIDs []string) ([]int64, error) {
	result := make([]int64, 0, len(cartIDs))
	seen := make(map[int64]struct{}, len(cartIDs))
	for _, cartID := range cartIDs {
		value := strings.TrimSpace(cartID)
		if value == "" {
			continue
		}
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return nil, errors.New("购物车ID格式错误")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	if len(result) == 0 {
		return nil, errors.New("购物车ID不能为空")
	}
	return result, nil
}

// buildSeckillCacheItem 构建秒杀商品预下单快照。
func (s *IApiShopOrderServiceImpl) buildSeckillCacheItem(c *gin.Context, secKillID int64, skuID int64, quantity int64) (*models.OrderCacheItem, error) {
	if skuID <= 0 {
		return nil, errors.New("商品规格不能为空")
	}
	seckill, err := s.seckillDao.GetByID(c, secKillID)
	if err != nil {
		return nil, errors.New("读取秒杀商品失败")
	}
	if seckill == nil {
		return nil, errors.New("秒杀商品不存在")
	}
	if seckill.IsShow != 1 || seckill.Status != 1 {
		return nil, errors.New("秒杀商品已下架")
	}
	if seckill.OnceNum > 0 && quantity > int64(seckill.OnceNum) {
		return nil, errors.New("超过秒杀单次限购数量")
	}

	goods, sku, err := s.loadGoodsAndSkuByGoodsID(c, seckill.ProductID, skuID, quantity)
	if err != nil {
		return nil, err
	}
	goodsName := strings.TrimSpace(seckill.Title)
	if goodsName == "" {
		goodsName = goods.GoodsName
	}
	if err := s.validateCacheItemStock(goodsName, sku.SkuName, quantity, minInt64(seckill.Stock, sku.Quantity), "秒杀商品"); err != nil {
		return nil, err
	}
	imageURL := strings.TrimSpace(seckill.Image)
	item := s.assembleCacheItem(c, goods, sku, quantity, seckill.Price, goodsName, imageURL)
	item.SecKillId = secKillID
	item.SeckillInfo = models2.FromatSeckillMainInfo(seckill)
	return item, nil
}

// buildCombinationCacheItem 构建拼团商品预下单快照。
func (s *IApiShopOrderServiceImpl) buildCombinationCacheItem(c *gin.Context, combinationID int64, skuID int64, quantity int64) (*models.OrderCacheItem, error) {
	if skuID <= 0 {
		return nil, errors.New("商品规格不能为空")
	}
	combination, err := s.combDao.GetByID(c, combinationID)
	if err != nil {
		return nil, errors.New("读取拼团商品失败")
	}
	if combination == nil {
		return nil, errors.New("拼团商品不存在")
	}
	if combination.IsShow != 1 {
		return nil, errors.New("拼团商品已下架")
	}
	if combination.OnceNum > 0 && quantity > combination.OnceNum {
		return nil, errors.New("超过拼团单次限购数量")
	}

	goods, sku, err := s.loadGoodsAndSkuByGoodsCode(c, strings.TrimSpace(combination.ProductID), skuID, quantity)
	if err != nil {
		return nil, err
	}
	goodsName := strings.TrimSpace(combination.Title)
	if goodsName == "" {
		goodsName = goods.GoodsName
	}
	if err := s.validateCacheItemStock(goodsName, sku.SkuName, quantity, minInt64(combination.Stock, sku.Quantity), "拼团商品"); err != nil {
		return nil, err
	}
	imageURL := strings.TrimSpace(combination.Image)
	item := s.assembleCacheItem(c, goods, sku, quantity, combination.Price, goodsName, imageURL)
	item.CombinationId = combinationID
	item.CombinationInfo = models2.FormatCombinationMainInfo(combination)
	return item, nil
}

// buildDirectCacheItem 构建普通立即购买商品预下单快照。
func (s *IApiShopOrderServiceImpl) buildDirectCacheItem(c *gin.Context, userID int64, skuID int64, quantity int64) (*models.OrderCacheItem, error) {
	if skuID <= 0 {
		return nil, errors.New("商品规格不能为空")
	}
	sku, err := s.skuDao.GetByID(c, skuID)
	if err != nil {
		return nil, errors.New("读取商品规格失败")
	}
	if sku == nil {
		return nil, errors.New("商品规格不存在")
	}
	goods, err := s.goodsDao.GetByGoodsID(c, sku.GoodsID)
	if err != nil {
		return nil, errors.New("读取商品信息失败")
	}
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	if err := s.validateCacheItemStock(goods.GoodsName, sku.SkuName, quantity, sku.Quantity, ""); err != nil {
		return nil, err
	}
	price := s.applyGoodsDiscount(c, userID, goods, sku.SkuID, sku.RetailPrice)
	return s.assembleCacheItem(c, goods, sku, quantity, price, goods.GoodsName, ""), nil
}

// loadGoodsAndSkuByGoodsID 按商品主键和规格ID加载商品快照所需数据。
func (s *IApiShopOrderServiceImpl) loadGoodsAndSkuByGoodsID(c *gin.Context, goodsID int64, skuID int64, quantity int64) (*models.Goods, *shopmodels.GoodsSku, error) {
	goods, err := s.goodsDao.GetByID(c, goodsID)
	if err != nil {
		return nil, nil, errors.New("读取商品信息失败")
	}
	if goods == nil {
		return nil, nil, errors.New("商品不存在")
	}
	sku, err := s.skuDao.GetByID(c, skuID)
	if err != nil {
		return nil, nil, errors.New("读取商品规格失败")
	}
	if sku == nil {
		return nil, nil, errors.New("商品规格不存在")
	}
	if sku.GoodsID != goods.GoodsID {
		return nil, nil, errors.New("商品规格与活动商品不匹配")
	}
	if goods.IsOnSale != 1 {
		return nil, nil, errors.New("商品已下架")
	}
	if goods.Quantity <= 0 || sku.Quantity <= 0 {
		return nil, nil, s.newStockInsufficientError(goods.GoodsName, sku.SkuName, quantity, 0, "")
	}
	return goods, sku, nil
}

// loadGoodsAndSkuByGoodsCode 按商品业务ID和规格ID加载商品快照所需数据。
func (s *IApiShopOrderServiceImpl) loadGoodsAndSkuByGoodsCode(c *gin.Context, goodsCode string, skuID int64, quantity int64) (*models.Goods, *shopmodels.GoodsSku, error) {
	if goodsCode == "" {
		return nil, nil, errors.New("活动商品数据异常")
	}
	goods, err := s.goodsDao.GetByGoodsID(c, goodsCode)
	if err != nil {
		return nil, nil, errors.New("读取商品信息失败")
	}
	if goods == nil {
		return nil, nil, errors.New("商品不存在")
	}
	sku, err := s.skuDao.GetByID(c, skuID)
	if err != nil {
		return nil, nil, errors.New("读取商品规格失败")
	}
	if sku == nil {
		return nil, nil, errors.New("商品规格不存在")
	}
	if sku.GoodsID != goods.GoodsID {
		return nil, nil, errors.New("商品规格与活动商品不匹配")
	}
	if goods.IsOnSale != 1 {
		return nil, nil, errors.New("商品已下架")
	}
	if goods.Quantity <= 0 || sku.Quantity <= 0 {
		return nil, nil, s.newStockInsufficientError(goods.GoodsName, sku.SkuName, quantity, 0, "")
	}
	return goods, sku, nil
}

// assembleCacheItem 统一组装预下单缓存项，优先保留活动价和活动展示信息。
func (s *IApiShopOrderServiceImpl) assembleCacheItem(
	c *gin.Context,
	goods *models.Goods,
	sku *shopmodels.GoodsSku,
	quantity int64,
	price float64,
	goodsName string,
	imageURL string,
) *models.OrderCacheItem {
	finalGoodsName := strings.TrimSpace(goodsName)
	if finalGoodsName == "" {
		finalGoodsName = goods.GoodsName
	}
	finalImageURL := strings.TrimSpace(imageURL)

	if finalImageURL == "" {
		finalImageURL = strings.TrimSpace(goods.ImageURL)
	}

	if finalImageURL == "" {
		finalImageURL = strings.TrimSpace(sku.ImageURL)
	}
	return &models.OrderCacheItem{
		GoodsID:     goods.ID,
		SkuID:       int64(sku.ID),
		GoodsName:   finalGoodsName,
		SkuName:     sku.SkuName,
		ImageURL:    fileUtils.BuildAbsoluteURL(c, finalImageURL),
		Price:       price,
		Quantity:    quantity,
		TotalAmount: price * float64(quantity),
	}
}

func (s *IApiShopOrderServiceImpl) buildCacheItemsFromRequest(c *gin.Context, userID int64, reqs []*models.OrderCacheItemReq) ([]*models.OrderCacheItem, error) {
	if len(reqs) == 0 {
		return nil, errors.New("订单商品不能为空")
	}
	items := make([]*models.OrderCacheItem, 0, len(reqs))
	for _, req := range reqs {
		if req == nil {
			continue
		}
		if req.Quantity <= 0 {
			return nil, errors.New("商品数量必须大于0")
		}
		item, err := s.buildSingleCacheItem(c, userID, req.GoodsID, req.SkuID, req.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *IApiShopOrderServiceImpl) buildSingleCacheItem(c *gin.Context, userID int64, goodsID int64, skuID int64, quantity int64) (*models.OrderCacheItem, error) {
	goods, err := s.goodsDao.GetByID(c, goodsID)
	if err != nil {
		return nil, errors.New("读取商品信息失败")
	}
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	sku, err := s.skuDao.GetByID(c, skuID)
	if err != nil {
		return nil, errors.New("读取商品规格失败")
	}
	if sku == nil {
		return nil, errors.New("商品规格不存在")
	}
	if err := s.validateCacheItemStock(goods.GoodsName, sku.SkuName, quantity, sku.Quantity, ""); err != nil {
		return nil, err
	}
	price := s.applyGoodsDiscount(c, userID, goods, sku.SkuID, sku.RetailPrice)
	item := &models.OrderCacheItem{
		GoodsID:     goodsID,
		SkuID:       skuID,
		GoodsName:   goods.GoodsName,
		SkuName:     sku.SkuName,
		ImageURL:    fileUtils.BuildAbsoluteURL(c, sku.ImageURL),
		Price:       price,
		Quantity:    quantity,
		TotalAmount: price * float64(quantity),
	}
	return item, nil
}

// validateCacheItemStock 校验缓存商品的库存是否满足本次购买数量。
func (s *IApiShopOrderServiceImpl) validateCacheItemStock(goodsName string, skuName string, quantity int64, availableStock int64, scene string) error {
	if quantity <= availableStock {
		return nil
	}
	return s.newStockInsufficientError(goodsName, skuName, quantity, availableStock, scene)
}

// newStockInsufficientError 组装带商品信息的库存不足错误提示。
func (s *IApiShopOrderServiceImpl) newStockInsufficientError(goodsName string, skuName string, quantity int64, availableStock int64, scene string) error {
	message := s.buildStockInsufficientDetail(goodsName, skuName, quantity, availableStock)
	if scene != "" {
		return errors.New(scene + "库存不足: " + message)
	}
	return errors.New("库存不足: " + message)
}

// buildStockInsufficientDetail 组装库存不足明细，包含商品名、规格、购买数和剩余库存。
func (s *IApiShopOrderServiceImpl) buildStockInsufficientDetail(goodsName string, skuName string, quantity int64, availableStock int64) string {
	name := strings.TrimSpace(goodsName)
	if currentSkuName := strings.TrimSpace(skuName); currentSkuName != "" {
		name = strings.TrimSpace(name + " " + currentSkuName)
	}
	if name == "" {
		name = "商品"
	}
	return fmt.Sprintf("%s(购买%d，剩余%d)", name, quantity, availableStock)
}

func (s *IApiShopOrderServiceImpl) resolveConfirmAddress(c *gin.Context, userID int64, addressID int64) (*models.ShopUserAddressApp, error) {
	if addressID > 0 {
		address, err := s.addressDao.GetByID(c, addressID)
		if err != nil {
			return nil, err
		}
		if address == nil || address.UserID != userID {
			return nil, errors.New("收货地址不存在")
		}
		return address, nil
	}
	list, err := s.addressDao.List(c, userID)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Rows) == 0 {
		return nil, nil
	}
	for _, address := range list.Rows {
		if address != nil && address.IsDefault == 1 {
			return address, nil
		}
	}
	return list.Rows[0], nil
}

func (s *IApiShopOrderServiceImpl) recalculateOrderAmounts(c *gin.Context, userID int64, data *models.OrderCacheData) {
	data.GoodsAmount = sumCacheItems(data.Items)
	if data.DeliveryType == "" {
		data.DeliveryType = defaultDeliveryType
	}
	data.FreightAmount = 0

	// 折扣已在构建缓存时应用，这里不再重复计算折扣
	// discountAmount 保持为 0（cache 时已折扣）
	data.DiscountAmount = 0
	data.PayAmount = data.GoodsAmount + data.FreightAmount
	if data.PayAmount < 0 {
		data.PayAmount = 0
	}
}

// applyGoodsDiscount 对单个商品应用用户折扣，返回折扣后价格
func (s *IApiShopOrderServiceImpl) applyGoodsDiscount(c *gin.Context, userID int64, goods *models.Goods, skuID string, price float64) float64 {
	if s.discountService == nil || userID == 0 || goods == nil || price <= 0 {
		return price
	}
	discountedPrice, hasDiscount := s.discountService.CalculateDiscountPrice(
		c, userID, goods.GoodsID, skuID, strconv.FormatInt(goods.ShopCategoryId, 10), price)
	if hasDiscount {
		return discountedPrice
	}
	return price
}

func sumCacheItems(items []*models.OrderCacheItem) float64 {
	var total float64
	for _, item := range items {
		if item != nil {
			total += item.TotalAmount
		}
	}
	return total
}

// minInt64 返回两个 int64 中的较小值。
func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// collectOrderItemSkuIDs 收集订单商品中的 SKU ID 并去重。
func collectOrderItemSkuIDs(items []*models.OrderCacheItem) []int64 {
	result := make([]int64, 0, len(items))
	seen := make(map[int64]struct{}, len(items))
	for _, item := range items {
		if item == nil || item.SkuID <= 0 {
			continue
		}
		if _, ok := seen[item.SkuID]; ok {
			continue
		}
		seen[item.SkuID] = struct{}{}
		result = append(result, item.SkuID)
	}
	return result
}

// buildERPOrderSet 将商城订单请求转换为 ERP 订单保存参数。
func (s *IApiShopOrderServiceImpl) buildERPOrderSet(
	orderNo string,
	shopUser *models.User,
	address *models.ShopUserAddressApp,
	cacheData *models.OrderCacheData,
	req *models.OrderCreateReq,
) *erpordermodels.OrderSet {
	details := make([]*erpordermodels.OrderDetailSet, 0, len(cacheData.Items))
	for index, item := range cacheData.Items {
		if item == nil {
			continue
		}
		details = append(details, &erpordermodels.OrderDetailSet{
			OID:            fmt.Sprintf("%s-%d", orderNo, index+1),
			EShopGoodsID:   fmt.Sprintf("%d", item.GoodsID),
			EShopGoodsName: item.GoodsName,
			EShopSkuID:     fmt.Sprintf("%d", item.SkuID),
			EShopSkuName:   item.SkuName,
			NumIID:         item.GoodsID,
			SkuID:          item.SkuID,
			Num:            float64(item.Quantity),
			Payment:        item.TotalAmount,
			PicPath:        item.ImageURL,
		})
	}

	return &erpordermodels.OrderSet{
		Tid:                  orderNo,
		BuyerNick:            s.buildOrderBuyerNick(shopUser),
		BuyerMessage:         strings.TrimSpace(req.Remark),
		SellerMemo:           strings.TrimSpace(req.Remark),
		Total:                cacheData.GoodsAmount,
		Privilege:            cacheData.DiscountAmount,
		PostFee:              cacheData.FreightAmount,
		ReceiverName:         address.ReceiverName,
		ReceiverProvince:     address.ProvinceCode,
		ReceiverProvinceName: address.ProvinceName,
		ReceiverCity:         address.CityCode,
		ReceiverCityName:     address.CityName,
		ReceiverDistrict:     address.DistrictCode,
		ReceiverDistrictName: address.DistrictName,
		ReceiverAddress:      address.DetailAddress,
		ReceiverPhone:        address.ReceiverMobile,
		ReceiverMobile:       address.ReceiverMobile,
		Status:               s.shopStatusToERPStatus(orderConstant.OrderStatusPending),
		OrderType:            "shop",
		Details:              details,
		Accounts: []*erpordermodels.OrderAccountSet{
			{
				FinanceCode: "PAY_AMOUNT",
				Total:       cacheData.PayAmount,
			},
		},
	}
}

// toShopOrder 将 ERP 订单模型转换为商城订单模型。
func (s *IApiShopOrderServiceImpl) toShopOrder(order *erpordermodels.Order) *models.Order {
	if order == nil {
		return nil
	}
	return &models.Order{
		ID:                    int64(order.ID),
		OrderNo:               order.Tid,
		TotalAmount:           order.Total,
		PayAmount:             order.Total - order.Privilege + order.PostFee,
		FreightAmount:         order.PostFee,
		DiscountAmount:        order.Privilege,
		Status:                s.erpStatusToShopStatus(order.Status),
		PayTime:               order.PayTime,
		ReceiverName:          order.ReceiverName,
		ReceiverPhone:         s.firstNonEmpty(order.ReceiverMobile, order.ReceiverPhone),
		ReceiverProvince:      s.firstNonEmpty(order.ReceiverProvinceName, order.ReceiverProvince),
		ReceiverCity:          s.firstNonEmpty(order.ReceiverCityName, order.ReceiverCity),
		ReceiverDistrict:      s.firstNonEmpty(order.ReceiverDistrictName, order.ReceiverDistrict),
		ReceiverDetailAddress: order.ReceiverAddress,
		Remark:                s.firstNonEmpty(order.SellerMemo, order.BuyerMessage),
		DeptID:                order.DeptID,
		State:                 order.State,
		BaseEntity:            order.BaseEntity,
	}
}

// toShopOrderVO 将 ERP 订单详情转换为商城订单视图。
func (s *IApiShopOrderServiceImpl) toShopOrderVO(order *erpordermodels.Order) *models.OrderVO {
	if order == nil {
		return nil
	}
	items := make([]*models.OrderItem, 0, len(order.Details))
	for _, detail := range order.Details {
		if detail == nil {
			continue
		}
		items = append(items, s.toShopOrderItem(detail))
	}
	return &models.OrderVO{
		Order: *s.toShopOrder(order),
		Items: items,
	}
}

// toShopOrderItem 将 ERP 订单明细转换为商城订单商品明细。
func (s *IApiShopOrderServiceImpl) toShopOrderItem(detail *erpordermodels.OrderDetail) *models.OrderItem {
	if detail == nil {
		return nil
	}
	price := detail.Payment
	if detail.Num > 0 {
		price = detail.Payment / detail.Num
	}
	return &models.OrderItem{
		ID:          int64(detail.ID),
		OrderID:     int64(detail.OrderID),
		OrderNo:     detail.Tid,
		GoodsID:     s.firstNonEmpty(detail.EShopGoodsID, fmt.Sprintf("%d", detail.NumIID)),
		SkuID:       s.firstNonEmpty(detail.EShopSkuID, fmt.Sprintf("%d", detail.SkuID)),
		GoodsName:   detail.EShopGoodsName,
		SkuName:     detail.EShopSkuName,
		ImageURL:    detail.PicPath,
		Price:       price,
		Quantity:    int32(detail.Num),
		TotalAmount: detail.Payment,
		DeptID:      detail.DeptID,
		State:       detail.State,
		BaseEntity:  detail.BaseEntity,
	}
}

// listERPOrders 查询当前商城用户在 ERP 表中的订单列表。
func (s *IApiShopOrderServiceImpl) listERPOrders(c *gin.Context, shopUser *models.User, query *models.OrderQuery) (*models.OrderListData, error) {
	db := s.erpOrderBaseQuery(c).Where("buyer_nick = ?", s.buildOrderBuyerNick(shopUser))
	if query.Status != nil {
		db = db.Where("status = ?", s.shopStatusToERPStatus(*query.Status))
	}
	if strings.TrimSpace(query.OrderNo) != "" {
		db = db.Where("tid LIKE ?", "%"+strings.TrimSpace(query.OrderNo)+"%")
	}
	page := query.Page
	if page < 1 {
		page = 1
	}
	size := query.Size
	if size < 1 {
		size = 10
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*erpordermodels.Order, 0)
	if err := db.Order("id DESC").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	if err := s.attachERPOrderDetails(c, rows); err != nil {
		return nil, err
	}
	data := &models.OrderListData{
		Rows:  make([]*models.OrderVO, 0, len(rows)),
		Total: total,
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		data.Rows = append(data.Rows, s.toShopOrderVO(row))
	}
	return data, nil
}

// updateERPOrderStatus 更新 ERP 订单状态。
func (s *IApiShopOrderServiceImpl) updateERPOrderStatus(c *gin.Context, id int64, shopUser *models.User, status int32) (int64, error) {
	result := s.erpOrderBaseQuery(c).
		Where("id = ?", id).
		Where("buyer_nick = ?", s.buildOrderBuyerNick(shopUser)).
		Updates(map[string]interface{}{
			"status":      s.shopStatusToERPStatus(status),
			"update_time": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

// cancelERPOrder 将 ERP 订单标记为已取消。
func (s *IApiShopOrderServiceImpl) cancelERPOrder(c *gin.Context, id int64, shopUser *models.User, reason string) (int64, error) {
	result := s.erpOrderBaseQuery(c).
		Where("id = ?", id).
		Where("buyer_nick = ?", s.buildOrderBuyerNick(shopUser)).
		Updates(map[string]interface{}{
			"status":      s.shopStatusToERPStatus(orderConstant.OrderStatusCancelled),
			"seller_memo": strings.TrimSpace(reason),
			"update_time": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

// getERPOrderStatistics 统计当前商城用户在 ERP 表中的订单状态数量。
func (s *IApiShopOrderServiceImpl) getERPOrderStatistics(c *gin.Context, shopUser *models.User) (*models.OrderStatistics, error) {
	stats := &models.OrderStatistics{}
	baseQuery := s.erpOrderBaseQuery(c).Where("buyer_nick = ?", s.buildOrderBuyerNick(shopUser))
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(orderConstant.OrderStatusPending)).
		Count(&stats.PendingPay).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(orderConstant.OrderStatusPaid)).
		Count(&stats.PendingSend).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(orderConstant.OrderStatusShipped)).
		Count(&stats.PendingReceive).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(orderConstant.OrderStatusCompleted)).
		Count(&stats.Completed).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(orderConstant.OrderStatusCancelled)).
		Count(&stats.Cancelled).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

// attachERPOrderDetails 批量挂载 ERP 订单明细。
func (s *IApiShopOrderServiceImpl) attachERPOrderDetails(c *gin.Context, orders []*erpordermodels.Order) error {
	if len(orders) == 0 {
		return nil
	}
	orderIDs := make([]uint64, 0, len(orders))
	for _, order := range orders {
		if order == nil {
			continue
		}
		orderIDs = append(orderIDs, order.ID)
	}
	if len(orderIDs) == 0 {
		return nil
	}
	details := make([]*erpordermodels.OrderDetail, 0)
	if err := s.db.WithContext(c).
		Table("erp_order_detail").
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&details).Error; err != nil {
		return err
	}
	detailMap := make(map[uint64][]*erpordermodels.OrderDetail)
	for _, detail := range details {
		if detail == nil {
			continue
		}
		detailMap[detail.OrderID] = append(detailMap[detail.OrderID], detail)
	}
	for _, order := range orders {
		if order == nil {
			continue
		}
		order.Details = detailMap[order.ID]
	}
	return nil
}

// erpOrderBaseQuery 构建 ERP 订单表基础查询。
func (s *IApiShopOrderServiceImpl) erpOrderBaseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).
		Table("erp_order").
		Where("state = ?", commonStatus.NORMAL)
}

// buildOrderBuyerNick 生成商城订单在 ERP 表中的买家标识。
func (s *IApiShopOrderServiceImpl) buildOrderBuyerNick(shopUser *models.User) string {
	if shopUser == nil {
		return ""
	}
	return fmt.Sprintf("shop-user-%s", shopUser.UserID)
}

// buildOrderOwnerCandidates 构造商城用户在 ERP 订单中的可能归属标识。
func (s *IApiShopOrderServiceImpl) buildOrderOwnerCandidates(shopUser *models.User) []string {
	if shopUser == nil {
		return []string{""}
	}
	candidates := make([]string, 0, 6)
	addCandidate := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		for _, item := range candidates {
			if item == value {
				return
			}
		}
		candidates = append(candidates, value)
	}
	addCandidate(s.buildOrderBuyerNick(shopUser))
	addCandidate(fmt.Sprintf("%d", shopUser.ID))
	addCandidate(shopUser.UserID)
	addCandidate(shopUser.Username)
	addCandidate(shopUser.Nickname)
	addCandidate(shopUser.WechatOpenid)
	if len(candidates) == 0 {
		candidates = append(candidates, "")
	}
	return candidates
}

// isOrderOwnedByUser 校验 ERP 订单是否属于当前商城用户。
func (s *IApiShopOrderServiceImpl) isOrderOwnedByUser(order *erpordermodels.Order, shopUser *models.User) bool {
	if order == nil || shopUser == nil {
		return false
	}
	for _, candidate := range s.buildOrderOwnerCandidates(shopUser) {
		if order.BuyerNick == candidate {
			return true
		}
	}
	return false
}

// shopStatusToERPStatus 将商城订单状态转换为 ERP 状态值。
func (s *IApiShopOrderServiceImpl) shopStatusToERPStatus(status int32) string {
	switch status {
	case orderConstant.OrderStatusPending:
		return orderConstant.ERPStatusNoPay
	case orderConstant.OrderStatusPaid:
		return orderConstant.ERPStatusPayed
	case orderConstant.OrderStatusShipped:
		return orderConstant.ERPStatusSended
	case orderConstant.OrderStatusPartShipped:
		return orderConstant.ERPStatusPartSend
	case orderConstant.OrderStatusCompleted:
		return orderConstant.ERPStatusTradeSuccess
	case orderConstant.OrderStatusCancelled:
		return orderConstant.ERPStatusTradeClosed
	case orderConstant.OrderStatusAftersale:
		return orderConstant.ERPStatusAftersale
	default:
		return orderConstant.ERPStatusNoPay
	}
}

// erpStatusToShopStatus 将 ERP 状态值转换为商城订单状态。
func (s *IApiShopOrderServiceImpl) erpStatusToShopStatus(status string) int32 {
	switch strings.TrimSpace(status) {
	case orderConstant.ERPStatusNoPay:
		return orderConstant.OrderStatusPending
	case orderConstant.ERPStatusPayed:
		return orderConstant.OrderStatusPaid
	case orderConstant.ERPStatusSended:
		return orderConstant.OrderStatusShipped
	case orderConstant.ERPStatusPartSend:
		return orderConstant.OrderStatusPartShipped
	case orderConstant.ERPStatusTradeSuccess:
		return orderConstant.OrderStatusCompleted
	case orderConstant.ERPStatusTradeClosed:
		return orderConstant.OrderStatusCancelled
	case orderConstant.ERPStatusAftersale:
		return orderConstant.OrderStatusAftersale
	default:
		return orderConstant.OrderStatusPending
	}
}

// firstNonEmpty 返回第一个非空字符串。
func (s *IApiShopOrderServiceImpl) firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
