package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	erporderdao "nova-factory-server/app/business/erp/order/orderdao"
	erpordermodels "nova-factory-server/app/business/erp/order/ordermodels"
	activityDao "nova-factory-server/app/business/shop/activity/dao"
	models2 "nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/baizeContext"
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

// IApiShopOrderServiceImpl 提供订单相关的业务实现。
type IApiShopOrderServiceImpl struct {
	cache        cache.Cache
	db           *gorm.DB
	orderDao     erporderdao.IOrderDao
	orderItemDao dao.IApiShopOrderItemDao
	userDao      dao.IApiShopWechatUserDao
	addressDao   dao.IApiShopAddressDao
	seckillDao   activityDao.IShopSeckillDao
	combDao      activityDao.IShopCombinationDao
	goodsDao     dao.IApiShopGoodsDao
	skuDao       dao.IApiShopSkuDao
}

// NewIApiShopOrderServiceImpl 创建订单服务实现。
func NewIApiShopOrderServiceImpl(
	cache cache.Cache,
	db *gorm.DB,
	orderDao erporderdao.IOrderDao,
	orderItemDao dao.IApiShopOrderItemDao,
	userDao dao.IApiShopWechatUserDao,
	addressDao dao.IApiShopAddressDao,
	seckillDao activityDao.IShopSeckillDao,
	combDao activityDao.IShopCombinationDao,
	goodsDao dao.IApiShopGoodsDao,
	skuDao dao.IApiShopSkuDao,
) service.IApiShopOrderService {
	return &IApiShopOrderServiceImpl{
		cache:        cache,
		db:           db,
		orderDao:     orderDao,
		orderItemDao: orderItemDao,
		userDao:      userDao,
		addressDao:   addressDao,
		seckillDao:   seckillDao,
		combDao:      combDao,
		goodsDao:     goodsDao,
		skuDao:       skuDao,
	}
}

