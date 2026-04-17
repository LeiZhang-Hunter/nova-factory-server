package impl

import (
	"errors"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShopAddressServiceImpl 提供商城用户地址相关业务能力。
type ShopAddressServiceImpl struct {
	dao     dao.IShopAddressDao
	userDao dao.IShopUserDao
}

// NewShopAddressService 创建商城用户地址服务。
func NewShopAddressService(dao dao.IShopAddressDao, userDao dao.IShopUserDao) service.IShopAddressService {
	return &ShopAddressServiceImpl{
		dao:     dao,
		userDao: userDao,
	}
}

// Set 新增或修改商城用户地址。
func (s *ShopAddressServiceImpl) Set(c *gin.Context, req *models.AddressSetReq) (*models.Address, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	req.ReceiverMobile = strings.TrimSpace(req.ReceiverMobile)
	req.ProvinceCode = strings.TrimSpace(req.ProvinceCode)
	req.ProvinceName = strings.TrimSpace(req.ProvinceName)
	req.CityCode = strings.TrimSpace(req.CityCode)
	req.CityName = strings.TrimSpace(req.CityName)
	req.DistrictCode = strings.TrimSpace(req.DistrictCode)
	req.DistrictName = strings.TrimSpace(req.DistrictName)
	req.StreetCode = strings.TrimSpace(req.StreetCode)
	req.StreetName = strings.TrimSpace(req.StreetName)
	req.DetailAddress = strings.TrimSpace(req.DetailAddress)
	req.PostalCode = strings.TrimSpace(req.PostalCode)
	req.AddressLabel = strings.TrimSpace(req.AddressLabel)
	if req.ReceiverMobile == "" {
		return nil, errors.New("收货人手机号不能为空")
	}

	user, err := s.userDao.GetByMobile(c, req.ReceiverMobile)
	if err != nil {
		zap.L().Error("get user info by mobile error", zap.String("mobile", req.ReceiverMobile), zap.Error(err))
		return nil, err
	}
	if user == nil {
		zap.L().Error("not found user info", zap.String("mobile", req.ReceiverMobile), zap.Error(err))
		return nil, errors.New("未找到对应商城用户")
	}
	req.UserID = user.UserID
	req.ReceiverName = resolveReceiverName(user)
	return s.dao.Set(c, req)
}

// GetByID 根据主键查询商城用户地址。
func (s *ShopAddressServiceImpl) GetByID(c *gin.Context, id int64) (*models.Address, error) {
	return s.dao.GetByID(c, id)
}

// List 查询商城用户地址列表。
func (s *ShopAddressServiceImpl) List(c *gin.Context, req *models.AddressQuery) (*models.AddressListData, error) {
	return s.dao.List(c, req)
}

// Remove 删除商城用户地址。
func (s *ShopAddressServiceImpl) Remove(c *gin.Context, ids []int64) error {
	return s.dao.Remove(c, ids)
}

func resolveReceiverName(user *models.User) string {
	if user == nil {
		return ""
	}
	if value := strings.TrimSpace(user.ContactName); value != "" {
		return value
	}
	if value := strings.TrimSpace(user.Nickname); value != "" {
		return value
	}
	return strings.TrimSpace(user.Username)
}
