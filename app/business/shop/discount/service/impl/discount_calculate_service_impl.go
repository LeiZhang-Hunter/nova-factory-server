package impl

import (
	"math"

	"nova-factory-server/app/business/shop/discount/dao"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// DiscountCalculateServiceImpl 折扣价格计算服务实现
type DiscountCalculateServiceImpl struct {
	dao dao.IDiscountRuleDao
}

// NewDiscountCalculateService 创建折扣计算服务
func NewDiscountCalculateService(dao dao.IDiscountRuleDao) discountservice.IDiscountCalculateService {
	return &DiscountCalculateServiceImpl{dao: dao}
}

// CalculateDiscountPrice 计算商品折扣价
func (s *DiscountCalculateServiceImpl) CalculateDiscountPrice(c *gin.Context, userID int64, goodsID int64, skuID int64, categoryID int64, originalPrice float64) (float64, bool) {
	if userID == 0 || originalPrice <= 0 {
		return originalPrice, false
	}

	// 优先查 SKU 级折扣
	if skuID != 0 {
		skuRule, err := s.dao.GetValidRule(c, userID, "sku", skuID)
		if err == nil && skuRule != nil {
			// discountRate 直接使用（存储的是小数，如 0.50 表示 50%）
			discountedPrice := math.Round(originalPrice*skuRule.DiscountRate*100) / 100
			return discountedPrice, true
		}
	}

	// 其次查分类级折扣
	if categoryID != 0 {
		categoryRule, err := s.dao.GetValidRule(c, userID, "category", categoryID)
		if err == nil && categoryRule != nil {
			discountedPrice := math.Round(originalPrice*categoryRule.DiscountRate*100) / 100
			return discountedPrice, true
		}
	}

	return originalPrice, false
}

// CalculateSkuDiscountPrice 计算SKU的折扣价
func (s *DiscountCalculateServiceImpl) CalculateSkuDiscountPrice(c *gin.Context, userID int64, categoryID int64, skus []*discountservice.SkuPrice) ([]*discountservice.SkuPrice, error) {
	if len(skus) == 0 || userID == 0 {
		return skus, nil
	}

	// 构建 SKU ID 列表
	skuIDs := make([]int64, 0, len(skus))
	for _, sku := range skus {
		if sku != nil && sku.SkuID != 0 {
			skuIDs = append(skuIDs, sku.SkuID)
		}
	}

	// 1. 批量查询 SKU 规则
	skuRules, _ := s.dao.ListValidRulesByTargets(c, userID, "sku", skuIDs)
	skuRuleMap := make(map[int64]*models.UserDiscountRule, len(skuRules))
	for _, r := range skuRules {
		skuRuleMap[r.TargetID] = r
	}

	// 2. 查分类规则（兜底）
	var categoryRule *models.UserDiscountRule
	if categoryID != 0 {
		categoryRule, _ = s.dao.GetValidRule(c, userID, "category", categoryID)
	}

	// 3. 逐个 SKU 计算折扣价
	for _, sku := range skus {
		if sku == nil {
			continue
		}
		if rule, ok := skuRuleMap[sku.SkuID]; ok {
			// SKU 专属折扣优先
			sku.RetailPrice = math.Round(sku.RetailPrice*rule.DiscountRate*100) / 100
		} else if categoryRule != nil {
			// 分类折扣兜底
			sku.RetailPrice = math.Round(sku.RetailPrice*categoryRule.DiscountRate*100) / 100
		}
	}

	return skus, nil
}

// BatchCalculateDiscountPrices 批量计算商品折扣价（内部批量查库一次，用于列表场景消除 N+1）
func (s *DiscountCalculateServiceImpl) BatchCalculateDiscountPrices(c *gin.Context, userID int64, goods []*discountservice.GoodsWithPrice) map[int64]float64 {
	if len(goods) == 0 || userID == 0 {
		return nil
	}

	// 1. 收集所有 SKU ID 和分类 ID
	skuIDSet := make(map[int64]struct{})
	catIDSet := make(map[int64]struct{})
	for _, g := range goods {
		for _, skuID := range g.SkuIDs {
			if skuID != 0 {
				skuIDSet[skuID] = struct{}{}
			}
		}
		if g.CategoryID != 0 {
			catIDSet[g.CategoryID] = struct{}{}
		}
	}
	skuIDs := make([]int64, 0, len(skuIDSet))
	for id := range skuIDSet {
		skuIDs = append(skuIDs, id)
	}
	catIDs := make([]int64, 0, len(catIDSet))
	for id := range catIDSet {
		catIDs = append(catIDs, id)
	}

	// 2. 批量查询（2次 DB）
	skuRules, _ := s.dao.ListValidRulesByTargets(c, userID, "sku", skuIDs)
	catRules, _ := s.dao.ListValidRulesByTargets(c, userID, "category", catIDs)

	// 3. 构建内存 Map
	skuRuleMap := make(map[int64]*models.UserDiscountRule, len(skuRules))
	for _, r := range skuRules {
		skuRuleMap[r.TargetID] = r
	}
	catRuleMap := make(map[int64]*models.UserDiscountRule, len(catRules))
	for _, r := range catRules {
		catRuleMap[r.TargetID] = r
	}

	// 4. 内存中计算折扣价（按 SKU 级别）
	result := make(map[int64]float64, len(goods))
	for _, g := range goods {
		var rate float64 = 1.0
		found := false

		// 优先匹配 SKU 折扣
		for _, skuID := range g.SkuIDs {
			if rule, ok := skuRuleMap[skuID]; ok {
				rate = rule.DiscountRate
				found = true
				break
			}
		}

		// SKU 没有匹配，再匹配分类折扣
		if !found && g.CategoryID != 0 {
			if rule, ok := catRuleMap[g.CategoryID]; ok {
				rate = rule.DiscountRate
				found = true
			}
		}

		if found && rate > 0 && rate < 1.0 {
			result[g.GoodsID] = math.Round(g.RetailPrice*rate*100) / 100
		}
	}
	return result
}
