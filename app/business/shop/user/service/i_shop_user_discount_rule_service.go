package service

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopUserDiscountRuleService 商城用户折扣规则服务接口
type IShopUserDiscountRuleService interface {
	// Create 创建折扣规则
	Create(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error)
	// BatchCreate 批量创建折扣规则
	BatchCreate(c *gin.Context, req *models.BatchDiscountRuleCreate) (int64, error)
	// GetByID 根据ID查询
	GetByID(c *gin.Context, id int64) (*models.UserDiscountRule, error)
	// ListByUserID 根据用户ID查询所有规则
	ListByUserID(c *gin.Context, userID int64, page, size int64) (*models.UserDiscountRuleListData, error)
	// ListByUserIDAndType 根据用户ID和类型查询规则
	ListByUserIDAndType(c *gin.Context, userID int64, targetType string, page, size int64) (*models.UserDiscountRuleListData, error)
	// ListByUserIDs 批量查询用户的折扣规则
	ListByUserIDs(c *gin.Context, userIDs []int64, targetType string, targetIDs []string) ([]*models.UserDiscountRule, error)
	// Update 更新折扣规则
	Update(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error)
	// DeleteByID 软删除单条规则
	DeleteByID(c *gin.Context, id int64) error
	// DeleteByIDs 批量软删除规则
	DeleteByIDs(c *gin.Context, ids []int64) error
	// GetValidRule 获取有效的折扣规则（用于价格计算）
	GetValidRule(c *gin.Context, userID int64, targetType string, targetID string) (*models.UserDiscountRule, error)
	// GetUsedTargetIds 获取用户所有已使用的目标ID（去重）
	GetUsedTargetIds(c *gin.Context, userID int64) (goodsIds []string, categoryIds []string, err error)
}
