package impl

import (
	"errors"
	"fmt"
	shopusermodels "nova-factory-server/app/business/shop/user/models"

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
func (s *IApiShopWechatUserDaoImpl) GetByOpenid(c *gin.Context, openid string) (*shopusermodels.User, error) {
	var item shopusermodels.User
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

// GetByAccount 根据用户名或手机号查询商城用户，不带 dept_id 过滤。
func (s *IApiShopWechatUserDaoImpl) GetByAccount(c *gin.Context, account string) (*shopusermodels.User, error) {
	var item shopusermodels.User
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("(username = ? OR mobile = ?)", account, account).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// BindWechatOpenid binds a WeChat openid to a shop_user without dept filtering.
func (s *IApiShopWechatUserDaoImpl) BindWechatOpenid(c *gin.Context, id int64, openid string) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"wechat_openid": openid,
			"update_time":   gorm.Expr("NOW()"),
		}).Error
}

// CreateWechatUser 创建微信用户。
func (s *IApiShopWechatUserDaoImpl) CreateWechatUser(c *gin.Context, req *models.WechatUserCreate) (*shopusermodels.User, error) {
	// 微信登录无 session，使用默认 dept_id=0
	const defaultDeptID int64 = 0
	model := &shopusermodels.User{
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
func (s *IApiShopWechatUserDaoImpl) GetByID(c *gin.Context, id int64) (*shopusermodels.User, error) {
	var item shopusermodels.User
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
func (s *IApiShopWechatUserDaoImpl) GetByUserID(c *gin.Context, userId int64) (*shopusermodels.User, error) {
	var item shopusermodels.User
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userId).
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

// UpdatePassword 更新商城用户密码（bcrypt 哈希后写入）。
func (s *IApiShopWechatUserDaoImpl) UpdatePassword(c *gin.Context, userID int64, newPassword string) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", userID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"password":    newPassword,
			"update_time": gorm.Expr("NOW()"),
		}).Error
}

// UpdateProfile 更新商城用户个人资料（昵称、手机号、邮箱）。
func (s *IApiShopWechatUserDaoImpl) UpdateProfile(c *gin.Context, userID int64, updates map[string]interface{}) error {
	updates["update_time"] = gorm.Expr("NOW()")
	return s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", userID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(updates).Error
}
