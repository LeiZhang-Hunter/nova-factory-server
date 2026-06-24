package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	activityModels "nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/api/dao"
	api "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/business/shop/product/shopmodels"
	shopConstant "nova-factory-server/app/constant/shop"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/fileUtils"
	orderUtils "nova-factory-server/app/utils/order"
	"nova-factory-server/app/utils/stringUtils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// cartDefaultDelivery 立即购买快照的默认配送方式，保持与订单确认缓存结构一致。
	cartDefaultDelivery   = "express"
	cartBuyNowCachePrefix = "shop:app:cart:buy-now:"
	cartBuyNowCacheTTL    = time.Hour
)

// cartResolvedItem 是 GenCart 保存购物车或立即购买快照前的统一商品上下文。
// 普通商品、秒杀商品、拼团商品都会先解析成这个结构，后续流程不再关心请求来源。
type cartResolvedItem struct {
	GoodsID         int64                               // 基础商品主键 ID。
	SkuID           int64                               // 商品规格主键 ID。
	GoodsName       string                              // 下单或购物车展示用商品名称，活动商品优先使用活动标题。
	SkuName         string                              // 下单或购物车展示用规格名称。
	ImageURL        string                              // 下单或购物车展示用图片，活动图优先，其次 SKU 图和商品图。
	Price           float64                             // 当前购买价格，普通商品取折扣后 SKU 价，活动商品取活动价。
	Quantity        int64                               // 本次购买数量。
	AvailableStock  int64                               // 当前可售库存，取活动库存和 SKU 库存的较小值。
	ProductType     int32                               // 购物车商品类型：普通、秒杀、拼团。
	ActivityID      int64                               // 秒杀或拼团活动商品 ID，普通商品为 0。
	PinkID          int64                               // 拼团团队 ID，仅拼团加入已有团队时有值。
	SeckillInfo     *activityModels.SeckillMainInfo     // 立即购买确认页需要的秒杀活动快照。
	CombinationInfo *activityModels.CombinationMainInfo // 立即购买确认页需要的拼团活动快照。
}

type IApiShopCartServiceImpl struct {
	cache           cache.Cache
	goodsDao        dao.IApiShopGoodsDao
	goodsSkuDao     dao.IApiShopSkuDao
	shopUserDao     dao.IApiShopWechatUserDao
	seckillDao      dao.IApiShopSeckillDao
	combinationDao  dao.IApiShopCombinationDao
	pinkDao         dao.IApiShopPinkDao
	iShopCartDao    dao.IApiShopCartDao
	discountService discountservice.IDiscountCalculateService
}

// NewIApiShopCartServiceImpl 购物车服务
func NewIApiShopCartServiceImpl(
	cache cache.Cache,
	goodsDao dao.IApiShopGoodsDao,
	goodsSkuDao dao.IApiShopSkuDao,
	shopUserDao dao.IApiShopWechatUserDao,
	seckillDao dao.IApiShopSeckillDao,
	combinationDao dao.IApiShopCombinationDao,
	pinkDao dao.IApiShopPinkDao,
	iShopCartDao dao.IApiShopCartDao,
	discountService discountservice.IDiscountCalculateService,
) service.IApiShopCartService {
	return &IApiShopCartServiceImpl{
		cache:           cache,
		goodsDao:        goodsDao,
		goodsSkuDao:     goodsSkuDao,
		iShopCartDao:    iShopCartDao,
		shopUserDao:     shopUserDao,
		seckillDao:      seckillDao,
		combinationDao:  combinationDao,
		pinkDao:         pinkDao,
		discountService: discountService,
	}
}

