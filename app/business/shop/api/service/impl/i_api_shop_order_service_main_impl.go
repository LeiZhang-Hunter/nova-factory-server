package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	models2 "nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/api/callback"
	shopordermodels "nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/business/shop/product/shopmodels"
	shopusermodels "nova-factory-server/app/business/shop/user/models"
	orderConstant "nova-factory-server/app/constant/order"
	shopConstant "nova-factory-server/app/constant/shop"
	"nova-factory-server/app/datasource/objectFile"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/observer"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	apimodels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/utils/fileUtils"
	order2 "nova-factory-server/app/utils/order"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"gorm.io/gorm"
)

const (
	orderCachePrefix      = "shop:app:order:cache:"
	orderCreateLockPrefix = "shop:app:order:create:"
	orderCacheTTL         = 10 * time.Minute
	orderCreateLockTTL    = 15 * time.Second
	defaultDeliveryType   = "express"
)

// Confirm 根据 cartId 构建确认单，内部生成 orderKey 并写入预订单缓存。
func (s *IApiShopOrderServiceImpl) Confirm(c *gin.Context, userID int64, req *apimodels.OrderConfirmReq) (*apimodels.OrderConfirmResp, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if req == nil || strings.TrimSpace(req.CartIDValue()) == "" {
		return nil, errors.New("cartId不能为空")
	}

	items, cartIDs, err := s.buildConfirmCacheItems(c, userID, req)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}

	address, err := s.resolveConfirmAddress(c, userID, req.AddressID)
	if err != nil {
		return nil, err
	}
	goodsAmount := sumCacheItems(items)
	cacheData := &apimodels.OrderCacheData{
		OrderKey:       order2.GenerateOrderNo(),
		UserID:         userID,
		Items:          items,
		DeliveryType:   s.resolveConfirmDeliveryType(req),
		GoodsAmount:    goodsAmount,
		FreightAmount:  0,
		DiscountAmount: 0,
		PayAmount:      goodsAmount,
		CartIDs:        cartIDs,
		BuyNow:         req.BuyNow,
	}
	if address != nil {
		cacheData.AddressID = address.ID
	}
	if err := s.fillOrderItemsStock(c, cacheData.Items); err != nil {
		return nil, err
	}

	s.recalculateOrderAmounts(c, userID, cacheData)
	if err := s.saveOrderCache(c, cacheData); err != nil {
		return nil, err
	}

	return &apimodels.OrderConfirmResp{
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
func (s *IApiShopOrderServiceImpl) getItemAvailableStock(c *gin.Context, item *apimodels.OrderCacheItem) (int64, error) {
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
func (s *IApiShopOrderServiceImpl) fillOrderItemsStock(c *gin.Context, items []*apimodels.OrderCacheItem) error {
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
func (s *IApiShopOrderServiceImpl) validateCreateItemsStock(c *gin.Context, items []*apimodels.OrderCacheItem) error {
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
func (s *IApiShopOrderServiceImpl) loadOrderItemSkuMap(c *gin.Context, items []*apimodels.OrderCacheItem) (map[int64]*shopmodels.GoodsSku, error) {
	skuIDs := collectOrderItemSkuIDs(items)
	if len(skuIDs) == 0 {
		return map[int64]*shopmodels.GoodsSku{}, nil
	}

	skuList, err := s.skuDao.ListBySkuIDs(c, skuIDs)
	if err != nil {
		return nil, errors.New("读取商品库存失败")
	}

	if len(skuList) == 0 {
		return nil, errors.New("商品规格不存在")
	}

	skuMap := make(map[int64]*shopmodels.GoodsSku, len(skuList))
	for _, sku := range skuList {
		if sku == nil {
			continue
		}
		skuMap[sku.SkuID] = sku
	}
	return skuMap, nil
}

// calcItemAvailableStock 计算单个商品当前可用库存，活动商品按活动库存与 SKU 库存的较小值返回。
func (s *IApiShopOrderServiceImpl) calcItemAvailableStock(item *apimodels.OrderCacheItem, sku *shopmodels.GoodsSku) int64 {
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
func (s *IApiShopOrderServiceImpl) buildStockInsufficientMessage(item *apimodels.OrderCacheItem) string {
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
func (s *IApiShopOrderServiceImpl) Create(c *gin.Context, userID int64, req *apimodels.OrderCreateReq) (*shopordermodels.Order, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if req == nil || strings.TrimSpace(req.OrderKey) == "" {
		return nil, errors.New("orderKey不能为空")
	}
	req.OrderKey = strings.TrimSpace(req.OrderKey)

	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	if existingOrder, err := s.apiOrderDao.GetByTid(c, req.OrderKey); err != nil {
		return nil, errors.New("读取订单信息失败")
	} else if existingOrder != nil {
		if !s.isOrderOwnedByUser(existingOrder, shopUser) {
			return nil, errors.New("无权操作该订单")
		}
		if err := s.ensureExistingOrderCreated(existingOrder); err != nil {
			return nil, err
		}
		return existingOrder, nil
	}

	lockKey := orderCreateLockPrefix + req.OrderKey
	if !s.cache.SetNX(context.Background(), lockKey, "1", orderCreateLockTTL) {
		return nil, errors.New("订单正在处理中，请勿重复提交")
	}
	defer s.cache.Del(context.Background(), lockKey)

	if existingOrder, err := s.apiOrderDao.GetByTid(c, req.OrderKey); err != nil {
		return nil, errors.New("读取订单信息失败")
	} else if existingOrder != nil {
		if !s.isOrderOwnedByUser(existingOrder, shopUser) {
			return nil, errors.New("无权操作该订单")
		}
		if err := s.ensureExistingOrderCreated(existingOrder); err != nil {
			return nil, err
		}
		return existingOrder, nil
	}
	// 读取预订单缓存
	cacheData, err := s.getOrderCache(c, req.OrderKey)
	if err != nil {
		return nil, err
	}
	if cacheData == nil {
		return nil, errors.New("支付订单已经失效")
	}
	if cacheData.UserID != userID {
		return nil, errors.New("无权操作该预订单")
	}
	if len(cacheData.Items) == 0 {
		return nil, errors.New("预订单商品不能为空")
	}
	// 读取收货地址
	address, err := s.resolveConfirmAddress(c, userID, cacheData.AddressID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("请先选择收货地址")
	}
	// 验证商品库存
	if err := s.validateCreateItemsStock(c, cacheData.Items); err != nil {
		return nil, err
	}
	// 重新计算订单金额
	s.recalculateOrderAmounts(c, userID, cacheData)

	orderNo := req.OrderKey
	if s.db == nil {
		return nil, errors.New("数据库连接不存在")
	}
	err = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		c.Set("db", tx)
		defer c.Set("db", nil)
		for _, item := range cacheData.Items {
			if item == nil {
				continue
			}
			if item.Quantity <= 0 {
				return errors.New("下单库存不能为负数")
			}
			// 扣除商品库存
			if err := s.deductOrderItemStockWithLock(c, item); err != nil {
				return err
			}
			// 扣除活动库存
			if err := s.deductOrderItemActivityStockWithLock(c, item); err != nil {
				return err
			}
		}
		// 读取默认账户信息
		defaultAccount := s.iAccountDao.GetDefaultAccountNo(c)
		orderData := s.buildERPOrderSet(orderNo, shopUser, address, cacheData, req, defaultAccount)
		// 创建订单
		if err := s.syncCreatedOrder(tx, c, orderData, cacheData); err != nil {
			return fmt.Errorf("创建订单失败")
		}

		if len(cacheData.CartIDs) > 0 {
			if err := s.cartDao.DeleteByIds(c, userID, cacheData.CartIDs); err != nil {
				return fmt.Errorf("删除购物车记录失败: %v", err)
			}
		}
		// 删除订单缓存
		s.cache.Del(context.Background(), s.buildOrderCacheKey(req.OrderKey))
		return nil
	})

	if err != nil {
		if isOrderStockError(err) {
			return nil, err
		}
		zap.L().Error("订单创建失败", zap.Error(err))
		return nil, errors.New(err.Error())
	}

	latestOrder, err := s.apiOrderDao.GetByTid(c, orderNo)
	if err != nil || latestOrder == nil {
		return nil, errors.New("读取订单信息失败")
	}
	return latestOrder, nil
}

// GetByID 获取订单详情，包含商品明细。
func (s *IApiShopOrderServiceImpl) GetByID(c *gin.Context, id int64) (*apimodels.ApiOrderVO, error) {
	if id == 0 {
		return nil, errors.New("订单ID不能为空")
	}

	order, err := s.apiOrderDao.GetByID(c, uint64(id))
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	orderIDs := []uint64{order.ID}
	details, err := s.detailDao.ListByOrderIDs(c, orderIDs)
	if err != nil {
		return nil, err
	}
	//accounts, err := s.accountDao.ListByOrderIDs(c, orderIDs)
	//if err != nil {
	//	return nil, err
	//}
	order.Details = details
	//order.Accounts = accounts

	return apimodels.ToApiShopOrderVO(order), nil
}

// List 获取当前用户的订单列表。
func (s *IApiShopOrderServiceImpl) List(c *gin.Context, userID int64, query *apimodels.OrderQuery) (*apimodels.OrderListData, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if query == nil {
		query = &apimodels.OrderQuery{}
	}

	query.UserID = userID

	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	if s.apiOrderDao == nil {
		return nil, errors.New("订单DAO不存在")
	}
	list, err := s.apiOrderDao.ListShopOrders(c, shopUser, query)
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
func (s *IApiShopOrderServiceImpl) UpdateStatus(c *gin.Context, userID int64, req *apimodels.OrderStatusReq) error {
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
	order, err := s.apiOrderDao.GetByID(c, uint64(req.ID))
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
	currentStatus := orderConstant.ErpStatusToShopStatus(order.Status)
	if !s.isValidStatusTransition(currentStatus, req.Status) {
		return errors.New("非法的状态流转")
	}

	rowsAffected, err := s.apiOrderDao.UpdateERPOrderStatus(c, req.ID, shopUser, req.Status)
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

// BatchUpdateStatus 更新订单状态，验证状态流转合法性,从观察者里触发，调用观察者 OnOrderStatusChange,就不要再调用这个函数了，防止多次更新
func (s *IApiShopOrderServiceImpl) BatchUpdateStatus(c *gin.Context, userID int64, statuses []apimodels.OrderStatus) error {
	if userID == 0 {
		return errors.New("用户未登录")
	}
	if len(statuses) == 0 {
		return errors.New("订单ID不能为空")
	}

	var ids []int64 = make([]int64, 0, len(statuses))
	for _, status := range statuses {
		ids = append(ids, int64(status.ID))
	}

	statusesMap := make(map[int64]apimodels.OrderStatus)
	for _, status := range statuses {
		statusesMap[status.ID] = status
	}

	// 验证用户权限
	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return errors.New("商城用户不存在")
	}
	orders, err := s.apiOrderDao.GetByIDs(c, ids)
	if err != nil {
		return errors.New("订单不存在")
	}

	orderList := make(map[int32][]int64)

	for _, order := range orders {
		if !s.isOrderOwnedByUser(order, shopUser) {
			zap.L().Error("无权操作此订单")
			continue
		}
		requestStatus, ok := statusesMap[int64(order.ID)]
		if !ok {
			continue
		}
		// 验证状态流转
		currentStatus := orderConstant.ErpStatusToShopStatus(order.Status)
		if !s.isValidStatusTransition(currentStatus, requestStatus.Status) {
			zap.L().Error("非法的状态流转")
			continue
		}
		_, ok = orderList[currentStatus]
		if !ok {
			orderList[currentStatus] = make([]int64, 0)
		}
	}

	for status, idsArr := range orderList {
		rowsAffected, err := s.apiOrderDao.BatchUpdateERPOrderStatus(c, idsArr, shopUser, status)
		if err != nil {
			return fmt.Errorf("更新订单状态失败: %v", err)
		}
		if rowsAffected == 0 {
			continue
		}
	}

	return nil
}

// Pay 支付订单，调用微信V3 JSAPI预下单并返回调起支付参数。
func (s *IApiShopOrderServiceImpl) Pay(c *gin.Context, userID int64, id int64) (*apimodels.OrderPayResp, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if id == 0 {
		return nil, errors.New("订单ID不能为空")
	}
	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	order, err := s.apiOrderDao.GetByID(c, uint64(id))
	if err != nil || order == nil {
		return nil, errors.New("订单不存在")
	}
	if !s.isOrderOwnedByUser(order, shopUser) {
		return nil, errors.New("无权操作此订单")
	}
	if orderConstant.ErpStatusToShopStatus(order.Status) != orderConstant.OrderStatusPending {
		return nil, errors.New("只能支付待支付的订单")
	}
	//if order.SyncStatus != shopConstant.OrderSyncStatusSuccess {
	//	return nil, errors.New("订单尚未同步管家婆，暂不能支付")
	//}
	if shopUser.WechatOpenid == "" {
		return nil, errors.New("用户未绑定微信")
	}

	// 读取微信配置
	cfgMap, err := s.loadWechatConfig(c)
	if err != nil {
		return nil, err
	}
	appId := cfgMap["wechat_mini_program_app_id"]
	mchId := cfgMap["wechat_pay_mch_id"]
	apiV3Key := cfgMap["wechat_pay_api_v3_key"]
	serialNo := cfgMap["wechat_pay_serial_no"]
	privateKeyPath := cfgMap["wechat_pay_private_key_path"]
	notifyUrl := cfgMap["wechat_pay_notify_url"]
	//platformPublicKeyPath := cfgMap["wechat_pay_platform_public_key_path"]
	if appId == "" || mchId == "" || apiV3Key == "" || serialNo == "" || privateKeyPath == "" || notifyUrl == "" {
		return nil, errors.New("微信支付配置不完整，请在后台管理配置微信支付参数")
	}
	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥文件失败: %v", err)
	}

	// 初始化微信V3客户端
	client, err := wechat.NewClientV3(mchId, serialNo, apiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	// 预下单（金额元转分，四舍五入防浮点截断）
	total := int64(math.Round(order.Total * 100))
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	bm.Set("appid", appId).
		Set("mchid", mchId).
		Set("description", "订单支付: "+order.Tid).
		Set("out_trade_no", order.Tid).
		Set("time_expire", expire).
		Set("notify_url", notifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", total).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", shopUser.WechatOpenid)
		})

	wxRsp, err := client.V3TransactionJsapi(c.Request.Context(), bm)
	if err != nil {
		return nil, fmt.Errorf("微信预下单失败: %v", err)
	}
	if wxRsp.Code != 0 {
		return nil, fmt.Errorf("微信预下单失败: %s", wxRsp.Error)
	}

	// 生成小程序调起支付签名
	paySign, err := client.PaySignOfApplet(appId, wxRsp.Response.PrepayId)
	if err != nil {
		return nil, fmt.Errorf("生成支付签名失败: %v", err)
	}

	return &apimodels.OrderPayResp{
		AppId:     paySign.AppId,
		TimeStamp: paySign.TimeStamp,
		NonceStr:  paySign.NonceStr,
		Package:   "prepay_id=" + wxRsp.Response.PrepayId,
		SignType:  "RSA",
		PaySign:   paySign.PaySign,
	}, nil
}

var wechatConfigKeys = []string{
	"wechat_mini_program_app_id",
	"wechat_pay_mch_id",
	"wechat_pay_api_v3_key",
	"wechat_pay_serial_no",
	"wechat_pay_private_key_path",
	"wechat_pay_notify_url",
	"wechat_pay_platform_public_key_path",
}

// loadWechatConfig 批量读取微信配置并转为 key→value map。
func (s *IApiShopOrderServiceImpl) loadWechatConfig(c *gin.Context) (map[string]string, error) {
	rows, err := s.configDao.GetByConfigKeys(c, wechatConfigKeys)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付配置失败: %v", err)
	}
	cfgMap := make(map[string]string)
	for _, row := range rows {
		cfgMap[row.ConfigKey] = row.ConfigValue
	}
	return cfgMap, nil
}

// HandleWechatNotify 处理微信支付回调。
// func (s *IApiShopOrderServiceImpl) HandleWechatNotify(c *gin.Context, outTradeNo, transactionId, notifyRaw, mchId, appid, payerOpenid string, notifyTotalInt int64) error {
func (s *IApiShopOrderServiceImpl) HandleWechatNotify(e event.ZOrderStatusSyncReqEvent) error {
	c := e.GetCtx()
	var notifyData apimodels.PayNotifyData
	if orderInfo, ok := e.Metadata()["order"]; ok {
		notifyData = orderInfo.(apimodels.PayNotifyData)
	} else {
		return errors.New("订单信息不存在")
	}
	order, err := s.apiOrderDao.GetByTid(c, notifyData.OutTradeNo)
	if err != nil || order == nil {
		return errors.New("订单不存在")
	}

	// 幂等：已支付则跳过
	//if order.Status == orderConstant.ERPStatusPayed {
	//	return nil
	//}
	if order.Status != orderConstant.ERPStatusNoPay {
		return errors.New("订单状态错误")
	}
	rawBytes, _ := json.Marshal(notifyData)
	notifyRaw := string(rawBytes)
	// 金额校验
	orderTotal := int64(math.Round(order.Total * 100))
	if orderTotal != notifyData.Amount.Total {
		return fmt.Errorf("金额校验失败: 订单金额%d分, 回调金额%d分", orderTotal, notifyData.Amount.Total)
	}

	now := time.Now()
	order.TransactionID = notifyData.TransactionID
	order.NotifyRaw = notifyRaw
	order.MchID = notifyData.MchID
	order.AppID = notifyData.AppID
	order.PayerOpenid = notifyData.Payer.Openid
	order.PayTime = &now
	order.Status = orderConstant.ERPStatusPayed
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		return s.apiOrderDao.MarkOrderPaidWithTx(c, tx, order.ID, order.PayTime, notifyData.TransactionID, notifyRaw, notifyData.MchID, notifyData.AppID, notifyData.Payer.Openid)
	})
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
	order, err := s.apiOrderDao.GetByID(c, uint64(id))
	if err != nil {
		return errors.New("订单不存在")
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	if !s.isOrderOwnedByUser(order, shopUser) {
		return errors.New("无权操作此订单")
	}

	if orderConstant.ErpStatusToShopStatus(order.Status) != orderConstant.OrderStatusPending {
		return errors.New("只能取消待支付的订单")
	}

	order.ID = uint64(id)
	order.Status = orderConstant.ShopStatusToERPStatus(orderConstant.OrderStatusCancelled)
	orderEvent := shopordermodels.NewOrderStatusSyncEvent([]*shopordermodels.Order{
		order,
	}, orderConstant.Normal)
	orderEvent.WithCache(s.cache)
	orderEvent.WithDB(s.db)
	orderEvent.WithCtx(c)
	orderEvent.WithTransaction(true)
	orderEvent.WithCallback(callback.NewShopOrderStatusApiCallback(orderEvent))
	err = observer.GetNotifier().OnOrderStatusChange(orderEvent)
	if err != nil {
		zap.L().Error("on order status change err", zap.Error(err))
		return err
	}

	rowsAffected, err := s.apiOrderDao.CancelERPOrder(c, id, shopUser, reason)
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
	order, err := s.apiOrderDao.GetByID(c, uint64(id))
	if err != nil {
		return errors.New("订单不存在")
	}
	if order == nil {
		return errors.New("订单不存在")
	}

	if !s.isOrderOwnedByUser(order, shopUser) {
		return errors.New("无权操作此订单")
	}

	if orderConstant.ErpStatusToShopStatus(order.Status) != orderConstant.OrderStatusShipped {
		return errors.New("只能确认已发货的订单")
	}

	rowsAffected, err := s.apiOrderDao.UpdateERPOrderStatus(c, id, shopUser, orderConstant.OrderStatusCompleted)
	if err != nil {
		return fmt.Errorf("确认收货失败: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("订单状态已更新，请刷新后重试")
	}

	return nil
}

func isOrderStockError(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "库存不足") || strings.Contains(message, "购买数量必须大于0")
}

// deductOrderItemStockWithLock 在事务内锁定 SKU 行，确认库存充足后再扣减。
func (s *IApiShopOrderServiceImpl) deductOrderItemStockWithLock(c *gin.Context, item *apimodels.OrderCacheItem) error {
	if item == nil {
		return nil
	}
	if item.Quantity <= 0 {
		return errors.New("购买数量必须大于0")
	}

	sku, err := s.skuDao.GetBySkuIDForUpdate(c, item.SkuID)
	if err != nil {
		return errors.New("读取商品库存失败")
	}
	if sku == nil {
		return errors.New("sku不存在")
	}
	availableStock := sku.Quantity
	if item.Quantity > availableStock {
		return errors.New("下单失败，库存不足: " + s.buildStockInsufficientDetail(item.GoodsName, item.SkuName, item.Quantity, availableStock))
	}
	if err := s.skuDao.DeductStockBySkuId(c, item.SkuID, item.Quantity); err != nil {
		return fmt.Errorf("扣减库存失败: %v", err)
	}
	// 扣减 SKU 库存后重新计算商品总库存
	if err := s.recalcGoodsStockByGoodsID(c, sku.GoodsID); err != nil {
		return err
	}
	return nil
}

// deductOrderItemActivityStockWithLock 在事务内锁定活动库存行，确认库存充足后再扣减。
func (s *IApiShopOrderServiceImpl) deductOrderItemActivityStockWithLock(c *gin.Context, item *apimodels.OrderCacheItem) error {
	if item == nil {
		return nil
	}
	if item.SecKillId > 0 {
		seckill, err := s.seckillDao.GetByIDForUpdate(c, item.SecKillId)
		if err != nil {
			return errors.New("读取秒杀库存失败")
		}
		availableStock := int64(0)
		if seckill != nil {
			availableStock = seckill.Stock
		}
		if item.Quantity > availableStock {
			return errors.New("下单失败，库存不足: " + s.buildStockInsufficientDetail(item.GoodsName, item.SkuName, item.Quantity, availableStock))
		}
		if err := s.seckillDao.DeductStock(c, item.SecKillId, item.Quantity); err != nil {
			return fmt.Errorf("扣减秒杀库存失败: %v", err)
		}
		// 扣减秒杀库存后重新计算商品总库存
		var goods apimodels.Goods
		if err := s.getTxDB(c).WithContext(c).Table("shop_goods").
			Where("id = ?", seckill.ProductID).
			First(&goods).Error; err != nil {
			return fmt.Errorf("查询商品信息失败: %v", err)
		}
		if err := s.recalcGoodsStockByGoodsID(c, goods.GoodsID); err != nil {
			return err
		}
	}
	if item.CombinationId > 0 {
		combination, err := s.combDao.GetByIDForUpdate(c, item.CombinationId)
		if err != nil {
			return errors.New("读取拼团库存失败")
		}
		availableStock := int64(0)
		if combination != nil {
			availableStock = combination.Stock
		}
		if item.Quantity > availableStock {
			return errors.New("下单失败，库存不足: " + s.buildStockInsufficientDetail(item.GoodsName, item.SkuName, item.Quantity, availableStock))
		}
		if err := s.combDao.DeductStock(c, item.CombinationId, item.Quantity); err != nil {
			return fmt.Errorf("扣减拼团库存失败: %v", err)
		}
	}
	return nil
}

// GetStatistics 获取当前用户各状态订单数量统计。
func (s *IApiShopOrderServiceImpl) GetStatistics(c *gin.Context, userID int64) (*apimodels.OrderStatistics, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	shopUser, err := s.userDao.GetByUserID(c, userID)
	if err != nil || shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	return s.apiOrderDao.GetERPOrderStatistics(c, shopUser)
}

func (s *IApiShopOrderServiceImpl) syncCreatedOrder(tx *gorm.DB, c *gin.Context, order *shopordermodels.OrderSet, cacheData *apimodels.OrderCacheData) error {
	if order == nil {
		return errors.New("订单不存在")
	}
	orderEvent := shopordermodels.BuildShopOrderSyncEvent(order)
	orderEvent.WithCache(s.cache)
	orderEvent.WithDB(tx)
	//orderEvent.WithCallback(callback.NewOrderSyncRequestCallback(c, orderEvent, order))
	if err := observer.GetNotifier().OnOrderChanged(orderEvent); err != nil {
		return fmt.Errorf("订单同步观察者失败: %v", err)
	}
	return nil
}

func (s *IApiShopOrderServiceImpl) ensureExistingOrderCreated(order *shopordermodels.Order) error {
	if order == nil {
		return errors.New("订单不存在")
	}
	if strings.TrimSpace(order.Status) == orderConstant.ERPStatusTradeClosed {
		return errors.New("订单同步管家婆失败，请重新下单")
	}
	if order.SyncStatus == shopConstant.OrderSyncStatusFailed {
		return errors.New("订单同步管家婆失败，请重新下单")
	}
	if order.SyncStatus != shopConstant.OrderSyncStatusSuccess {
		return errors.New("订单尚未同步管家婆，请稍后重试")
	}
	return nil
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

func (s *IApiShopOrderServiceImpl) saveOrderCache(c *gin.Context, data *apimodels.OrderCacheData) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	s.cache.Set(context.Background(), s.buildOrderCacheKey(data.OrderKey), string(body), orderCacheTTL)
	return nil
}

func (s *IApiShopOrderServiceImpl) getOrderCache(c *gin.Context, orderKey string) (*apimodels.OrderCacheData, error) {
	val, err := s.cache.Get(context.Background(), s.buildOrderCacheKey(orderKey))
	if err != nil {
		return nil, errors.New("预订单已失效，请重新确认商品")
	}
	var data apimodels.OrderCacheData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, errors.New("预订单数据异常")
	}
	return &data, nil
}

