package impl

import (
	"errors"
	"fmt"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopWechatUserDaoImpl 提供微信用户数据访问能力。
type IApiShopWechatUserDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopWechatUserDaoImpl 创建微信用户 DAO。
func NewIApiShopWechatUserDaoImpl(ms *gorm.DB) dao.IApiShopWechatUserDao {
	return &IApiShopWechatUserDaoImpl{
		db:        ms,
		tableName: "shop_user",
	}
}

// GetByOpenid 根据微信openid查询商城用户。
func (s *IApiShopWechatUserDaoImpl) GetByOpenid(c *gin.Context, openid string) (*models.User, error) {
	var item models.User
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("wechat_openid = ?", openid).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// CreateWechatUser 创建微信用户。
func (s *IApiShopWechatUserDaoImpl) CreateWechatUser(c *gin.Context, req *models.WechatUserCreate) (*models.User, error) {
	// 微信登录无 session，使用默认 dept_id=0
	const defaultDeptID int64 = 0
	model := &models.User{
		ID:           snowflake.GenID(),
		UserID:       fmt.Sprintf("%d", snowflake.GenID()),
		Username:     req.Username,
		Nickname:     req.Nickname,
		Avatar:       req.Avatar,
		UserType:     req.UserType,
		Status:       req.Status,
		WechatOpenid: req.Openid,
		DeptID:       defaultDeptID,
		State:        commonStatus.NORMAL,
	}
	if model.Status == nil {
		model.Status = boolPtr(false)
	}
	model.SetCreateBy(defaultDeptID)
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

// GetByID 根据用户ID查询商城用户（不带 dept_id 过滤，用于小程序）。
func (s *IApiShopWechatUserDaoImpl) GetByID(c *gin.Context, id int64) (*models.User, error) {
	var item models.User
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

func boolPtr(v bool) *bool {
	return &v
}