// Cache 固化本次下单商品快照。
func (s *IApiShopOrderServiceImpl) Cache(c *gin.Context, userID int64, req *models.OrderCacheReq) (*models.OrderCacheResp, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	item, err := s.buildCacheItems(c, req)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("订单商品不能为空")
	}

	if _, err := s.userDao.GetByID(c, userID); err != nil {
		return nil, errors.New("读取用户信息失败")
	}

	orderKey := order2.GenerateOrderNo()
	cacheData := &models.OrderCacheData{
		OrderKey:       orderKey,
		UserID:         userID,
		Item:           item,
		DeliveryType:   defaultDeliveryType,
		GoodsAmount:    sumCacheItems([]*models.OrderCacheItem{item}),
		FreightAmount:  0,
		DiscountAmount: 0,
		PayAmount:      sumCacheItems([]*models.OrderCacheItem{item}),
	}

	if err := s.saveOrderCache(c, cacheData); err != nil {
		return nil, err
	}

	return &models.OrderCacheResp{
		OrderKey:      orderKey,
		Item:          item,
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

	s.recalculateOrderAmounts(cacheData)
	if err := s.saveOrderCache(c, cacheData); err != nil {
		return nil, err
	}

	return &models.OrderConfirmResp{
		OrderKey:       cacheData.OrderKey,
		Address:        address,
		Item:           cacheData.Item,
		GoodsAmount:    cacheData.GoodsAmount,
		FreightAmount:  cacheData.FreightAmount,
		DiscountAmount: cacheData.DiscountAmount,
		PayAmount:      cacheData.PayAmount,
		DeliveryType:   cacheData.DeliveryType,
		ExpireSeconds:  int64(orderCacheTTL / time.Second),
	}, nil
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
	if cacheData.Item == nil {
		return nil, errors.New("预订单商品不能为空")
	}

	address, err := s.resolveConfirmAddress(c, userID, cacheData.AddressID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("请先选择收货地址")
	}

	s.recalculateOrderAmounts(cacheData)

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
		if err := s.skuDao.DeductStock(c, cacheData.Item.SkuID, cacheData.Item.Quantity); err != nil {
			return fmt.Errorf("扣减库存失败: %v", err)
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

	shopUser, err := s.userDao.GetByID(c, userID)
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
	shopUser, err := s.userDao.GetByID(c, userID)
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

	shopUser, err := s.userDao.GetByID(c, userID)
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

	if s.erpStatusToShopStatus(order.Status) != models.OrderStatusPending {
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

	shopUser, err := s.userDao.GetByID(c, userID)
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

	if s.erpStatusToShopStatus(order.Status) != models.OrderStatusShipped {
		return errors.New("只能确认已发货的订单")
	}

	rowsAffected, err := s.updateERPOrderStatus(c, id, shopUser, models.OrderStatusCompleted)
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
	shopUser, err := s.userDao.GetByID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	return s.getERPOrderStatistics(c, shopUser)
}

// isValidStatusTransition 验证订单状态流转是否合法。
func (s *IApiShopOrderServiceImpl) isValidStatusTransition(from, to int32) bool {
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

func (s *IApiShopOrderServiceImpl) buildCacheItems(c *gin.Context, req *models.OrderCacheReq) (*models.OrderCacheItem, error) {
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
		item, err := s.buildDirectCacheItem(c, req.SkuID, req.Quantity)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	return nil, errors.New("请选择需要下单的商品")
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
	if quantity > seckill.Stock {
		return nil, errors.New("秒杀库存不足")
	}
	if seckill.OnceNum > 0 && quantity > int64(seckill.OnceNum) {
		return nil, errors.New("超过秒杀单次限购数量")
	}

	goods, sku, err := s.loadGoodsAndSkuByGoodsID(c, seckill.ProductID, skuID)
	if err != nil {
		return nil, err
	}
	goodsName := strings.TrimSpace(seckill.Title)
	if goodsName == "" {
		goodsName = goods.GoodsName
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
	if quantity > combination.Stock {
		return nil, errors.New("拼团库存不足")
	}
	if combination.OnceNum > 0 && quantity > combination.OnceNum {
		return nil, errors.New("超过拼团单次限购数量")
	}

	goods, sku, err := s.loadGoodsAndSkuByGoodsCode(c, strings.TrimSpace(combination.ProductID), skuID)
	if err != nil {
		return nil, err
	}
	goodsName := strings.TrimSpace(combination.Title)
	if goodsName == "" {
		goodsName = goods.GoodsName
	}
	imageURL := strings.TrimSpace(combination.Image)
	item := s.assembleCacheItem(c, goods, sku, quantity, combination.Price, goodsName, imageURL)
	item.CombinationId = combinationID
	item.CombinationInfo = models2.FormatCombinationMainInfo(combination)
	return item, nil
}

// buildDirectCacheItem 构建普通立即购买商品预下单快照。
func (s *IApiShopOrderServiceImpl) buildDirectCacheItem(c *gin.Context, skuID int64, quantity int64) (*models.OrderCacheItem, error) {
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
	if quantity > sku.Quantity {
		return nil, errors.New("库存不足")
	}
	return s.assembleCacheItem(c, goods, sku, quantity, sku.RetailPrice, goods.GoodsName, ""), nil
}

// loadGoodsAndSkuByGoodsID 按商品主键和规格ID加载商品快照所需数据。
func (s *IApiShopOrderServiceImpl) loadGoodsAndSkuByGoodsID(c *gin.Context, goodsID int64, skuID int64) (*models.Goods, *shopmodels.GoodsSku, error) {
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
		return nil, nil, errors.New("库存不足")
	}
	return goods, sku, nil
}

// loadGoodsAndSkuByGoodsCode 按商品业务ID和规格ID加载商品快照所需数据。
func (s *IApiShopOrderServiceImpl) loadGoodsAndSkuByGoodsCode(c *gin.Context, goodsCode string, skuID int64) (*models.Goods, *shopmodels.GoodsSku, error) {
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
		return nil, nil, errors.New("库存不足")
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

func (s *IApiShopOrderServiceImpl) buildCacheItemsFromRequest(c *gin.Context, reqs []*models.OrderCacheItemReq) ([]*models.OrderCacheItem, error) {
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
		item, err := s.buildSingleCacheItem(c, req.GoodsID, req.SkuID, req.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *IApiShopOrderServiceImpl) buildSingleCacheItem(c *gin.Context, goodsID int64, skuID int64, quantity int64) (*models.OrderCacheItem, error) {
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
	if quantity > sku.Quantity {
		return nil, errors.New("库存不足")
	}
	item := &models.OrderCacheItem{
		GoodsID:     goodsID,
		SkuID:       skuID,
		GoodsName:   goods.GoodsName,
		SkuName:     sku.SkuName,
		ImageURL:    fileUtils.BuildAbsoluteURL(c, sku.ImageURL),
		Price:       sku.RetailPrice,
		Quantity:    quantity,
		TotalAmount: sku.RetailPrice * float64(quantity),
	}
	return item, nil
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

func (s *IApiShopOrderServiceImpl) recalculateOrderAmounts(data *models.OrderCacheData) {
	data.GoodsAmount = sumCacheItems([]*models.OrderCacheItem{data.Item})
	if data.DeliveryType == "" {
		data.DeliveryType = defaultDeliveryType
	}
	data.FreightAmount = 0
	data.DiscountAmount = 0
	data.PayAmount = data.GoodsAmount + data.FreightAmount - data.DiscountAmount
	if data.PayAmount < 0 {
		data.PayAmount = 0
	}
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

// buildERPOrderSet 将商城订单请求转换为 ERP 订单保存参数。
func (s *IApiShopOrderServiceImpl) buildERPOrderSet(
	orderNo string,
	shopUser *models.User,
	address *models.ShopUserAddressApp,
	cacheData *models.OrderCacheData,
	req *models.OrderCreateReq,
) *erpordermodels.OrderSet {
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
		Status:               s.shopStatusToERPStatus(models.OrderStatusPending),
		OrderType:            "shop",
		Details: []*erpordermodels.OrderDetailSet{
			{
				OID:            fmt.Sprintf("%s-1", orderNo),
				EShopGoodsID:   fmt.Sprintf("%d", cacheData.Item.GoodsID),
				EShopGoodsName: cacheData.Item.GoodsName,
				EShopSkuID:     fmt.Sprintf("%d", cacheData.Item.SkuID),
				EShopSkuName:   cacheData.Item.SkuName,
				NumIID:         cacheData.Item.GoodsID,
				SkuID:          cacheData.Item.SkuID,
				Num:            float64(cacheData.Item.Quantity),
				Payment:        cacheData.Item.TotalAmount,
				PicPath:        cacheData.Item.ImageURL,
			},
		},
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
	db := s.erpOrderBaseQuery(c).Where("buyer_nick IN ?", s.buildOrderOwnerCandidates(shopUser))
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
		Where("buyer_nick IN ?", s.buildOrderOwnerCandidates(shopUser)).
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
		Where("buyer_nick IN ?", s.buildOrderOwnerCandidates(shopUser)).
		Updates(map[string]interface{}{
			"status":      s.shopStatusToERPStatus(models.OrderStatusCancelled),
			"seller_memo": strings.TrimSpace(reason),
			"update_time": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

// getERPOrderStatistics 统计当前商城用户在 ERP 表中的订单状态数量。
func (s *IApiShopOrderServiceImpl) getERPOrderStatistics(c *gin.Context, shopUser *models.User) (*models.OrderStatistics, error) {
	stats := &models.OrderStatistics{}
	baseQuery := s.erpOrderBaseQuery(c).Where("buyer_nick IN ?", s.buildOrderOwnerCandidates(shopUser))
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(models.OrderStatusPending)).
		Count(&stats.PendingPay).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(models.OrderStatusPaid)).
		Count(&stats.PendingSend).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(models.OrderStatusShipped)).
		Count(&stats.PendingReceive).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(models.OrderStatusCompleted)).
		Count(&stats.Completed).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", s.shopStatusToERPStatus(models.OrderStatusCancelled)).
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
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
}

// buildOrderBuyerNick 生成商城订单在 ERP 表中的买家标识。
func (s *IApiShopOrderServiceImpl) buildOrderBuyerNick(shopUser *models.User) string {
	if shopUser == nil {
		return ""
	}
	return fmt.Sprintf("shop-user-%d", shopUser.ID)
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
	return models.GetStatusText(status)
}

// erpStatusToShopStatus 将 ERP 状态值转换为商城订单状态。
func (s *IApiShopOrderServiceImpl) erpStatusToShopStatus(status string) int32 {
	switch strings.TrimSpace(status) {
	case models.GetStatusText(models.OrderStatusPaid):
		return models.OrderStatusPaid
	case models.GetStatusText(models.OrderStatusShipped):
		return models.OrderStatusShipped
	case models.GetStatusText(models.OrderStatusCompleted):
		return models.OrderStatusCompleted
	case models.GetStatusText(models.OrderStatusCancelled):
		return models.OrderStatusCancelled
	default:
		return models.OrderStatusPending
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
