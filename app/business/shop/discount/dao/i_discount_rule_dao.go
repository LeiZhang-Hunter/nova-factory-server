package dao

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IDiscountRuleDao 折扣规则数据访问接口（用于价格计算）
type IDiscountRuleDao interface {
	// GetValidRule 获取用户对指定目标的的有效折扣规则
	GetValidRule(c *gin.Context, userID int64, targetType string, targetID int64) (*models.UserDiscountRule, error)
	// ListByUserAndGoods 查询用户对商品的折扣规则（优先）和分类折扣规则
	ListByUserAndGoods(c *gin.Context, userID int64, goodsID int64, categoryID int64) (*models.UserDiscountRule, error)
	// ListValidRulesByTargets 批量查询用户的有效折扣规则（用于列表场景消除 N+1）
	ListValidRulesByTargets(c *gin.Context, userID int64, targetType string, targetIDs []int64) ([]*models.UserDiscountRule, error)
}
