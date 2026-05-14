package dao

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopUserDiscountRuleDao 商城用户折扣规则数据访问接口
type IShopUserDiscountRuleDao interface {
	// Create 创建折扣规则
	Create(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error)
	// BatchCreate 批量创建折扣规则
	BatchCreate(c *gin.Context, rules []*models.UserDiscountRule) error
	// GetByID 根据ID查询
	GetByID(c *gin.Context, id int64) (*models.UserDiscountRule, error)
	// GetByUserAndTarget 根据用户ID和目标查询唯一规则
	GetByUserAndTarget(c *gin.Context, userID int64, targetType string, targetID string) (*models.UserDiscountRule, error)
	// ListByUserID 根据用户ID查询所有规则
	ListByUserID(c *gin.Context, userID int64, page, size int64) (*models.UserDiscountRuleListData, error)
	// ListByUserIDAndType 根据用户ID和类型查询规则
	ListByUserIDAndType(c *gin.Context, userID int64, targetType string, page, size int64) (*models.UserDiscountRuleListData, error)
	// ListByTarget 根据目标查询所有用户的规则
	ListByTarget(c *gin.Context, targetType string, targetID string, page, size int64) (*models.UserDiscountRuleListData, error)
	// Update 更新折扣规则
	Update(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error)
	// DeleteByID 软删除单条规则
	DeleteByID(c *gin.Context, id int64) error
	// DeleteByIDs 批量软删除规则
	DeleteByIDs(c *gin.Context, ids []int64) error
	// GetUsedTargetIdsByUserID 获取用户所有已使用的目标ID（去重）
	GetUsedTargetIdsByUserID(c *gin.Context, userID int64) (goodsIds []string, categoryIds []string, err error)
}