// GenCart 统一处理 App 端加入购物车和立即购买，普通/秒杀/拼团都会返回 cartId。
func (s *IApiShopCartServiceImpl) GenCart(c *gin.Context, req *api.CartSetDataReq) (*api.CartSetDataResp, error) {
	// 基础参数先在入口统一拦截，避免后续 DAO 查询出现无意义的空值。
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	if req.SkuID == 0 {
		return nil, errors.New("SKU ID不能为空")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("购买数量必须大于0")
	}

	// App 用户信息用于写入购物车归属和用户名快照。
	userId := baizeContext.GetUserId(c)
	userInfo, err := s.shopUserDao.GetByUserID(c, userId)
	if err != nil {
		zap.L().Error("get user info failed", zap.Error(err))
		return nil, errors.New("读取用户信息失败")
	}
	if userInfo == nil {
		return nil, errors.New("用户不存在")
	}

	// 下面把普通商品、秒杀商品、拼团商品归一成同一个商品上下文，后续保存逻辑只依赖 item。
	var item *cartResolvedItem
	if req.SecKillID > 0 && req.CombinationID > 0 {
		return nil, errors.New("不能同时选择秒杀和拼团商品")
	}

	if req.SecKillID > 0 {
		// 秒杀商品以 secKillId 为活动身份，价格使用秒杀价，库存取秒杀库存和 SKU 库存的较小值。
		seckill, err := s.seckillDao.GetByID(c, req.SecKillID)
		if err != nil {
			return nil, errors.New("读取秒杀商品失败")
		}
		if seckill == nil {
			return nil, errors.New("秒杀商品不存在")
		}
		if seckill.IsShow != 1 || seckill.Status != 1 {
			return nil, errors.New("秒杀商品已下架")
		}
		if seckill.Num > 0 && req.Quantity > int64(seckill.Num) {
			return nil, errors.New("超过秒杀单次限购数量")
		}
		if req.GoodsID > 0 && req.GoodsID != seckill.ProductID {
			return nil, errors.New("商品与秒杀活动不匹配")
		}

		goods, sku, err := s.loadGoodsAndSkuByID(c, seckill.ProductID, req.SkuID)
		if err != nil {
			return nil, err
		}
		availableStock := minCartStock(seckill.Stock, sku.Quantity)
		if err := validateCartStock(stringUtils.FirstNonEmpty(seckill.Title, goods.GoodsName), sku.SkuName, req.Quantity, availableStock); err != nil {
			return nil, err
		}
		item = &cartResolvedItem{
			GoodsID:        goods.ID,
			SkuID:          int64(sku.ID),
			GoodsName:      stringUtils.FirstNonEmpty(seckill.Title, goods.GoodsName),
			SkuName:        sku.SkuName,
			ImageURL:       stringUtils.FirstNonEmpty(seckill.Image, sku.ImageURL, goods.ImageURL),
			Price:          seckill.Price,
			Quantity:       req.Quantity,
			AvailableStock: availableStock,
			ProductType:    shopConstant.CartProductTypeSeckill,
			ActivityID:     req.SecKillID,
			SeckillInfo: &activityModels.SeckillMainInfo{
				ID:         seckill.ID,
				ActivityID: seckill.ActivityID,
				ProductID:  seckill.ProductID,
				Price:      seckill.Price,
				Cost:       seckill.Cost,
				OtPrice:    seckill.OtPrice,
				Stock:      seckill.Stock,
				StartTime:  seckill.StartTime,
				StopTime:   seckill.StopTime,
				Status:     seckill.Status,
				IsPostage:  seckill.IsPostage,
				IsHot:      seckill.IsHot,
				Num:        seckill.Num,
				IsShow:     seckill.IsShow,
				TimeID:     seckill.TimeID,
				Quota:      seckill.Quota,
				QuotaShow:  seckill.QuotaShow,
				OnceNum:    seckill.Num,
			},
		}
	} else if req.CombinationID > 0 {
		// 拼团商品以 combinationId 为活动身份，pinkId 存在时表示加入已有团队。
		combination, err := s.combinationDao.GetByID(c, req.CombinationID)
		if err != nil {
			return nil, errors.New("读取拼团商品失败")
		}
		if combination == nil {
			return nil, errors.New("拼团商品不存在")
		}
		if combination.IsShow != 1 {
			return nil, errors.New("拼团商品已下架")
		}
		if combination.OnceNum > 0 && req.Quantity > combination.OnceNum {
			return nil, errors.New("超过拼团单次限购数量")
		}

		if req.PinkID > 0 {
			pink, err := s.pinkDao.GetByID(c, req.PinkID)
			if err != nil {
				return nil, errors.New("读取拼团记录失败")
			}
			if pink == nil {
				return nil, errors.New("拼团记录不存在")
			}
			value := strings.TrimSpace(pink.StopTime)
			if value == "" {
				return nil, errors.New("拼团记录数据异常")
			}
			var stopTime time.Time
			if unix, err := strconv.ParseInt(value, 10, 64); err == nil {
				stopTime = time.Unix(unix, 0)
			} else {
				// 尝试多种时间格式解析拼团结束时间，支持标准格式、RFC3339和日期格式。
				for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"} {
					parsed, parseErr := time.ParseInLocation(layout, value, time.Local)
					if parseErr == nil {
						stopTime = parsed
						break
					}
				}
			}
			if stopTime.IsZero() {
				return nil, errors.New("拼团记录数据异常")
			}
			if time.Now().After(stopTime) {
				return nil, errors.New("拼团已到期")
			}
		}

		// 拼团表中的 ProductID 是基础商品的 goods_code，需要通过 goods_code 找回商品。
		goodsCode := combination.ProductID
		if goodsCode == 0 {
			return nil, errors.New("活动商品数据异常")
		}
		goods, err := s.goodsDao.GetByGoodsID(c, goodsCode)
		if err != nil {
			return nil, errors.New("读取商品信息失败")
		}
		if goods == nil {
			return nil, errors.New("商品不存在")
		}
		_, sku, err := s.loadSkuForGoods(c, goods, req.SkuID)
		if err != nil {
			return nil, err
		}
		if req.GoodsID > 0 && req.GoodsID != goods.ID {
			return nil, errors.New("商品与拼团活动不匹配")
		}
		availableStock := minCartStock(combination.Stock, sku.Quantity)
		if err := validateCartStock(stringUtils.FirstNonEmpty(combination.Title, goods.GoodsName), sku.SkuName, req.Quantity, availableStock); err != nil {
			return nil, err
		}
		item = &cartResolvedItem{
			GoodsID:        goods.ID,
			SkuID:          int64(sku.ID),
			GoodsName:      stringUtils.FirstNonEmpty(combination.Title, goods.GoodsName),
			SkuName:        sku.SkuName,
			ImageURL:       stringUtils.FirstNonEmpty(combination.Image, sku.ImageURL, goods.ImageURL),
			Price:          combination.Price,
			Quantity:       req.Quantity,
			AvailableStock: availableStock,
			ProductType:    shopConstant.CartProductTypeCombination,
			ActivityID:     req.CombinationID,
			PinkID:         req.PinkID,
			CombinationInfo: &activityModels.CombinationMainInfo{
				ID:            combination.ID,
				ProductID:     combination.ProductID,
				MerID:         combination.MerID,
				Attr:          combination.Attr,
				Price:         combination.Price,
				Sort:          combination.Sort,
				Sales:         combination.Sales,
				Stock:         combination.Stock,
				IsHost:        combination.IsHost,
				IsShow:        combination.IsShow,
				IsPostage:     combination.IsPostage,
				Postage:       combination.Postage,
				StartTime:     combination.StartTime,
				StopTime:      combination.StopTime,
				EffectiveTime: combination.EffectiveTime,
				Browse:        combination.Browse,
				UnitName:      combination.UnitName,
				Weight:        combination.Weight,
				Volume:        combination.Volume,
				Num:           combination.Num,
				OnceNum:       combination.OnceNum,
				Quota:         combination.Quota,
				QuotaShow:     combination.QuotaShow,
				Virtual:       combination.Virtual,
				HomeModuleIDs: combination.HomeModuleIDs,
			},
		}
	} else {
		// 普通商品以 goodsId + skuId 为身份，价格和库存来自基础商品/SKU。
		if req.GoodsID <= 0 {
			return nil, errors.New("商品ID不能为空")
		}
		goods, sku, err := s.loadGoodsAndSkuByID(c, req.GoodsID, req.SkuID)
		if err != nil {
			return nil, err
		}
		availableStock := minCartStock(goods.Quantity, sku.Quantity)
		if err := validateCartStock(goods.GoodsName, sku.SkuName, req.Quantity, availableStock); err != nil {
			return nil, err
		}
		price := sku.RetailPrice
		if s.discountService != nil && userId != 0 && price > 0 {
			discountedPrice, hasDiscount := s.discountService.CalculateDiscountPrice(
				c, userId, goods.GoodsID, sku.SkuID, goods.ShopCategoryId, price)
			if hasDiscount {
				price = discountedPrice
			}
		}
		item = &cartResolvedItem{
			GoodsID:        goods.ID,
			SkuID:          int64(sku.ID),
			GoodsName:      goods.GoodsName,
			SkuName:        sku.SkuName,
			ImageURL:       stringUtils.FirstNonEmpty(sku.ImageURL, goods.ImageURL),
			Price:          price,
			Quantity:       req.Quantity,
			AvailableStock: availableStock,
			ProductType:    shopConstant.CartProductTypeNormal,
		}
	}

	if req.IsBuyNow() {
		cartID, err := s.saveBuyNowCartSnapshot(c, userId, item)
		if err != nil {
			zap.L().Error("save buy now cart snapshot failed", zap.Error(err))
			return nil, errors.New("保存立即购买快照失败")
		}
		return &api.CartSetDataResp{
			Mode:   shopConstant.CartModeBuyNow,
			CartID: cartID,
		}, nil
	}

	cartInfo, err := s.iShopCartDao.Save(c, &api.CartSetData{
		UserID:      userId,
		Username:    userInfo.Username,
		GoodsID:     item.GoodsID,
		SkuID:       uint64(item.SkuID),
		GoodsName:   item.GoodsName,
		SkuName:     item.SkuName,
		ImageURL:    item.ImageURL,
		RetailPrice: item.Price,
		Quantity:    item.Quantity,
		ProductType: item.ProductType,
		ActivityID:  item.ActivityID,
		PinkID:      item.PinkID,
		State:       shopConstant.CartStateNormal,
	})
	if err != nil {
		zap.L().Error("save cart info failed", zap.Error(err))
		return nil, errors.New("保存购物车失败")
	}

	return &api.CartSetDataResp{
		Mode:   shopConstant.CartModeCart,
		CartID: strconv.FormatInt(cartInfo.ID, 10),
		Cart:   cartInfo,
	}, nil
}