func (s *IApiShopOrderServiceImpl) getBuyNowCartCache(c *gin.Context, cartID string) (*apimodels.OrderCacheData, error) {
	val, err := s.cache.Get(context.Background(), buildBuyNowCartCacheKey(cartID))
	if err != nil {
		return nil, errors.New("立即购买商品已失效，请重新选择商品")
	}
	var data apimodels.OrderCacheData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, errors.New("立即购买商品数据异常")
	}
	return &data, nil
}

func (s *IApiShopOrderServiceImpl) buildConfirmCacheItems(c *gin.Context, userID int64, req *apimodels.OrderConfirmReq) ([]*apimodels.OrderCacheItem, []int64, error) {
	if req == nil {
		return nil, nil, errors.New("参数不能为空")
	}
	if req.BuyNow {
		cacheData, err := s.getBuyNowCartCache(c, strings.TrimSpace(req.CartIDValue()))
		if err != nil {
			return nil, nil, err
		}
		if cacheData.UserID != userID {
			return nil, nil, errors.New("无权操作该立即购买商品")
		}
		if len(cacheData.Items) == 0 {
			return nil, nil, errors.New("立即购买商品不存在")
		}
		return cacheData.Items, nil, nil
	}

	idList, err := s.parseCartIDString(req.CartIDValue())
	if err != nil {
		return nil, nil, err
	}
	items, err := s.buildCartCacheItemsByState(c, userID, idList, false)
	if err != nil {
		return nil, nil, err
	}
	return items, idList, nil
}

