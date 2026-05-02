package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	productdao "nova-factory-server/app/business/shop/product/shopdao"
	shopmodels "nova-factory-server/app/business/shop/product/shopmodels"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const seckillBatchSize = 100

// ShopSeckillActivityServiceImpl 提供秒杀活动及其关联秒杀商品的业务能力。
type ShopSeckillActivityServiceImpl struct {
	dao        dao.IShopSeckillActivityDao
	seckillDao dao.IShopSeckillDao
	goodsDao   productdao.IShopGoodsDao
	skuDao     productdao.IShopSkuDao
}

// NewShopSeckillActivityService 创建秒杀活动服务。
func NewShopSeckillActivityService(
	dao dao.IShopSeckillActivityDao,
	seckillDao dao.IShopSeckillDao,
	goodsDao productdao.IShopGoodsDao,
	skuDao productdao.IShopSkuDao,
) service.IShopSeckillActivityService {
	return &ShopSeckillActivityServiceImpl{
		dao:        dao,
		seckillDao: seckillDao,
		goodsDao:   goodsDao,
		skuDao:     skuDao,
	}
}

// Set 保存秒杀活动，并按 productInfos 同步保存关联的秒杀商品记录。
func (s *ShopSeckillActivityServiceImpl) Set(c *gin.Context, req *models.SeckillActivitySet) (*models.SeckillActivity, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.TimeIDs = splitCSV(normalizeCSV(req.TimeIDs))
	productInfos := compactProductInfos(req.ProductInfos)
	if len(productInfos) == 0 {
		return nil, fmt.Errorf("请选择参与活动的商品")
	}

	goodsMap, err := s.loadActivityGoods(c, productInfos)
	if err != nil {
		return nil, err
	}
	skuMap, err := s.loadActivitySkus(c, goodsMap)
	if err != nil {
		return nil, err
	}
	if err := validateActivityProductAttrs(productInfos, goodsMap, skuMap); err != nil {
		return nil, err
	}

	req.ProductInfos = productInfos
	activity, err := s.dao.Set(c, req)
	if err != nil {
		return nil, err
	}
	if err := s.syncActivitySeckillGoods(c, activity, req, goodsMap, skuMap); err != nil {
		return nil, err
	}
	return activity, nil
}

// DeleteByIDs 删除秒杀活动。
func (s *ShopSeckillActivityServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 根据主键获取秒杀活动。
func (s *ShopSeckillActivityServiceImpl) GetByID(c *gin.Context, id int64) (*models.SeckillActivity, error) {
	activity, err := s.dao.GetByID(c, id)
	if err != nil || activity == nil {
		return activity, err
	}
	productInfos, err := s.loadActivityProductInfos(c, activity.ID)
	if err != nil {
		return nil, err
	}
	activity.ProductInfos = productInfos
	return activity, nil
}

// GetByIDs 根据主键批量获取秒杀活动。
func (s *ShopSeckillActivityServiceImpl) GetByIDs(c *gin.Context, ids []int64) ([]*models.SeckillActivity, error) {
	return s.dao.GetByIDs(c, ids)
}

// List 查询秒杀活动列表。
func (s *ShopSeckillActivityServiceImpl) List(c *gin.Context, req *models.SeckillActivityQuery) (*models.SeckillActivityListData, error) {
	return s.dao.List(c, req)
}

// normalizeCSV 清洗时间段ID数组并转换为逗号分隔字符串。
func normalizeCSV(raw []string) string {
	result := make([]string, 0, len(raw))
	for _, part := range raw {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	return strings.Join(result, ",")
}

// splitCSV 将逗号分隔字符串恢复为数组，供后续 DAO 继续使用数组入参。
func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

// compactProductInfos 清洗活动商品列表，移除空项并规整规格信息。
func compactProductInfos(productInfos []*models.SeckillActivityProductInfo) []*models.SeckillActivityProductInfo {
	result := make([]*models.SeckillActivityProductInfo, 0, len(productInfos))
	for _, productInfo := range productInfos {
		if productInfo == nil || productInfo.ID <= 0 {
			continue
		}
		productInfo.Attrs = compactProductAttrs(productInfo.Attrs)
		result = append(result, productInfo)
	}
	return result
}

// compactProductAttrs 清洗活动商品规格列表。
func compactProductAttrs(attrs []*models.SeckillActivityProductInfoAttr) []*models.SeckillActivityProductInfoAttr {
	result := make([]*models.SeckillActivityProductInfoAttr, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			continue
		}
		attr.SkuID = strings.TrimSpace(attr.SkuID)
		if attr.SkuID == "" {
			continue
		}
		result = append(result, attr)
	}
	return result
}

// loadActivityGoods 加载活动中涉及的商品信息。
func (s *ShopSeckillActivityServiceImpl) loadActivityGoods(
	c *gin.Context,
	productInfos []*models.SeckillActivityProductInfo,
) (map[int64]*shopmodels.Goods, error) {
	goodsMap := make(map[int64]*shopmodels.Goods, len(productInfos))
	for _, productInfo := range productInfos {
		goods, err := s.goodsDao.GetByID(c, productInfo.ID)
		if err != nil {
			return nil, err
		}
		if goods == nil {
			return nil, fmt.Errorf("商品不存在: %d", productInfo.ID)
		}
		goodsMap[productInfo.ID] = goods
	}
	return goodsMap, nil
}

// loadActivitySkus 批量加载活动商品对应的规格信息。
func (s *ShopSeckillActivityServiceImpl) loadActivitySkus(
	c *gin.Context,
	goodsMap map[int64]*shopmodels.Goods,
) (map[int64]map[string]*shopmodels.GoodsSku, error) {
	goodsIDs := make([]string, 0, len(goodsMap))
	goodsIndex := make(map[string]int64, len(goodsMap))
	for id, goods := range goodsMap {
		if goods == nil || strings.TrimSpace(goods.GoodsID) == "" {
			continue
		}
		goodsIDs = append(goodsIDs, goods.GoodsID)
		goodsIndex[goods.GoodsID] = id
	}

	skuRows, err := s.skuDao.ListByGoodsIDs(c, goodsIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]map[string]*shopmodels.GoodsSku, len(goodsMap))
	for _, sku := range skuRows {
		if sku == nil {
			continue
		}
		productID, ok := goodsIndex[sku.GoodsID]
		if !ok {
			continue
		}
		if result[productID] == nil {
			result[productID] = make(map[string]*shopmodels.GoodsSku)
		}
		result[productID][sku.SkuID] = sku
	}
	return result, nil
}