func buildBuyNowCartCacheKey(cartID string) string {
	return cartBuyNowCachePrefix + cartID
}

func (s *IApiShopCartServiceImpl) saveBuyNowCartSnapshot(c *gin.Context, userID int64, item *cartResolvedItem) (string, error) {
	if item == nil {
		return "", errors.New("订单商品不能为空")
	}
	cartID := orderUtils.GenerateOrderNo()
	orderItem := &api.OrderCacheItem{
		CombinationId:     boolToInt64(item.ProductType == shopConstant.CartProductTypeCombination, item.ActivityID),
		SecKillId:         boolToInt64(item.ProductType == shopConstant.CartProductTypeSeckill, item.ActivityID),
		SeckillInfo:       item.SeckillInfo,
		CombinationInfo:   item.CombinationInfo,
		PinkId:            item.PinkID,
		GoodsID:           item.GoodsID,
		SkuID:             item.SkuID,
		GoodsName:         item.GoodsName,
		SkuName:           item.SkuName,
		ImageURL:          fileUtils.BuildAbsoluteURL(c, item.ImageURL),
		Price:             item.Price,
		Quantity:          item.Quantity,
		AvailableStock:    item.AvailableStock,
		StockInsufficient: item.AvailableStock < item.Quantity,
		TotalAmount:       item.Price * float64(item.Quantity),
	}
	cacheData := &api.OrderCacheData{
		OrderKey:       cartID,
		UserID:         userID,
		Items:          []*api.OrderCacheItem{orderItem},
		DeliveryType:   cartDefaultDelivery,
		GoodsAmount:    orderItem.TotalAmount,
		FreightAmount:  0,
		DiscountAmount: 0,
		PayAmount:      orderItem.TotalAmount,
		BuyNow:         true,
	}
	body, err := json.Marshal(cacheData)
	if err != nil {
		return "", err
	}
	s.cache.Set(context.Background(), buildBuyNowCartCacheKey(cartID), string(body), cartBuyNowCacheTTL)
	return cartID, nil
}