func (s *IApiShopOrderServiceImpl) buildCartCacheItemsByState(c *gin.Context, userID int64, idList []int64, buyNow bool) ([]*apimodels.OrderCacheItem, error) {
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}
	if len(idList) == 0 {
		return nil, errors.New("购物车ID不能为空")
	}

	state := shopConstant.CartStateNormal
	if buyNow {
		state = shopConstant.CartStateBuyNow
	}
	cartList, err := s.cartDao.ListByIDsAndState(c, userID, idList, state)
	if err != nil {
		return nil, errors.New("读取购物车商品失败")
	}
	if len(cartList) == 0 {
		return nil, errors.New("购物车商品不存在")
	}

	cartMap := make(map[int64]*apimodels.CartDto, len(cartList))
	for _, cart := range cartList {
		if cart != nil {
			cartMap[cart.ID] = cart
		}
	}

	items := make([]*apimodels.OrderCacheItem, 0, len(idList))
	for _, id := range idList {
		cartInfo, ok := cartMap[id]
		if !ok {
			return nil, errors.New("部分购物车商品不存在")
		}
		item, err := s.buildCartInfoCacheItem(c, userID, cartInfo)
		if err != nil {
			return nil, err
		}
		item.CartID = cartInfo.ID
		items = append(items, item)
	}
	return items, nil
}