// validateActivityProductAttrs 校验活动商品规格配置是否有效。
func validateActivityProductAttrs(
	productInfos []*models.SeckillActivityProductInfo,
	goodsMap map[int64]*shopmodels.Goods,
	skuMap map[int64]map[string]*shopmodels.GoodsSku,
) error {
	for _, productInfo := range productInfos {
		goods := goodsMap[productInfo.ID]
		if goods == nil {
			return fmt.Errorf("商品不存在: %d", productInfo.ID)
		}
		if len(productInfo.Attrs) == 0 {
			continue
		}

		productSkus := skuMap[productInfo.ID]
		if len(productSkus) == 0 {
			return fmt.Errorf("商品没有可用规格: %d", productInfo.ID)
		}
		for _, attr := range productInfo.Attrs {
			if attr == nil || attr.Status == 0 {
				continue
			}
			sku := productSkus[attr.SkuID]
			if sku == nil {
				return fmt.Errorf("商品规格不存在: %s", attr.SkuID)
			}
			if attr.Price <= 0 {
				return fmt.Errorf("请填写商品规格活动价: %s", attr.SkuID)
			}
			if attr.Quota < 0 {
				return fmt.Errorf("请填写商品规格限量: %s", attr.SkuID)
			}
			if attr.Quota > sku.Quantity {
				return fmt.Errorf("商品规格限量超过库存: %s", attr.SkuID)
			}
		}
	}
	return nil
}

