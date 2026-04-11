package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"strconv"

	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopUserDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopUserDao(ms *gorm.DB) dao.IShopUserDao {
	return &ShopUserDaoImpl{
		db:        ms,
		tableName: "shop_user",
	}
}

func (s *ShopUserDaoImpl) Create(c *gin.Context, req *models.UserUpsert) (*models.User, error) {
	model := &models.User{
		ID:           snowflake.GenID(),
		UserID:       req.UserID,
		Username:     req.Username,
		Nickname:     req.Nickname,
		Mobile:       req.Mobile,
		Email:        req.Email,
		Password:     req.Password,
		UserType:     req.UserType,
		Avatar:       req.Avatar,
		CompanyName:  req.CompanyName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		DeptID:       baizeContext.GetDeptId(c),
		State:        commonStatus.NORMAL,
	}
	if model.UserID == "" {
		model.UserID = strconv.FormatInt(snowflake.GenID(), 10)
	}
	if model.Status == nil {
		model.Status = boolPtr(false)
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	model.DeptID = baizeContext.GetDeptId(c)
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopUserDaoImpl) Update(c *gin.Context, req *models.UserUpsert) (*models.User, error) {
	updates := &models.User{
		ID:           req.ID,
		UserID:       req.UserID,
		Username:     req.Username,
		Nickname:     req.Nickname,
		Mobile:       req.Mobile,
		Email:        req.Email,
		Password:     req.Password,
		UserType:     req.UserType,
		Avatar:       req.Avatar,
		CompanyName:  req.CompanyName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		DeptID:       baizeContext.GetDeptId(c),
	}
	if updates.Status == nil {
		updates.Status = boolPtr(false)
	}
	updates.SetUpdateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select("user_id", "username", "nickname", "mobile", "email", "password", "user_type", "avatar", "company_name", "contact_name", "contact_phone", "status", "update_by", "update_time").
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, req.ID)
}

func (s *ShopUserDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": gorm.Expr("NOW()"),
		}).Error
}

func (s *ShopUserDaoImpl) GetByID(c *gin.Context, id int64) (*models.User, error) {
	var item models.User
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopUserDaoImpl) List(c *gin.Context, req *models.UserQuery) (*models.UserListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.Username != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ?", "%"+req.Username+"%", "%"+req.Username+"%")
	}
	if req.Mobile != "" {
		db = db.Where("mobile = ?", req.Mobile)
	}
	if req.UserType > 0 {
		db = db.Where("user_type = ?", req.UserType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.User, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.UserListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func boolPtr(v bool) *bool {
	return &v
}
