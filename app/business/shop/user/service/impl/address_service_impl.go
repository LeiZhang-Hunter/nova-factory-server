package impl

import (
	"errors"
	"github.com/itmisx/go_regions"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"strconv"
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
	req.CityCode = strings.TrimSpace(req.CityCode)
	req.DistrictCode = strings.TrimSpace(req.DistrictCode)
	req.StreetCode = strings.TrimSpace(req.StreetCode)
	req.ProvinceName = ""
	req.CityName = ""
	req.DistrictName = ""
	req.StreetName = ""
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
	req.ReceiverName = resolveReceiverName(user)
	if err = fillRegionNames(req); err != nil {
		return nil, err
	}

	if req.ProvinceName == "" {
		return nil, errors.New("省份输入错误")
	}

	if req.CityName == "" {
		return nil, errors.New("市输入错误")
	}

	if req.DistrictName == "" {
		return nil, errors.New("输入区错误")
	}

	if req.StreetName == "" {
		return nil, errors.New("街道不能为空")
	}

	req.UserID = user.ID
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

func fillRegionNames(req *models.AddressSetReq) error {
	var err error
	req.ProvinceName, err = regionNameByCode(req.ProvinceCode, "省")
	if err != nil {
		return err
	}
	req.CityName, err = regionNameByCode(req.CityCode, "市")
	if err != nil {
		return err
	}
	req.DistrictName, err = regionNameByCode(req.DistrictCode, "区")
	if err != nil {
		return err
	}
	req.StreetName, err = regionNameByCode(req.StreetCode, "街道")
	if err != nil {
		return err
	}
	return nil
}

func regionNameByCode(code string, label string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", nil
	}
	id, err := strconv.Atoi(code)
	if err != nil {
		return "", errors.New(label + "编码格式不正确")
	}
	info := go_regions.RegionInfo(id)
	if info == nil {
		return "", errors.New(label + "编码不存在")
	}
	return strings.TrimSpace(info.Name), nil
}

func (s *ShopAddressServiceImpl) Query(c *gin.Context, req *models.UserAddressInfoQuery) (*models.AddressListData, error) {
	user, err := s.userDao.GetByUsername(c, req.Username)
	if err != nil {
		zap.L().Error("get username error", zap.Error(err))
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	list, err := s.dao.List(c, &models.AddressQuery{
		UserId: user.ID,
		Page:   1,
		Size:   20,
	})
	if err != nil {
		zap.L().Error("list error", zap.Error(err))
		return list, err
	}

	return list, nil
}