// syncActivitySeckillGoods 按活动商品列表同步保存秒杀商品。
func (s *ShopSeckillActivityServiceImpl) syncActivitySeckillGoods(
	c *gin.Context,
	activity *models.SeckillActivity,
	req *models.SeckillActivitySet,
	goodsMap map[int64]*shopmodels.Goods,
	skuMap map[int64]map[string]*shopmodels.GoodsSku,
) error {
	currentRows, err := s.findSeckillsByActivityID(c, activity.ID)
	if err != nil {
		return err
	}
	currentMap := make(map[int64]*models.Seckill, len(currentRows))
	for _, row := range currentRows {
		if row == nil {
			continue
		}
		currentMap[row.ProductID] = row
	}

	keepProductIDs := make(map[int64]struct{}, len(req.ProductInfos))
	createReqs := make([]*models.SeckillSet, 0, len(req.ProductInfos))
	updateReqs := make([]*models.SeckillSet, 0, len(req.ProductInfos))
	for _, productInfo := range req.ProductInfos {
		goods := goodsMap[productInfo.ID]
		if goods == nil {
			return fmt.Errorf("商品不存在: %d", productInfo.ID)
		}
		seckillReq := buildSeckillSetFromActivity(
			activity,
			req,
			productInfo,
			goods,
			skuMap[productInfo.ID],
			currentMap[productInfo.ID],
		)
		if seckillReq.ID > 0 {
			updateReqs = append(updateReqs, seckillReq)
		} else {
			createReqs = append(createReqs, seckillReq)
		}
		keepProductIDs[productInfo.ID] = struct{}{}
	}

	if len(createReqs) == 0 && len(updateReqs) == 0 {
		return errors.New("不存在要更新和写入的秒杀商品")
	}

	if len(createReqs) != 0 {
		if err := s.seckillDao.BatchCreate(c, createReqs, seckillBatchSize); err != nil {
			return err
		}
	}

	if len(updateReqs) != 0 {
		if err := s.seckillDao.BatchUpdate(c, updateReqs, seckillBatchSize); err != nil {
			return err
		}
	}

	removeIDs := make([]int64, 0)
	for productID, current := range currentMap {
		if _, ok := keepProductIDs[productID]; ok {
			continue
		}
		removeIDs = append(removeIDs, current.ID)
	}
	if len(removeIDs) > 0 {
		return s.seckillDao.DeleteByIDs(c, removeIDs)
	}
	return nil
}

