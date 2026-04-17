package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopAddressDaoImpl 提供商城用户地址表的数据访问能力。
type ShopAddressDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopAddressDao 创建商城用户地址 DAO。
func NewShopAddressDao(ms *gorm.DB) dao.IShopAddressDao {
	return &ShopAddressDaoImpl{
		db:        ms,
		tableName: "shop_user_address",
	}
}

// Set 新增或修改商城用户地址，并在设置默认地址时清理同用户其他默认地址。
func (s *ShopAddressDaoImpl) Set(c *gin.Context, req *models.AddressSetReq) (*models.Address, error) {
	deptID := baizeContext.GetDeptId(c)
	userID := baizeContext.GetUserId(c)
	status := int32(1)
	if req.Status != nil {
		status = *req.Status
	}

	var result *models.Address
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if req.IsDefault == 1 {
			clearDefault := tx.Table(s.tableName).
				Where("user_id = ?", req.UserID).
				Where("dept_id = ?", deptID).
				Where("state = ?", commonStatus.NORMAL)
			if req.ID > 0 {
				clearDefault = clearDefault.Where("id <> ?", req.ID)
			}
			if err := clearDefault.Updates(map[string]interface{}{
				"is_default":  0,
				"update_by":   userID,
				"update_time": gorm.Expr("NOW()"),
			}).Error; err != nil {
				return err
			}
		}

		if req.ID > 0 {
			existing, err := s.getByIDTx(c, tx, req.ID)
			if err != nil {
				return err
			}
			if existing == nil {
				return errors.New("地址不存在")
			}

			if err := tx.Table(s.tableName).
				Where("id = ?", req.ID).
				Where("dept_id = ?", deptID).
				Where("state = ?", commonStatus.NORMAL).
				Updates(map[string]interface{}{
					"user_id":         req.UserID,
					"receiver_name":   req.ReceiverName,
					"receiver_mobile": req.ReceiverMobile,
					"province_code":   req.ProvinceCode,
					"province_name":   req.ProvinceName,
					"city_code":       req.CityCode,
					"city_name":       req.CityName,
					"district_code":   req.DistrictCode,
					"district_name":   req.DistrictName,
					"street_code":     req.StreetCode,
					"street_name":     req.StreetName,
					"detail_address":  req.DetailAddress,
					"postal_code":     req.PostalCode,
					"address_label":   req.AddressLabel,
					"is_default":      req.IsDefault,
					"status":          status,
					"update_by":       userID,
					"update_time":     gorm.Expr("NOW()"),
				}).Error; err != nil {
				return err
			}
			result, err = s.getByIDTx(c, tx, req.ID)
			return err
		}

		model := &models.Address{
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
			StreetCode:     req.StreetCode,
			StreetName:     req.StreetName,
			DetailAddress:  req.DetailAddress,
			PostalCode:     req.PostalCode,
			AddressLabel:   req.AddressLabel,
			IsDefault:      req.IsDefault,
			Status:         status,
			DeptID:         deptID,
			State:          commonStatus.NORMAL,
		}
		model.SetCreateBy(userID)
		if err := tx.Table(s.tableName).Create(model).Error; err != nil {
			return err
		}
		result = model
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List 查询商城用户地址列表。
func (s *ShopAddressDaoImpl) List(c *gin.Context, req *models.AddressQuery) (*models.AddressListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.UserID != "" {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.ReceiverName != "" {
		db = db.Where("receiver_name LIKE ?", "%"+req.ReceiverName+"%")
	}
	if req.ReceiverMobile != "" {
		db = db.Where("receiver_mobile = ?", req.ReceiverMobile)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
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
	rows := make([]*models.Address, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Order("is_default DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.AddressListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// Remove 软删除商城用户地址。
func (s *ShopAddressDaoImpl) Remove(c *gin.Context, ids []int64) error {
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

func (s *ShopAddressDaoImpl) getByIDTx(c *gin.Context, tx *gorm.DB, id int64) (*models.Address, error) {
	var item models.Address
	if err := tx.WithContext(c).Table(s.tableName).
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
