package impl

import (
	"math"
	"strconv"

	"nova-factory-server/app/business/shop/api/models"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// GetByID 获取商品详情
func (s *IApiShopGoodsServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
	goods, err := s.shopGoodsService.GetByID(c, id)
	if err != nil || goods == nil {
		return goods, err
	}
	s.applyDiscountPrice(c, goods)
	return goods, nil
}

// List 获取商品列表
func (s *IApiShopGoodsServiceImpl) List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error) {
	data, err := s.dao.List(c, query)
	if err != nil || data == nil || len(data.Rows) == 0 {
		return data, err
	}
	// 应用折扣价格（商品级别，无SKU，批量查询消除 N+1）
	s.applyDiscountPriceForList(c, data.Rows)
	return data, nil
}

// ListRepurchase 获取用户复购商品列表
func (s *IApiShopGoodsServiceImpl) ListRepurchase(c *gin.Context, userID int64, query *models.GoodsQuery) (*models.GoodsListData, error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	size := query.Size
	if size <= 0 {
		size = 20
	}
	query.Page = page
	query.Size = size
	data, err := s.dao.ListByUserPurchased(c, userID, query)
	if err != nil || data == nil || len(data.Rows) == 0 {
		return data, err
	}
	// 应用折扣价格（商品级别，无SKU，批量查询消除 N+1）
	s.applyDiscountPriceForList(c, data.Rows)
	return data, nil
}

// applyDiscountPrice 应用折扣价格（完整版，包含SKU）
func (s *IApiShopGoodsServiceImpl) applyDiscountPrice(c *gin.Context, goods *shopmodels.Goods) {
	if goods == nil || s.discountService == nil {
		return
	}
	userID := baizeContext.GetUserIdSafe(c)
	if userID == 0 {
		return
	}

	// 计算商品级别的折扣价
	if len(goods.Skus) > 0 {
		skuPrices := make([]*discountservice.SkuPrice, 0, len(goods.Skus))
		for _, sku := range goods.Skus {
			skuPrices = append(skuPrices, &discountservice.SkuPrice{
				ID:          int64(sku.ID),
				SkuID:       sku.SkuID,
				RetailPrice: sku.RetailPrice,
			})
		}
		discountedSkus, err := s.discountService.CalculateSkuDiscountPrice(c, userID, strconv.FormatInt(goods.ShopCategoryId, 10), skuPrices)
		if err == nil && discountedSkus != nil {
			skuMap := make(map[int64]float64, len(discountedSkus))
			for _, ds := range discountedSkus {
				skuMap[ds.ID] = ds.RetailPrice
			}
			for _, sku := range goods.Skus {
				if discountedPrice, ok := skuMap[int64(sku.ID)]; ok {
					sku.RetailPrice = discountedPrice
				}
			}
		}
	} else if goods.RetailPrice > 0 {
		// 没有SKU的情况，直接计算商品零售价的折扣
		discountedPrice, hasDiscount := s.discountService.CalculateDiscountPrice(c, userID, goods.GoodsID, "", strconv.FormatInt(goods.ShopCategoryId, 10), goods.RetailPrice)
		if hasDiscount {
			goods.RetailPrice = math.Round(discountedPrice*100) / 100
		}
	}
}

// applyDiscountPriceSimple 应用折扣价格（简洁版，仅商品级别，无SKU）
func (s *IApiShopGoodsServiceImpl) applyDiscountPriceSimple(c *gin.Context, goods *models.Goods) {
	if goods == nil || s.discountService == nil {
		return
	}
	userID := baizeContext.GetUserIdSafe(c)
	if userID == 0 {
		return
	}
	if goods.RetailPrice > 0 {
		discountedPrice, hasDiscount := s.discountService.CalculateDiscountPrice(c, userID, goods.GoodsID, "", strconv.FormatInt(goods.ShopCategoryId, 10), goods.RetailPrice)
		if hasDiscount {
			goods.RetailPrice = math.Round(discountedPrice*100) / 100
		}
	}
}

// applyDiscountPriceForList 批量应用折扣价格（用于列表场景，消除 N+1 查询）
func (s *IApiShopGoodsServiceImpl) applyDiscountPriceForList(c *gin.Context, goodsList []*models.Goods) {
	if len(goodsList) == 0 || s.discountService == nil {
		return
	}
	userID := baizeContext.GetUserIdSafe(c)
	if userID == 0 {
		return
	}
	// 构建批量输入
	batchGoods := make([]*discountservice.GoodsWithPrice, 0, len(goodsList))
	for _, goods := range goodsList {
		if goods != nil && goods.RetailPrice > 0 {
			// 列表场景暂无 SKU 信息，先传空数组
			batchGoods = append(batchGoods, &discountservice.GoodsWithPrice{
				GoodsID:     goods.GoodsID,
				SkuIDs:      []string{},
				CategoryID:  strconv.FormatInt(goods.ShopCategoryId, 10),
				RetailPrice: goods.RetailPrice,
			})
		}
	}
	if len(batchGoods) == 0 {
		return
	}
	// 一次批量查询 → 内存 Map
	priceMap := s.discountService.BatchCalculateDiscountPrices(c, userID, batchGoods)
	if len(priceMap) == 0 {
		return
	}
	// 应用折扣价
	for _, goods := range goodsList {
		if goods == nil {
			continue
		}
		if discountedPrice, ok := priceMap[goods.GoodsID]; ok {
			goods.RetailPrice = discountedPrice
		}
	}
}

func (s *IApiShopGoodsServiceImpl) loadSearchGoodsMap(c *gin.Context, rows []*shopmodels.GoodsVectorBatchSearchItem) (map[int64]*models.Goods, error) {
	ids := make([]int64, 0)
	names := make([]string, 0)
	seen := make(map[int64]struct{})
	for _, row := range rows {
		if row == nil {
			continue
		}
		for _, item := range row.Rows {
			if item == nil || item.SkuID == 0 {
				continue
			}
			if _, ok := seen[item.SkuID]; ok {
				continue
			}
			seen[item.SkuID] = struct{}{}
			ids = append(ids, item.SkuID)
			names = append(names, item.GoodsName)
		}
	}
	if len(ids) == 0 {
		return map[int64]*models.Goods{}, nil
	}

	skuRows, err := s.skuDao.ListByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	goodsRows := make([]*models.Goods, 0)

	for k, v := range skuRows {
		var goods models.Goods
		goods.ID = int64(v.ID)
		goods.GoodsID = v.GoodsID
		goods.GoodsName = names[k]
		//goods.GoodsName = v.SkuName
		goods.GoodsCode = v.GoodsID
		goods.OuterID = v.OuterID
		goods.Description = v.Description
		goods.Weight = v.Weight
		goods.WeightUnit = v.WeightUnit
		goods.Unit = v.Unit
		goods.IsOnSale = 1
		goods.Quantity = v.Quantity
		goodsRows = append(goodsRows, &goods)
	}
	s.applyDiscountPriceForList(c, goodsRows)

	goodsMap := make(map[int64]*models.Goods, len(skuRows))
	for _, goods := range goodsRows {
		if goods == nil {
			continue
		}
		goodsMap[goods.ID] = goods
	}

	return goodsMap, nil
}

func normalizeGoodsSearchLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func buildGoodsVectorSearchLimit(limit int) int {
	vectorLimit := limit * 3
	if vectorLimit < limit {
		return limit
	}
	if vectorLimit > 50 {
		return 50
	}
	return vectorLimit
}
