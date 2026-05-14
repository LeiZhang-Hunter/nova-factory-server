package impl

import (
	"errors"

	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	utiltime "nova-factory-server/app/utils/time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopUserDiscountRuleDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopUserDiscountRuleDao(ms *gorm.DB) dao.IShopUserDiscountRuleDao {
	return &ShopUserDiscountRuleDaoImpl{
		db:        ms,
		tableName: "shop_user_discount_rules",
	}
}

func (s *ShopUserDiscountRuleDaoImpl) Create(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error) {
	validFrom, _ := utiltime.ParseDateTime(req.ValidFrom)
	validTo, _ := utiltime.ParseDateTime(req.ValidTo)
	model := &models.UserDiscountRule{
		ID:           snowflake.GenID(),
		UserID:       req.UserID,
		TargetType:   req.TargetType,
		TargetID:     req.TargetID,
		DiscountRate: float64(req.DiscountRate) / 100,
		ValidFrom:    validFrom,
		ValidTo:      validTo,
		State:        commonStatus.NORMAL,
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopUserDiscountRuleDaoImpl) BatchCreate(c *gin.Context, rules []*models.UserDiscountRule) error {
	if len(rules) == 0 {
		return nil
	}
	return s.db.WithContext(c).Table(s.tableName).CreateInBatches(rules, 100).Error
}

func (s *ShopUserDiscountRuleDaoImpl) GetByID(c *gin.Context, id int64) (*models.UserDiscountRule, error) {
	var item models.UserDiscountRule
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopUserDiscountRuleDaoImpl) GetByUserAndTarget(c *gin.Context, userID int64, targetType string, targetID string) (*models.UserDiscountRule, error) {
	var item models.UserDiscountRule
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopUserDiscountRuleDaoImpl) ListByUserID(c *gin.Context, userID int64, page, size int64) (*models.UserDiscountRuleListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("state = ?", commonStatus.NORMAL)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.UserDiscountRule, 0)
	if err := db.Offset(int((page - 1) * size)).Limit(int(size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.UserDiscountRuleListData{Rows: rows, Total: total}, nil
}

func (s *ShopUserDiscountRuleDaoImpl) ListByUserIDAndType(c *gin.Context, userID int64, targetType string, page, size int64) (*models.UserDiscountRuleListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("target_type = ?", targetType).
		Where("state = ?", commonStatus.NORMAL)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.UserDiscountRule, 0)
	if err := db.Offset(int((page - 1) * size)).Limit(int(size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.UserDiscountRuleListData{Rows: rows, Total: total}, nil
}

func (s *ShopUserDiscountRuleDaoImpl) ListByTarget(c *gin.Context, targetType string, targetID string, page, size int64) (*models.UserDiscountRuleListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("state = ?", commonStatus.NORMAL)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.UserDiscountRule, 0)
	if err := db.Offset(int((page - 1) * size)).Limit(int(size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.UserDiscountRuleListData{Rows: rows, Total: total}, nil
}

func (s *ShopUserDiscountRuleDaoImpl) Update(c *gin.Context, req *models.UserDiscountRuleUpsert) (*models.UserDiscountRule, error) {
	validFrom, _ := utiltime.ParseDateTime(req.ValidFrom)
	validTo, _ := utiltime.ParseDateTime(req.ValidTo)
	updates := &models.UserDiscountRule{
		ID:           req.ID,
		DiscountRate: float64(req.DiscountRate) / 100,
		ValidFrom:    validFrom,
		ValidTo:      validTo,
	}
	updates.SetUpdateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Select("discount_rate", "valid_from", "valid_to", "update_by", "update_time").
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, req.ID)
}

func (s *ShopUserDiscountRuleDaoImpl) DeleteByID(c *gin.Context, id int64) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": gorm.Expr("NOW()"),
		}).Error
}

func (s *ShopUserDiscountRuleDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Delete(&models.UserDiscountRule{}).Error
}

func (s *ShopUserDiscountRuleDaoImpl) GetUsedTargetIdsByUserID(c *gin.Context, userID int64) (goodsIds []string, categoryIds []string, err error) {
	type result struct {
		TargetType string
		TargetID   string
	}
	var results []result
	err = s.db.WithContext(c).Table(s.tableName).
		Select("target_type, target_id").
		Where("user_id = ?", userID).
		Where("state = ?", commonStatus.NORMAL).
		Find(&results).Error
	if err != nil {
		return nil, nil, err
	}
	for _, r := range results {
		if r.TargetType == "goods" {
			goodsIds = append(goodsIds, r.TargetID)
		} else if r.TargetType == "category" {
			categoryIds = append(categoryIds, r.TargetID)
		}
	}
	return goodsIds, categoryIds, nil
}
