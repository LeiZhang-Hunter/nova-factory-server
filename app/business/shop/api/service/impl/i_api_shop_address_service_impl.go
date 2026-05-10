package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// IApiShopAddressServiceImpl 移动端地址服务实现
type IApiShopAddressServiceImpl struct {
	dao dao.IApiShopAddressDao
}

// NewIApiShopAddressServiceImpl  创建移动端地址服务
func NewIApiShopAddressServiceImpl(dao dao.IApiShopAddressDao) service.IApiShopAddressService {
	return &IApiShopAddressServiceImpl{dao: dao}
}

// Set 新增或修改地址
func (s *IApiShopAddressServiceImpl) Set(c *gin.Context, req *models.AddressSetReq) (*models.ShopUserAddressApp, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}

	// 校验
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	req.ReceiverMobile = strings.TrimSpace(req.ReceiverMobile)
	req.ProvinceCode = strings.TrimSpace(req.ProvinceCode)
	req.ProvinceName = strings.TrimSpace(req.ProvinceName)
	req.CityCode = strings.TrimSpace(req.CityCode)
	req.CityName = strings.TrimSpace(req.CityName)
	req.DistrictCode = strings.TrimSpace(req.DistrictCode)
	req.DistrictName = strings.TrimSpace(req.DistrictName)
	req.DetailAddress = strings.TrimSpace(req.DetailAddress)
	req.AddressLabel = strings.TrimSpace(req.AddressLabel)

	if req.ReceiverName == "" {
		return nil, errors.New("收货人姓名不能为空")
	}
	if req.ReceiverMobile == "" {
		return nil, errors.New("手机号不能为空")
	}
	if !isValidPhone(req.ReceiverMobile) {
		return nil, errors.New("手机号格式不正确")
	}
	if req.ProvinceCode == "" || req.ProvinceName == "" {
		return nil, errors.New("省份不能为空")
	}
	if req.CityCode == "" || req.CityName == "" {
		return nil, errors.New("城市不能为空")
	}
	if req.DetailAddress == "" {
		return nil, errors.New("详细地址不能为空")
	}

	return s.dao.Set(c, req)
}

// GetByID 根据 ID 查询
func (s *IApiShopAddressServiceImpl) GetByID(c *gin.Context, id int64) (*models.ShopUserAddressApp, error) {
	if id == 0 {
		return nil, errors.New("地址ID不能为空")
	}
	return s.dao.GetByID(c, id)
}

// List 查询用户地址列表
func (s *IApiShopAddressServiceImpl) List(c *gin.Context, userId int64) (*models.AddressListData, error) {
	if userId == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	return s.dao.List(c, userId)
}

// Remove 删除地址
func (s *IApiShopAddressServiceImpl) Remove(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("地址ID不能为空")
	}
	return s.dao.Remove(c, ids)
}

func isValidPhone(phone string) bool {
	return len(phone) == 11 && phone[0] == '1'
}

// Default 查询用户默认地址
func (s *IApiShopAddressServiceImpl) Default(c *gin.Context, userId int64) (*models.ShopUserAddressApp, error) {
	if userId == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	return s.dao.Default(c, userId)
}
