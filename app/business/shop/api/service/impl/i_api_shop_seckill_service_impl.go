package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// IApiShopSeckillServiceImpl 秒杀服务实现
type IApiShopSeckillServiceImpl struct {
	seckillDao       dao.IApiShopSeckillDao
	seckillConfigDao dao.IApiShopSeckillConfigDao
	shopGoodsService shopservice.IShopGoodsService
}

// NewIApiShopSeckillServiceImpl 创建秒杀服务
func NewIApiShopSeckillServiceImpl(
	seckillDao dao.IApiShopSeckillDao,
	seckillConfigDao dao.IApiShopSeckillConfigDao,
	shopGoodsService shopservice.IShopGoodsService,
) service.IApiShopSeckillService {
	return &IApiShopSeckillServiceImpl{
		seckillDao:       seckillDao,
		seckillConfigDao: seckillConfigDao,
		shopGoodsService: shopGoodsService,
	}
}

// ListConfigs 获取秒杀时间段列表（含当前状态）
func (s *IApiShopSeckillServiceImpl) ListConfigs(c *gin.Context) ([]*models.SeckillConfig, error) {
	data, err := s.seckillConfigDao.List(c, &models.SeckillConfigQuery{
		Status: func() *bool { v := true; return &v }(),
		Page:   1,
		Size:   100,
	})
	if err != nil {
		return nil, err
	}
	// 计算每个时段的当前状态
	for _, cfg := range data.Rows {
		s.computeConfigStatus(cfg)
	}
	return data.Rows, nil
}

// GetCurrentConfig 获取当前秒杀时段配置
func (s *IApiShopSeckillServiceImpl) GetCurrentConfig(c *gin.Context) (*models.SeckillConfig, error) {
	data, err := s.seckillConfigDao.List(c, &models.SeckillConfigQuery{
		Status: func() *bool { v := true; return &v }(),
		Page:   1,
		Size:   100,
	})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	currentMinutes := int64(now.Hour()*60 + now.Minute())

	for _, cfg := range data.Rows {
		beginMinutes := cfg.BeginClock * 60
		endMinutes := (cfg.BeginClock + cfg.ContinueClock) * 60
		if currentMinutes >= beginMinutes && currentMinutes < endMinutes {
			s.computeConfigStatus(cfg)
			return cfg, nil
		}
	}
	return nil, nil
}

// ListGoods 获取秒杀商品列表
func (s *IApiShopSeckillServiceImpl) ListGoods(c *gin.Context, query *models.SeckillQuery) (*models.SeckillListData, error) {
	// 只查询显示中的、未删除的
	query.IsShow = func() *int32 { v := int32(1); return &v }()
	query.Status = func() *int32 { v := int32(1); return &v }()

	data, err := s.seckillDao.List(c, query)
	if err != nil || data == nil {
		return data, err
	}

	// 计算每个商品的进度百分比
	for _, goods := range data.Rows {
		if goods != nil && goods.QuotaShow > 0 {
			goods.Percent = (goods.QuotaShow - goods.Quota) * 100 / goods.QuotaShow
			goods.StockPercent = goods.Percent
		}
	}

	return data, nil
}

// GetGoodsDetail 获取秒杀商品详情
func (s *IApiShopSeckillServiceImpl) GetGoodsDetail(c *gin.Context, id int64, timeID int64) (*models.Seckill, error) {
	goods, err := s.seckillDao.GetByID(c, id)
	if err != nil || goods == nil {
		return goods, err
	}
	applySeckillDetailStatus(c, s.seckillConfigDao, goods, timeID)
	if goods.QuotaShow > 0 {
		goods.StockPercent = (goods.QuotaShow - goods.Quota) * 100 / goods.QuotaShow
		goods.Percent = goods.StockPercent
	}
	goods.Stock = goods.Quota
	if s.shopGoodsService == nil || goods.ProductID == 0 {
		return goods, nil
	}
	product, err := s.shopGoodsService.GetByID(c, goods.ProductID)
	if err != nil || product == nil {
		return goods, err
	}
	goods.GoodsID = product.GoodsID
	goods.VideoURL = product.VideoURL
	goods.Gallery = buildSeckillGallery(goods, product)
	goods.Skus = mapSeckillSkus(goods, product)
	return goods, nil
}

