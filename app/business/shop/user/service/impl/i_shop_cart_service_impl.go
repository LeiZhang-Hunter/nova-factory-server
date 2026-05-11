package impl

import (
	"errors"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/utils/fileUtils"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShopCartServiceImpl 提供商城用户购物车相关业务能力。
type ShopCartServiceImpl struct {
	dao     dao.IShopCartDao
	userDao dao.IShopUserDao
}

// NewShopCartService 创建商城用户购物车服务。
func NewShopCartService(dao dao.IShopCartDao, userDao dao.IShopUserDao) service.IShopCartService {
	return &ShopCartServiceImpl{
		dao:     dao,
		userDao: userDao,
	}
}

// Set 新增或修改商城用户购物车项。
func (s *ShopCartServiceImpl) Set(c *gin.Context, req *models.CartSetReq) (*models.Cart, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	req.GoodsID = strings.TrimSpace(req.GoodsID)
	req.SkuID = strings.TrimSpace(req.SkuID)
	req.GoodsName = strings.TrimSpace(req.GoodsName)
	req.SkuName = strings.TrimSpace(req.SkuName)
	req.ImageURL = strings.TrimSpace(req.ImageURL)

	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.GoodsID == "" {
		return nil, errors.New("商品ID不能为空")
	}
	if req.SkuID == "" {
		return nil, errors.New("SKU ID不能为空")
	}
	if req.GoodsName == "" {
		return nil, errors.New("商品名称不能为空")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("购买数量必须大于0")
	}

	info, err := s.userDao.GetByUsername(c, req.Username)
	if err != nil {
		zap.L().Error("get user info fail", zap.String("username", req.Username), zap.Error(err))
		return nil, errors.New("读取用户错误")
	}

	if info == nil {
		return nil, errors.New("读取用户不存在")
	}

	user, err := s.userDao.GetByUserID(c, info.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("商城用户不存在")
	}

	req.UserID = user.ID
	return s.dao.Set(c, req)
}

// GetByID 根据主键查询商城用户购物车项。
func (s *ShopCartServiceImpl) GetByID(c *gin.Context, id int64) (*models.Cart, error) {
	if id == 0 {
		return nil, errors.New("购物车ID不能为空")
	}
	return s.dao.GetByID(c, id)
}

// List 查询商城用户购物车列表。
func (s *ShopCartServiceImpl) List(c *gin.Context, req *models.CartQuery) (*models.CartListData, error) {
	if req == nil {
		req = new(models.CartQuery)
	}
	req.GoodsID = strings.TrimSpace(req.GoodsID)
	req.SkuID = strings.TrimSpace(req.SkuID)
	ret, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	if ret == nil || len(ret.Rows) == 0 {
		return ret, nil
	}

	for k, v := range ret.Rows {
		ret.Rows[k].ImageURL = fileUtils.BuildAbsoluteURL(c, v.ImageURL)
	}
	return ret, nil
}

// Remove 删除商城用户购物车项。
func (s *ShopCartServiceImpl) Remove(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的购物车记录")
	}
	return s.dao.Remove(c, ids)
}
