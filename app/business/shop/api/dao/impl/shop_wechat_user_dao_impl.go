package impl

import (
	"errors"
	"fmt"

	"nova-factory-server/app/business/shop/api/dao"
	wechatModels "nova-factory-server/app/business/shop/api/models"
	userModels "nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopWechatUserDaoImpl 提供微信用户数据访问能力。
type ShopWechatUserDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopWechatUserDao 创建微信用户 DAO。
func NewShopWechatUserDao(ms *gorm.DB) dao.IShopWechatUserDao {
	return &ShopWechatUserDaoImpl{
		db:        ms,
		tableName: "shop_user",
	}
}

// GetByOpenid 根据微信openid查询商城用户。
func (s *ShopWechatUserDaoImpl) GetByOpenid(c *gin.Context, openid string) (*userModels.User, error) {
	var item userModels.User
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
func (s *ShopWechatUserDaoImpl) CreateWechatUser(c *gin.Context, req *wechatModels.WechatUserCreate) (*userModels.User, error) {
	// 微信登录无 session，使用默认 dept_id=0
	const defaultDeptID int64 = 0
	model := &userModels.User{
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

func boolPtr(v bool) *bool {
	return &v
}
