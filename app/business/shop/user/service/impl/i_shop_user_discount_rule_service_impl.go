package impl

import (
	"time"

	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type ShopUserDiscountRuleServiceImpl struct {
	dao dao.IShopUserDiscountRuleDao
}

func NewShopUserDiscountRuleService(dao dao.IShopUserDiscountRuleDao) service.IShopUserDiscountRuleService {
	return &ShopUserDiscountRuleServiceImpl{dao: dao}
}

func (s *ShopUserDiscountRuleServiceImpl) Create(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error) {
	return s.dao.Create(c, req)
}

func (s *ShopUserDiscountRuleServiceImpl) BatchCreate(c *gin.Context, req *models.BatchDiscountRuleCreate) (int64, error) {
	if len(req.UserIDs) == 0 || len(req.TargetIDs) == 0 {
		return 0, nil
	}
	// 限制批量大小
	maxUsers := 100
	maxTargets := 100
	if len(req.UserIDs) > maxUsers {
		req.UserIDs = req.UserIDs[:maxUsers]
	}
	if len(req.TargetIDs) > maxTargets {
		req.TargetIDs = req.TargetIDs[:maxTargets]
	}

	userID := baizeContext.GetUserId(c)
	rules := make([]*models.UserDiscountRule, 0, len(req.UserIDs)*len(req.TargetIDs))
	for _, uid := range req.UserIDs {
		for _, targetID := range req.TargetIDs {
			rules = append(rules, &models.UserDiscountRule{
				ID:           snowflake.GenID(),
				UserID:       uid,
				TargetType:   req.TargetType,
				TargetID:     targetID,
				DiscountRate: req.DiscountRate / 100,
				ValidFrom:    req.ValidFrom,
				ValidTo:      req.ValidTo,
				State:        0,
			})
		}
	}
	for _, rule := range rules {
		rule.SetCreateBy(userID)
	}

	if err := s.dao.BatchCreate(c, rules); err != nil {
		return 0, err
	}
	return int64(len(rules)), nil
}

func (s *ShopUserDiscountRuleServiceImpl) GetByID(c *gin.Context, id int64) (*models.UserDiscountRule, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopUserDiscountRuleServiceImpl) ListByUserID(c *gin.Context, userID int64, page, size int64) (*models.UserDiscountRuleListData, error) {
	return s.dao.ListByUserID(c, userID, page, size)
}

func (s *ShopUserDiscountRuleServiceImpl) ListByUserIDAndType(c *gin.Context, userID int64, targetType string, page, size int64) (*models.UserDiscountRuleListData, error) {
	return s.dao.ListByUserIDAndType(c, userID, targetType, page, size)
}

func (s *ShopUserDiscountRuleServiceImpl) ListByUserIDs(c *gin.Context, userIDs []int64, targetType string, targetIDs []string) ([]*models.UserDiscountRule, error) {
	// 批量查询用户对指定目标的折扣规则
	var rules []*models.UserDiscountRule
	for _, userID := range userIDs {
		for _, targetID := range targetIDs {
			rule, err := s.dao.GetByUserAndTarget(c, userID, targetType, targetID)
			if err != nil {
				return nil, err
			}
			if rule != nil {
				rules = append(rules, rule)
			}
		}
	}
	return rules, nil
}

func (s *ShopUserDiscountRuleServiceImpl) Update(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error) {
	return s.dao.Update(c, req)
}

func (s *ShopUserDiscountRuleServiceImpl) DeleteByID(c *gin.Context, id int64) error {
	return s.dao.DeleteByID(c, id)
}

func (s *ShopUserDiscountRuleServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

// GetValidRule 获取有效的折扣规则（检查有效期）
func (s *ShopUserDiscountRuleServiceImpl) GetValidRule(c *gin.Context, userID int64, targetType string, targetID string) (*models.UserDiscountRule, error) {
	rule, err := s.dao.GetByUserAndTarget(c, userID, targetType, targetID)
	if err != nil || rule == nil {
		return nil, err
	}
	// 检查有效期
	now := time.Now()
	if rule.ValidFrom != nil && now.Before(*rule.ValidFrom) {
		return nil, nil // 规则还未生效
	}
	if rule.ValidTo != nil && now.After(*rule.ValidTo) {
		return nil, nil // 规则已过期
	}
	return rule, nil
}

func (s *ShopUserDiscountRuleServiceImpl) GetUsedTargetIds(c *gin.Context, userID int64) (goodsIds []string, categoryIds []string, err error) {
	return s.dao.GetUsedTargetIdsByUserID(c, userID)
}