// findSeckillsByActivityID 查询活动当前关联的秒杀商品列表。
func (s *ShopSeckillActivityServiceImpl) findSeckillsByActivityID(c *gin.Context, activityID int64) ([]*models.Seckill, error) {
	data, err := s.seckillDao.List(c, &models.SeckillQuery{
		ActivityID: activityID,
		Page:       1,
		Size:       1000,
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return []*models.Seckill{}, nil
	}
	return data.Rows, nil
}

// loadActivityProductInfos 从秒杀商品记录中还原活动商品列表。
func (s *ShopSeckillActivityServiceImpl) loadActivityProductInfos(c *gin.Context, activityID int64) ([]*models.SeckillActivityProductInfo, error) {
	rows, err := s.findSeckillsByActivityID(c, activityID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []*models.SeckillActivityProductInfo{}, nil
	}
	productInfos := make([]*models.SeckillActivityProductInfo, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		productInfos = append(productInfos, &models.SeckillActivityProductInfo{
			ID:     row.ProductID,
			Status: row.Status,
			Sort:   row.Sort,
			IsHot:  row.IsHot,
			Attrs:  []*models.SeckillActivityProductInfoAttr{},
		})
	}
	return productInfos, nil
}

// buildSeckillSetFromActivity 将活动和商品信息组装为秒杀商品保存参数。
func buildSeckillSetFromActivity(
	activity *models.SeckillActivity,
	req *models.SeckillActivitySet,
	productInfo *models.SeckillActivityProductInfo,
	goods *shopmodels.Goods,
	productSkus map[string]*shopmodels.GoodsSku,
	current *models.Seckill,
) *models.SeckillSet {
	images := normalizeGoodsGalleryImages(goods.GalleryImages)
	image := strings.TrimSpace(goods.ImageURL)
	if image == "" && len(images) > 0 {
		image = images[0]
	}

	title := strings.TrimSpace(activity.Title)
	if title == "" {
		title = strings.TrimSpace(goods.GoodsName)
	}

	stats := buildSeckillStats(goods, productInfo, productSkus)
	seckillReq := &models.SeckillSet{
		ActivityID:   activity.ID,
		ProductID:    productInfo.ID,
		Image:        image,
		Images:       strings.Join(images, ","),
		Title:        title,
		Info:         strings.TrimSpace(goods.Description),
		Price:        stats.Price,
		Cost:         stats.Cost,
		OtPrice:      stats.OtPrice,
		Sort:         productInfo.Sort,
		Stock:        stats.Stock,
		Sales:        0,
		UnitName:     strings.TrimSpace(goods.Unit),
		Postage:      0,
		StartTime:    formatActivityBoundary(activity.StartDay, false),
		StopTime:     formatActivityBoundary(activity.EndDay, true),
		Status:       productInfo.Status,
		IsPostage:    0,
		IsHot:        productInfo.IsHot,
		Num:          clampInt64ToInt32(activity.Num),
		IsShow:       int32(activity.Status),
		TimeID:       activity.TimeIDs,
		TempID:       0,
		Weight:       goods.Weight,
		Volume:       0,
		Quota:        stats.Quota,
		QuotaShow:    stats.Quota,
		OnceNum:      clampInt64ToInt32(activity.OnceNum),
		Logistics:    "1,2",
		Freight:      2,
		CustomForm:   "",
		VirtualType:  0,
		IsCommission: activity.IsCommission,
	}
	if current != nil {
		seckillReq.ID = current.ID
		seckillReq.Sales = current.Sales
		if current.Postage > 0 {
			seckillReq.Postage = current.Postage
		}
		if current.TempID > 0 {
			seckillReq.TempID = current.TempID
		}
		if current.Volume > 0 {
			seckillReq.Volume = current.Volume
		}
		if current.Logistics != "" {
			seckillReq.Logistics = current.Logistics
		}
		if current.Freight > 0 {
			seckillReq.Freight = current.Freight
		}
		if current.CustomForm != "" {
			seckillReq.CustomForm = current.CustomForm
		}
		if current.VirtualType > 0 {
			seckillReq.VirtualType = current.VirtualType
		}
	}
	return seckillReq
}

type seckillStats struct {
	Price   float64
	Cost    float64
	OtPrice float64
	Stock   int64
	Quota   int64
}

// buildSeckillStats 根据商品活动规格汇总出秒杀商品主表需要的价格和库存数据。
func buildSeckillStats(
	goods *shopmodels.Goods,
	productInfo *models.SeckillActivityProductInfo,
	productSkus map[string]*shopmodels.GoodsSku,
) seckillStats {
	result := seckillStats{
		Price:   goods.RetailPrice,
		Cost:    0,
		OtPrice: goods.RetailPrice,
		Stock:   goods.Quantity,
		Quota:   goods.Quantity,
	}

	activeAttrs := collectActiveAttrs(productInfo.Attrs)
	if len(activeAttrs) == 0 {
		return result
	}

	minPrice := math.MaxFloat64
	minCost := math.MaxFloat64
	minOtPrice := math.MaxFloat64
	var totalQuota int64
	for _, attr := range activeAttrs {
		if attr.Price < minPrice {
			minPrice = attr.Price
		}
		if attr.Cost >= 0 && attr.Cost < minCost {
			minCost = attr.Cost
		}
		if attr.OtPrice >= 0 && attr.OtPrice < minOtPrice {
			minOtPrice = attr.OtPrice
		}
		totalQuota += attr.Quota
		if sku := productSkus[attr.SkuID]; sku != nil && sku.Quantity > 0 {
			result.Stock += 0
		}
	}
	result.Price = minPrice
	if minCost != math.MaxFloat64 {
		result.Cost = minCost
	}
	if minOtPrice != math.MaxFloat64 {
		result.OtPrice = minOtPrice
	}
	result.Stock = totalQuota
	result.Quota = totalQuota
	return result
}

// collectActiveAttrs 过滤参与活动的规格。
func collectActiveAttrs(attrs []*models.SeckillActivityProductInfoAttr) []*models.SeckillActivityProductInfoAttr {
	result := make([]*models.SeckillActivityProductInfoAttr, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil || attr.Status == 0 {
			continue
		}
		result = append(result, attr)
	}
	return result
}

// normalizeGoodsGalleryImages 兼容 JSON 数组或逗号分隔格式的商品图集字段。
func normalizeGoodsGalleryImages(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}

	var jsonImages []string
	if strings.HasPrefix(raw, "[") && json.Unmarshal([]byte(raw), &jsonImages) == nil {
		return compactStrings(jsonImages)
	}
	return compactStrings(strings.Split(raw, ","))
}

// compactStrings 清洗字符串数组，移除空值和首尾空格。
func compactStrings(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}

// formatActivityBoundary 将活动日期转换为秒杀商品使用的开始或结束时间字符串。
func formatActivityBoundary(day int32, endOfDay bool) string {
	raw := strconv.FormatInt(int64(day), 10)
	if len(raw) == 8 {
		if parsed, err := time.ParseInLocation("20060102", raw, time.Local); err == nil {
			if endOfDay {
				return parsed.Format("2006-01-02") + " 23:59:59"
			}
			return parsed.Format("2006-01-02") + " 00:00:00"
		}
	}

	parsed := time.Unix(int64(day), 0).In(time.Local)
	if endOfDay {
		return parsed.Format("2006-01-02") + " 23:59:59"
	}
	return parsed.Format("2006-01-02") + " 00:00:00"
}

// clampInt64ToInt32 将 int64 安全转换为 int32。
func clampInt64ToInt32(value int64) int32 {
	if value > 2147483647 {
		return 2147483647
	}
	if value < -2147483648 {
		return -2147483648
	}
	return int32(value)
}