func buildSeckillGallery(seckill *models.Seckill, product *shopmodels.Goods) []string {
	gallery := make([]string, 0)
	appendUniqueMedia(&gallery, seckill.Image)
	for _, item := range splitSeckillMedia(seckill.Images) {
		appendUniqueMedia(&gallery, item)
	}
	if product != nil {
		appendUniqueMedia(&gallery, product.ImageURL)
		for _, item := range product.GalleryImagesArray {
			appendUniqueMedia(&gallery, item)
		}
	}
	return gallery
}

func mapSeckillSkus(seckill *models.Seckill, product *shopmodels.Goods) []*models.SeckillSku {
	if product == nil || len(product.Skus) == 0 {
		return []*models.SeckillSku{}
	}
	rows := make([]*models.SeckillSku, 0, len(product.Skus))
	for _, sku := range product.Skus {
		if sku == nil {
			continue
		}
		gallery := make([]string, 0)
		appendUniqueMedia(&gallery, sku.ImageURL)
		for _, item := range sku.GalleryImagesArray {
			appendUniqueMedia(&gallery, item)
		}
		if len(gallery) == 0 {
			gallery = buildSeckillGallery(seckill, product)
		}
		imageURL := sku.ImageURL
		if imageURL == "" {
			if len(gallery) > 0 {
				imageURL = gallery[0]
			} else if product != nil {
				imageURL = product.ImageURL
			}
		}
		unit := sku.Unit
		if unit == "" && product != nil {
			unit = product.Unit
		}
		rows = append(rows, &models.SeckillSku{
			ID:             sku.ID,
			SkuID:          sku.SkuID,
			SkuName:        sku.SkuName,
			ImageURL:       imageURL,
			GalleryImages:  gallery,
			VideoURL:       chooseMediaURL(sku.VideoURL, product.VideoURL),
			RetailPrice:    seckill.Price,
			OriginalPrice:  sku.RetailPrice,
			Quantity:       sku.Quantity,
			AvailableStock: minSeckillStock(seckill.Quota, sku.Quantity),
			Unit:           unit,
			Weight:         sku.Weight,
			WeightUnit:     sku.WeightUnit,
		})
	}
	return rows
}

func splitSeckillMedia(raw string) []string {
	parts := strings.Split(strings.TrimSpace(raw), ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "[]" || part == "/[]" {
			continue
		}
		result = append(result, part)
	}
	return result
}

func appendUniqueMedia(target *[]string, value string) {
	value = strings.TrimSpace(value)
	if value == "" || value == "[]" || value == "/[]" {
		return
	}
	for _, item := range *target {
		if item == value {
			return
		}
	}
	*target = append(*target, value)
}

func chooseMediaURL(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && value != "[]" && value != "/[]" {
			return value
		}
	}
	return ""
}

func minSeckillStock(values ...int64) int64 {
	if len(values) == 0 {
		return 0
	}
	result := values[0]
	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}
	return result
}

func applySeckillDetailStatus(c *gin.Context, configDao dao.IApiShopSeckillConfigDao, goods *models.Seckill, timeID int64) {
	if goods == nil {
		return
	}

	now := time.Now()
	startTime, startErr := time.ParseInLocation("2006-01-02 15:04:05", goods.StartTime, time.Local)
	stopTime, stopErr := time.ParseInLocation("2006-01-02 15:04:05", goods.StopTime, time.Local)
	if startErr == nil && startTime.After(now) {
		goods.Status = 2
		return
	}
	if stopErr == nil && stopTime.Add(24*time.Hour).Before(now) {
		goods.Status = 0
		return
	}
	if goods.Status != 1 {
		goods.Status = 0
		return
	}
	if timeID == 0 || configDao == nil {
		return
	}
	cfg, err := configDao.GetByID(c, timeID)
	if err != nil || cfg == nil {
		goods.Status = 0
		return
	}
	currentHour := now.Hour()
	startHour := int(cfg.BeginClock)
	endHour := startHour + int(cfg.ContinueClock)
	if startHour <= currentHour && endHour > currentHour {
		goods.Status = 1
		return
	}
	if startHour > currentHour {
		goods.Status = 2
		return
	}
	goods.Status = 0
}

// computeConfigStatus 计算秒杀时段状态
// status: true=抢购中 false=即将开始/已结束
func (s *IApiShopSeckillServiceImpl) computeConfigStatus(cfg *models.SeckillConfig) {
	now := time.Now()
	currentMinutes := int64(now.Hour()*60 + now.Minute())
	beginMinutes := cfg.BeginClock * 60
	endMinutes := (cfg.BeginClock + cfg.ContinueClock) * 60

	if currentMinutes >= beginMinutes && currentMinutes < endMinutes {
		cfg.Status = true
	} else {
		cfg.Status = false
	}
}
