package impl

import (
	"errors"
	"time"

	"nova-factory-server/app/business/shop/discount/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/constant/commonStatus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DiscountRuleDaoImpl 折扣规则数据访问实现（用于价格计算）
type DiscountRuleDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewDiscountRuleDao 创建折扣规则DAO
func NewDiscountRuleDao(ms *gorm.DB) dao.IDiscountRuleDao {
	return &DiscountRuleDaoImpl{
		db:        ms,
		tableName: "shop_user_discount_rules",
	}
}

func (s *DiscountRuleDaoImpl) GetValidRule(c *gin.Context, userID int64, targetType string, targetID string) (*models.UserDiscountRule, error) {
	var item models.UserDiscountRule
	now := time.Now()
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("state = ?", commonStatus.NORMAL).
		Where("(valid_from IS NULL OR valid_from <= ?)", now).
		Where("(valid_to IS NULL OR valid_to >= ?)", now).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// ListByUserAndGoods 查询用户对商品的折扣规则（优先）和分类折扣规则
func (s *DiscountRuleDaoImpl) ListByUserAndGoods(c *gin.Context, userID int64, goodsID string, categoryID string) (*models.UserDiscountRule, error) {
	// 优先查询商品级折扣
	goodsRule, err := s.GetValidRule(c, userID, "goods", goodsID)
	if err != nil {
		return nil, err
	}
	if goodsRule != nil {
		return goodsRule, nil
	}
	// 其次查询分类级折扣
	if categoryID != "" {
		categoryRule, err := s.GetValidRule(c, userID, "category", categoryID)
		if err != nil {
			return nil, err
		}
		return categoryRule, nil
	}
	return nil, nil
}

// ListValidRulesByTargets 批量查询用户的有效折扣规则（用于列表场景消除 N+1）
func (s *DiscountRuleDaoImpl) ListValidRulesByTargets(c *gin.Context, userID int64, targetType string, targetIDs []string) ([]*models.UserDiscountRule, error) {
	if len(targetIDs) == 0 {
		return nil, nil
	}
	var items []*models.UserDiscountRule
	now := time.Now()
	err := s.db.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("target_type = ?", targetType).
		Where("target_id IN ?", targetIDs).
		Where("state = ?", commonStatus.NORMAL).
		Where("(valid_from IS NULL OR valid_from <= ?)", now).
		Where("(valid_to IS NULL OR valid_to >= ?)", now).
		Find(&items).Error
	return items, err
}
