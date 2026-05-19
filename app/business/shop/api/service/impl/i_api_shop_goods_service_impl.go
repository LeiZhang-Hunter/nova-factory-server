package impl

import (
	"errors"
	"math"
	"strconv"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// IApiShopGoodsServiceImpl 商品服务实现
type IApiShopGoodsServiceImpl struct {
	dao              dao.IApiShopGoodsDao
	shopGoodsService shopservice.IShopGoodsService
	discountService  discountservice.IDiscountCalculateService
}

// NewIApiShopGoodsServiceImpl  创建商品服务
func NewIApiShopGoodsServiceImpl(dao dao.IApiShopGoodsDao, shopGoodsService shopservice.IShopGoodsService, discountService discountservice.IDiscountCalculateService) service.IApiShopGoodsService {
	return &IApiShopGoodsServiceImpl{
		dao:              dao,
		shopGoodsService: shopGoodsService,
		discountService:  discountService,
	}
}

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
	data, err := s.dao.ListByUserPurchased(c, userID, query.CategoryId, page, size)
	if err != nil || data == nil || len(data.Rows) == 0 {
		return data, err
	}
	// 应用折扣价格（商品级别，无SKU，批量查询消除 N+1）
	s.applyDiscountPriceForList(c, data.Rows)
	return data, nil
}

// Search 按多个商品名称检索相似商品，并回填数据库中的最新商品数据
func (s *IApiShopGoodsServiceImpl) Search(c *gin.Context, req *models.GoodsSearchReq) (*models.GoodsSearchData, error) {
	if req == nil {
		return nil, errors.New("检索参数不能为空")
	}
	if len(req.GoodsNames) == 0 {
		return nil, errors.New("商品名称不能为空")
	}
	if req.Embedding == nil {
		return nil, errors.New("向量模型配置不能为空")
	}

	limit := normalizeGoodsSearchLimit(req.Limit)
	vectorData, err := s.shopGoodsService.BatchSearchVector(c, &shopmodels.GoodsVectorBatchSearchReq{
		Queries:   req.GoodsNames,
		Limit:     buildGoodsVectorSearchLimit(limit),
		Embedding: req.Embedding,
	})
	if err != nil {
		return nil, err
	}
	if vectorData == nil || len(vectorData.Rows) == 0 {
		return &models.GoodsSearchData{
			Rows:  make([]*models.GoodsSearchItem, 0),
			Total: 0,
		}, nil
	}

	goodsMap, err := s.loadSearchGoodsMap(c, vectorData.Rows)
	if err != nil {
		return nil, err
	}

	items := make([]*models.GoodsSearchItem, 0, len(vectorData.Rows))
	for _, row := range vectorData.Rows {
		if row == nil {
			continue
		}
		matches := make([]*models.GoodsSearchMatch, 0, limit)
		seen := make(map[int64]struct{}, len(row.Rows))
		for _, hit := range row.Rows {
			if hit == nil || hit.GoodsDBID == 0 {
				continue
			}
			if _, ok := seen[hit.GoodsDBID]; ok {
				continue
			}
			goods, ok := goodsMap[hit.GoodsDBID]
			if !ok || goods == nil {
				continue
			}
			seen[hit.GoodsDBID] = struct{}{}
			matches = append(matches, &models.GoodsSearchMatch{
				Score: hit.Score,
				Goods: goods,
			})
			if len(matches) >= limit {
				break
			}
		}
		items = append(items, &models.GoodsSearchItem{
			Query: row.Query,
			Rows:  matches,
			Total: int64(len(matches)),
		})
	}

	return &models.GoodsSearchData{
		Rows:  items,
		Total: int64(len(items)),
	}, nil
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
	seen := make(map[int64]struct{})
	for _, row := range rows {
		if row == nil {
			continue
		}
		for _, item := range row.Rows {
			if item == nil || item.GoodsDBID == 0 {
				continue
			}
			if _, ok := seen[item.GoodsDBID]; ok {
				continue
			}
			seen[item.GoodsDBID] = struct{}{}
			ids = append(ids, item.GoodsDBID)
		}
	}
	if len(ids) == 0 {
		return map[int64]*models.Goods{}, nil
	}

	goodsRows, err := s.dao.ListByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	s.applyDiscountPriceForList(c, goodsRows)

	goodsMap := make(map[int64]*models.Goods, len(goodsRows))
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
