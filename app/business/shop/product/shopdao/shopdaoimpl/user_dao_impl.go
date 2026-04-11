package shopdaoimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopUserDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopUserDao(ms *gorm.DB) shopdao.IShopUserDao {
	return &ShopUserDaoImpl{
		db:        ms,
		tableName: "shop_user",
	}
}

func (s *ShopUserDaoImpl) Create(c *gin.Context, req *shopmodels.UserUpsert) (*shopmodels.User, error) {
	model := &shopmodels.User{
		UserID:       req.UserID,
		Username:     req.Username,
		Nickname:     req.Nickname,
		Mobile:       req.Mobile,
		Email:        req.Email,
		Password:     req.Password,
		UserType:     req.UserType,
		CompanyName:  req.CompanyName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
	}
	if model.Status != 0 && model.Status != 1 {
		model.Status = 1
	}
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopUserDaoImpl) Update(c *gin.Context, req *shopmodels.UserUpsert) (*shopmodels.User, error) {
	updates := map[string]interface{}{
		"user_id":       req.UserID,
		"username":      req.Username,
		"nickname":      req.Nickname,
		"mobile":        req.Mobile,
		"email":         req.Email,
		"password":      req.Password,
		"user_type":     req.UserType,
		"company_name":  req.CompanyName,
		"contact_name":  req.ContactName,
		"contact_phone": req.ContactPhone,
		"status":        req.Status,
	}
	if req.Status != 0 && req.Status != 1 {
		delete(updates, "status")
	}
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

func (s *ShopUserDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Update("is_deleted", 1).Error
}

func (s *ShopUserDaoImpl) GetByID(c *gin.Context, id int64) (*shopmodels.User, error) {
	var item shopmodels.User
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).Where("is_deleted = 0").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ShopUserDaoImpl) List(c *gin.Context, req *shopmodels.UserQuery) (*shopmodels.UserListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).Where("is_deleted = 0")
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Mobile != "" {
		db = db.Where("mobile = ?", req.Mobile)
	}
	if req.UserType > 0 {
		db = db.Where("user_type = ?", req.UserType)
	}
	if req.Status == 0 || req.Status == 1 {
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
	rows := make([]*shopmodels.User, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopmodels.UserListData{
		Rows:  rows,
		Total: total,
	}, nil
}