// loadGoodsAndSkuByID 按基础商品主键加载商品和 SKU，并复用 SKU 归属校验。
func (s *IApiShopCartServiceImpl) loadGoodsAndSkuByID(c *gin.Context, goodsID int64, skuID int64) (*api.Goods, *shopmodels.GoodsSku, error) {
	goods, err := s.goodsDao.GetByID(c, goodsID)
	if err != nil {
		return nil, nil, errors.New("读取商品信息失败")
	}
	if goods == nil {
		return nil, nil, errors.New("商品不存在")
	}
	return s.loadSkuForGoods(c, goods, skuID)
}

// loadSkuForGoods 校验商品可售状态、SKU 存在性和 SKU 是否归属于该商品。
func (s *IApiShopCartServiceImpl) loadSkuForGoods(c *gin.Context, goods *api.Goods, skuID int64) (*api.Goods, *shopmodels.GoodsSku, error) {
	if goods.IsOnSale != 1 {
		return nil, nil, errors.New("商品已下架")
	}
	if (goods.GoodsID) == 0 {
		return nil, nil, errors.New("商品数据异常")
	}
	sku, err := s.goodsSkuDao.GetByID(c, skuID)
	if err != nil {
		return nil, nil, errors.New("读取商品规格失败")
	}
	if sku == nil {
		return nil, nil, errors.New("商品规格不存在")
	}
	if sku.GoodsID != goods.GoodsID {
		return nil, nil, errors.New("SKU不属于该商品")
	}
	if strings.TrimSpace(sku.SkuName) == "" {
		return nil, nil, errors.New("商品规格名称不能为空")
	}
	return goods, sku, nil
}