func (s *IApiShopOrderServiceImpl) buildCartInfoCacheItem(c *gin.Context, userID int64, cartInfo *apimodels.CartDto) (*apimodels.OrderCacheItem, error) {
	if cartInfo == nil {
		return nil, errors.New("购物车商品不存在")
	}
	switch cartInfo.ProductType {
	case shopConstant.CartProductTypeSeckill:
		return s.buildSeckillCacheItem(c, cartInfo.ActivityID, cartInfo.SkuID, cartInfo.Quantity)
	case shopConstant.CartProductTypeCombination:
		item, err := s.buildCombinationCacheItem(c, cartInfo.ActivityID, cartInfo.SkuID, cartInfo.Quantity)
		if err != nil {
			return nil, err
		}
		item.PinkId = cartInfo.PinkID
		return item, nil
	default:
		return s.buildSingleCacheItem(c, userID, cartInfo.GoodsID, cartInfo.SkuID, cartInfo.Quantity)
	}
}

func (s *IApiShopOrderServiceImpl) parseCartIDString(cartID string) ([]int64, error) {
	return s.parseCartIDs(strings.Split(cartID, ","))
}

func (s *IApiShopOrderServiceImpl) resolveConfirmDeliveryType(req *apimodels.OrderConfirmReq) string {
	if req == nil {
		return defaultDeliveryType
	}
	if deliveryType := strings.TrimSpace(req.DeliveryType); deliveryType != "" {
		return deliveryType
	}
	return defaultDeliveryType
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
func (s *IApiShopOrderServiceImpl) buildSeckillCacheItem(c *gin.Context, secKillID int64, skuID int64, quantity int64) (*apimodels.OrderCacheItem, error) {
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
func (s *IApiShopOrderServiceImpl) buildCombinationCacheItem(c *gin.Context, combinationID int64, skuID int64, quantity int64) (*apimodels.OrderCacheItem, error) {
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

	goods, sku, err := s.loadGoodsAndSkuByGoodsCode(c, combination.ProductID, skuID, quantity)
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
func (s *IApiShopOrderServiceImpl) buildDirectCacheItem(c *gin.Context, userID int64, skuID int64, quantity int64) (*apimodels.OrderCacheItem, error) {
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
func (s *IApiShopOrderServiceImpl) loadGoodsAndSkuByGoodsID(c *gin.Context, goodsID int64, skuID int64, quantity int64) (*apimodels.Goods, *shopmodels.GoodsSku, error) {
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
func (s *IApiShopOrderServiceImpl) loadGoodsAndSkuByGoodsCode(c *gin.Context, goodsId int64, skuID int64, quantity int64) (*apimodels.Goods, *shopmodels.GoodsSku, error) {
	if goodsId == 0 {
		return nil, nil, errors.New("活动商品数据异常")
	}
	goods, err := s.goodsDao.GetByGoodsID(c, goodsId)
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
	goods *apimodels.Goods,
	sku *shopmodels.GoodsSku,
	quantity int64,
	price float64,
	goodsName string,
	imageURL string,
) *apimodels.OrderCacheItem {
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
	return &apimodels.OrderCacheItem{
		GoodsID:     goods.GoodsID,
		SkuID:       int64(sku.SkuID),
		GoodsName:   finalGoodsName,
		SkuName:     sku.SkuName,
		ImageURL:    fileUtils.BuildAbsoluteURL(c, finalImageURL),
		Price:       price,
		Quantity:    quantity,
		TotalAmount: price * float64(quantity),
		OuterIid:    sku.OuterID,
	}
}

func (s *IApiShopOrderServiceImpl) buildSingleCacheItem(c *gin.Context, userID int64, goodsID int64, skuID int64, quantity int64) (*apimodels.OrderCacheItem, error) {
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
	item := &apimodels.OrderCacheItem{
		GoodsID:     goodsID,
		SkuID:       skuID,
		GoodsName:   goods.GoodsName,
		SkuName:     sku.SkuName,
		ImageURL:    fileUtils.BuildAbsoluteURL(c, sku.ImageURL),
		Price:       price,
		Quantity:    quantity,
		TotalAmount: price * float64(quantity),
		OuterIid:    sku.OuterID,
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

func (s *IApiShopOrderServiceImpl) resolveConfirmAddress(c *gin.Context, userID int64, addressID int64) (*apimodels.ShopUserAddressApp, error) {
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

func (s *IApiShopOrderServiceImpl) recalculateOrderAmounts(c *gin.Context, userID int64, data *apimodels.OrderCacheData) {
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
func (s *IApiShopOrderServiceImpl) applyGoodsDiscount(c *gin.Context, userID int64, goods *apimodels.Goods, skuID int64, price float64) float64 {
	if s.discountService == nil || userID == 0 || goods == nil || price <= 0 {
		return price
	}
	discountedPrice, hasDiscount := s.discountService.CalculateDiscountPrice(
		c, userID, goods.GoodsID, skuID, goods.ShopCategoryId, price)
	if hasDiscount {
		return discountedPrice
	}
	return price
}

func sumCacheItems(items []*apimodels.OrderCacheItem) float64 {
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
func collectOrderItemSkuIDs(items []*apimodels.OrderCacheItem) []int64 {
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

// buildOrderBuyerNick 生成商城订单在 ERP 表中的买家标识。
func (s *IApiShopOrderServiceImpl) buildOrderBuyerNick(shopUser *shopusermodels.User) string {
	if shopUser == nil {
		return ""
	}
	return fmt.Sprintf("shop-user-%s", shopUser.UserID)
}

// buildOrderOwnerCandidates 构造商城用户在 ERP 订单中的可能归属标识。
func (s *IApiShopOrderServiceImpl) buildOrderOwnerCandidates(shopUser *shopusermodels.User) []string {
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
func (s *IApiShopOrderServiceImpl) isOrderOwnedByUser(order *shopordermodels.Order, shopUser *shopusermodels.User) bool {
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

// getTxDB 获取上下文中事务 db，若无则返回默认 db。
func (s *IApiShopOrderServiceImpl) getTxDB(c *gin.Context) *gorm.DB {
	if tx, ok := c.Get("db"); ok {
		if txDB, ok := tx.(*gorm.DB); ok {
			return txDB
		}
	}
	return s.db
}

// recalcGoodsStockByGoodsID 根据 goodsID 汇总所有 SKU 库存并更新商品总库存。
func (s *IApiShopOrderServiceImpl) recalcGoodsStockByGoodsID(c *gin.Context, goodsID int64) error {
	skus, err := s.skuDao.ListByGoodsIDs(c, []int64{goodsID})
	if err != nil {
		return fmt.Errorf("查询商品SKU列表失败: %v", err)
	}
	var totalStock int64
	for _, sk := range skus {
		totalStock += sk.Quantity
	}
	if err := s.getTxDB(c).WithContext(c).Table("shop_goods").
		Where("goods_id = ?", goodsID).
		Update("quantity", totalStock).Error; err != nil {
		return fmt.Errorf("更新商品总库存失败: %v", err)
	}
	return nil
}
