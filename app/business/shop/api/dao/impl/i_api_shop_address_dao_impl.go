package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopAddressDaoImpl 移动端地址 DAO 实现
type IApiShopAddressDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopAddressDao 创建移动端地址 DAO
func NewShopAddressDao(ms *gorm.DB) dao.IApiShopAddressDao {
	return &IApiShopAddressDaoImpl{
		db:        ms,
		tableName: "shop_user_address",
	}
}

// Set 新增或修改地址
func (s *IApiShopAddressDaoImpl) Set(c *gin.Context, req *models.AddressSetReq) (*models.ShopUserAddressApp, error) {
	var result *models.ShopUserAddressApp
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 如果设置默认，先清除其他默认
		if req.IsDefault != nil && *req.IsDefault {
			if err := s.ClearDefault(c, req.UserID, 0); err != nil {
				return err
			}
		}

		if req.ID > 0 {
			// 更新
			if err := tx.Table(s.tableName).
				Where("id = ?", req.ID).
				Updates(map[string]interface{}{
					"receiver_name":   req.ReceiverName,
					"receiver_mobile": req.ReceiverMobile,
					"province_code":   req.ProvinceCode,
					"province_name":   req.ProvinceName,
					"city_code":       req.CityCode,
					"city_name":       req.CityName,
					"district_code":   req.DistrictCode,
					"district_name":   req.DistrictName,
					"detail_address":  req.DetailAddress,
					"address_label":   req.AddressLabel,
					"is_default":      req.IsDefault,
					"create_time":     gorm.Expr("NOW()"),
				}).Error; err != nil {
				return err
			}
		} else {
			// 新增
			model := &models.ShopUserAddressApp{
				ID:             snowflake.GenID(),
				UserID:         req.UserID,
				ReceiverName:   req.ReceiverName,
				ReceiverMobile: req.ReceiverMobile,
				ProvinceCode:   req.ProvinceCode,
				ProvinceName:   req.ProvinceName,
				CityCode:       req.CityCode,
				CityName:       req.CityName,
				DistrictCode:   req.DistrictCode,
				DistrictName:   req.DistrictName,
				DetailAddress:  req.DetailAddress,
				AddressLabel:   req.AddressLabel,
				IsDefault:      0,
			}
			if req.IsDefault != nil && *req.IsDefault {
				model.IsDefault = 1
			}
			if err := tx.Table(s.tableName).Create(model).Error; err != nil {
				return err
			}
			result = model
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return s.GetByID(c, req.ID)
	}
	return result, nil
}

// GetByID 根据 ID 查询
func (s *IApiShopAddressDaoImpl) GetByID(c *gin.Context, id int64) (*models.ShopUserAddressApp, error) {
	var item models.ShopUserAddressApp
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// List 查询用户地址列表
func (s *IApiShopAddressDaoImpl) List(c *gin.Context, userId int64) (*models.AddressListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).Where("user_id = ?", userId)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.ShopUserAddressApp, 0)
	if err := db.Order("is_default DESC, id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.AddressListData{Rows: rows, Total: total}, nil
}

// Remove 删除地址
func (s *IApiShopAddressDaoImpl) Remove(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Delete(nil).Error
}

// ClearDefault 清除用户的所有默认地址
func (s *IApiShopAddressDaoImpl) ClearDefault(c *gin.Context, userId int64, excludeId int64) error {
	db := s.db.WithContext(c).Table(s.tableName).Where("user_id = ?", userId)
	if excludeId > 0 {
		db = db.Where("id <> ?", excludeId)
	}
	return db.Update("is_default", 0).Error
}
