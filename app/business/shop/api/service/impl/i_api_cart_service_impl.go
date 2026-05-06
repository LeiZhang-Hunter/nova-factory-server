package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/api/dao"
	api "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/business/shop/product/shopdao"
	userDao "nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/utils/baizeContext"
)

type ShopCartServiceImpl struct {
	goodsDao     shopdao.IShopGoodsDao
	goodsSkuDao  shopdao.IShopSkuDao
	shopUserDao  userDao.IShopUserDao
	iShopCartDao dao.IApiShopCartDao
}

// NewShopCartServiceImpl 购物车服务
func NewShopCartServiceImpl(goodsDao shopdao.IShopGoodsDao,
	goodsSkuDao shopdao.IShopSkuDao,
	shopUserDao userDao.IShopUserDao,
	iShopCartDao dao.IApiShopCartDao) service.IApiShopCartService {
	return &ShopCartServiceImpl{
		goodsDao:     goodsDao,
		goodsSkuDao:  goodsSkuDao,
		iShopCartDao: iShopCartDao,
		shopUserDao:  shopUserDao,
	}
}

// GenCart 生成购物车，购物车下单
func (s *ShopCartServiceImpl) GenCart(c *gin.Context, req *api.CartSetDataReq) (*api.CartDto, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}

	if req.GoodsID == 0 {
		return nil, errors.New("商品ID不能为空")
	}
	if req.SkuID == 0 {
		return nil, errors.New("SKU ID不能为空")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("购买数量必须大于0")
	}

	info, err := s.goodsDao.GetByID(c, req.GoodsID)
	if err != nil {
		zap.L().Error("get goods info failed", zap.Error(err))
		return nil, errors.New("读取商品信息失败")
	}

	if info == nil {
		return nil, errors.New("商品信息不存在")
	}

	if info.GoodsID == "" {
		return nil, errors.New("商品id不能为空")
	}

	if info.GoodsName == "" {
		return nil, errors.New("读取商品名称信息失败")
	}

	skuInfo, err := s.goodsSkuDao.GetByID(c, req.SkuID)
	if err != nil {
		zap.L().Error("get goods sku info failed", zap.Error(err))
		return nil, errors.New("读取商品sku信息失败")
	}

	if skuInfo == nil {
		return nil, errors.New("商品sku信息不存在")
	}

	if skuInfo.GoodsID != info.GoodsID {
		return nil, errors.New("sku不属于这个商品")
	}

	if skuInfo.SkuName == "" {
		return nil, errors.New("读取商品sku名称信息失败")
	}

	userId := baizeContext.GetUserId(c)

	userInfo, err := s.shopUserDao.GetByID(c, userId)
	if err != nil {
		zap.L().Error("get user info failed", zap.Error(err))
		return nil, errors.New("读取用户信息失败")
	}

	if userInfo == nil {
		zap.L().Error("get user info failed", zap.Error(err))
		return nil, errors.New("用户不存在")
	}

	// 查询用户是否存在这个sku的购物车
	cartInfo, err := s.iShopCartDao.Save(c, &api.CartSetData{
		UserID:      userId,
		Username:    userInfo.Username,
		GoodsID:     info.GoodsID,
		SkuID:       skuInfo.SkuID,
		GoodsName:   info.GoodsName,
		SkuName:     skuInfo.SkuName,
		ImageURL:    skuInfo.ImageURL,
		RetailPrice: skuInfo.RetailPrice,
		Quantity:    req.Quantity,
	})
	if err != nil {
		zap.L().Error("save cart info failed", zap.Error(err))
		return nil, errors.New("保存购物车失败")
	}

	return cartInfo, nil
}