// validateCartStock 校验购买数量是否超过当前可售库存，并返回带商品/规格名的业务错误。
func validateCartStock(goodsName, skuName string, quantity int64, availableStock int64) error {
	if availableStock <= 0 {
		return errors.New("暂无库存")
	}
	if availableStock < quantity {
		name := strings.TrimSpace(goodsName)
		if sku := strings.TrimSpace(skuName); sku != "" {
			name = fmt.Sprintf("%s %s", name, sku)
		}
		if strings.TrimSpace(name) == "" {
			name = "商品"
		}
		return fmt.Errorf("%s库存不足", name)
	}
	return nil
}

// minCartStock 返回多个库存来源中的最小值。
// 活动商品需要同时受活动库存和基础 SKU 库存约束。
func minCartStock(values ...int64) int64 {
	var min int64
	for i, value := range values {
		if i == 0 || value < min {
			min = value
		}
	}
	return min
}

func boolToInt64(ok bool, value int64) int64 {
	if ok {
		return value
	}
	return 0
}

// List 查询用户购物车列表。
func (s *IApiShopCartServiceImpl) List(c *gin.Context) ([]*api.CartDto, error) {
	userID := baizeContext.GetUserId(c)
	//var userID int64 = 486436412988592128
	//if userID == 0 {
	//	return nil, errors.New("用户ID不能为空")
	//}
	ret, err := s.iShopCartDao.List(c, userID)
	if err != nil {
		return nil, err
	}

	skuIDs := make([]int64, 0, len(ret))
	skuIDSet := make(map[int64]struct{}, len(ret))
	for _, item := range ret {
		if _, ok := skuIDSet[item.SkuID]; ok {
			continue
		}
		skuIDSet[item.SkuID] = struct{}{}
		skuIDs = append(skuIDs, item.SkuID)
	}

	skuQuantityMap := make(map[int64]int64, len(skuIDs))
	if len(skuIDs) > 0 {
		skuList, skuErr := s.goodsSkuDao.ListByIDs(c, skuIDs)
		if skuErr != nil {
			zap.L().Error("batch get sku list failed", zap.Error(skuErr))
			return nil, errors.New("读取商品库存失败")
		}
		for _, skuInfo := range skuList {
			skuQuantityMap[int64(skuInfo.ID)] = skuInfo.Quantity
		}
	}

	for k, v := range ret {
		ret[k].ImageURL = fileUtils.BuildAbsoluteURL(c, v.ImageURL)
		skuQuantity, ok := skuQuantityMap[v.SkuID]
		if !ok {
			ret[k].IsStockEnough = false
			continue
		}
		ret[k].IsStockEnough = skuQuantity >= v.Quantity
	}
	return ret, nil
}

func (s *IApiShopCartServiceImpl) Remove(c *gin.Context, ids []string) error {
	return s.iShopCartDao.Remove(c, ids)
}
