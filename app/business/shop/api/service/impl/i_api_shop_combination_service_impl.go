package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// IApiShopCombinationServiceImpl 拼团服务实现
type IApiShopCombinationServiceImpl struct {
	combinationDao   dao.IApiShopCombinationDao
	pinkDao          dao.IApiShopPinkDao
	shopGoodsService shopservice.IShopGoodsService
}

// NewIApiShopCombinationServiceImpl 创建拼团服务
func NewIApiShopCombinationServiceImpl(
	combinationDao dao.IApiShopCombinationDao,
	pinkDao dao.IApiShopPinkDao,
	shopGoodsService shopservice.IShopGoodsService,
) service.IApiShopCombinationService {
	return &IApiShopCombinationServiceImpl{
		combinationDao:   combinationDao,
		pinkDao:          pinkDao,
		shopGoodsService: shopGoodsService,
	}
}

// List 获取拼团商品列表
func (s *IApiShopCombinationServiceImpl) List(c *gin.Context, query *models.CombinationQuery) (*models.CombinationListData, error) {
	// 只查询显示中的、未删除的
	query.IsShow = func() *int32 { v := int32(1); return &v }()
	return s.combinationDao.List(c, query)
}

// GetByID 获取拼团商品详情（含SKU和进行中的团）
func (s *IApiShopCombinationServiceImpl) GetByID(c *gin.Context, id int64) (*models.Combination, error) {
	combination, err := s.combinationDao.GetByID(c, id)
	if err != nil || combination == nil {
		return combination, err
	}

	// 填充商品扩展信息（gallery、videoUrl、skus）
	s.applyCombinationDetail(c, combination)
	return combination, nil
}

func (s *IApiShopCombinationServiceImpl) applyCombinationDetail(c *gin.Context, combination *models.Combination) {
	if s.shopGoodsService == nil {
		return
	}
	goodsKey := strings.TrimSpace(combination.ProductID)
	if goodsKey == "" {
		return
	}

	product, err := s.shopGoodsService.GetByGoodsID(c, goodsKey)
	if err != nil || product == nil {
		productIDInt, parseErr := strconv.ParseInt(goodsKey, 10, 64)
		if parseErr != nil {
			return
		}
		product, err = s.shopGoodsService.GetByID(c, productIDInt)
		if err != nil || product == nil {
			return
		}
	}

	combination.ProductID = strconv.FormatInt(product.ID, 10)
	combination.GoodsID = product.GoodsID
	combination.GoodsName = product.GoodsName
	combination.VideoURL = product.VideoURL
	combination.Gallery = buildCombinationGallery(combination, product)
	combination.Skus = mapCombinationSkus(combination, product)
}

func buildCombinationGallery(combination *models.Combination, product *shopmodels.Goods) []string {
	gallery := make([]string, 0)
	combinationAppendUniqueMedia(&gallery, combination.Image)
	for _, item := range splitCombinationMedia(combination.Images) {
		combinationAppendUniqueMedia(&gallery, item)
	}
	if product != nil {
		combinationAppendUniqueMedia(&gallery, product.ImageURL)
		for _, item := range product.GalleryImagesArray {
			combinationAppendUniqueMedia(&gallery, item)
		}
	}
	return gallery
}

func mapCombinationSkus(combination *models.Combination, product *shopmodels.Goods) []*models.CombinationSku {
	if product == nil || len(product.Skus) == 0 {
		return []*models.CombinationSku{}
	}
	rows := make([]*models.CombinationSku, 0, len(product.Skus))
	for _, sku := range product.Skus {
		if sku == nil {
			continue
		}
		gallery := make([]string, 0)
		combinationAppendUniqueMedia(&gallery, sku.ImageURL)
		for _, item := range sku.GalleryImagesArray {
			combinationAppendUniqueMedia(&gallery, item)
		}
		if len(gallery) == 0 {
			gallery = buildCombinationGallery(combination, product)
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
		rows = append(rows, &models.CombinationSku{
			ID:             sku.ID,
			SkuID:          sku.SkuID,
			SkuName:        sku.SkuName,
			ImageURL:       imageURL,
			GalleryImages:  gallery,
			VideoURL:       chooseCombinationMediaURL(sku.VideoURL, product.VideoURL),
			Price:          combination.Price, // 拼团价用活动的
			OriginalPrice:  sku.RetailPrice,   // 单买价用SKU原价
			Quantity:       sku.Quantity,
			AvailableStock: minCombinationStock(combination.Quota, sku.Quantity),
			Unit:           unit,
			Weight:         sku.Weight,
			WeightUnit:     sku.WeightUnit,
		})
	}
	return rows
}

func splitCombinationMedia(raw string) []string {
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

func combinationAppendUniqueMedia(target *[]string, value string) {
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

func chooseCombinationMediaURL(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && value != "[]" && value != "/[]" {
			return value
		}
	}
	return ""
}

func minCombinationStock(values ...int64) int64 {
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

// GetPinkList 获取正在进行中的团列表
func (s *IApiShopCombinationServiceImpl) GetPinkList(c *gin.Context, cid int64) (*models.PinkListData, error) {
	status := int32(1)
	isRefund := int32(0)
	return s.pinkDao.List(c, &models.PinkQuery{
		CID:      cid,
		Status:   &status,
		IsRefund: &isRefund,
		Page:     1,
		Size:     100,
	})
}

// IApiShopPinkServiceImpl 拼团记录服务实现
type IApiShopPinkServiceImpl struct {
	pinkDao dao.IApiShopPinkDao
}

// NewIApiShopPinkServiceImpl 创建拼团记录服务
func NewIApiShopPinkServiceImpl(pinkDao dao.IApiShopPinkDao) service.IApiShopPinkService {
	return &IApiShopPinkServiceImpl{
		pinkDao: pinkDao,
	}
}

// Create 创建拼团记录（开团或参团）
// 注意：实际开团由订单支付成功回调触发，这里只做参数校验
func (s *IApiShopPinkServiceImpl) Create(c *gin.Context, userID int64, combinationID int64, pinkID int64, orderID int64) (*models.Pink, error) {
	// 参数校验由 Controller 层处理
	// 实际拼团记录创建在订单支付回调中处理
	return nil, nil
}

// GetDetail 获取团详情（含当前人数）
func (s *IApiShopPinkServiceImpl) GetDetail(c *gin.Context, pinkID int64) (*models.Pink, error) {
	return s.pinkDao.GetByID(c, pinkID)
}

// GetPinkMemberCount 获取团内人数
func (s *IApiShopPinkServiceImpl) GetPinkMemberCount(c *gin.Context, pinkID int64) (int64, error) {
	return s.pinkDao.CountMembers(c, pinkID)
}

// ListMyPinks 获取用户的拼团记录
func (s *IApiShopPinkServiceImpl) ListMyPinks(c *gin.Context, userID int64) (*models.PinkListData, error) {
	return s.pinkDao.List(c, &models.PinkQuery{
		UID: userID,
	})
}
