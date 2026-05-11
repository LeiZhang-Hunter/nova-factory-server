package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	api "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/business/shop/product/shopdao"
	userDao "nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IApiShopCartServiceImpl struct {
	goodsDao     shopdao.IShopGoodsDao
	goodsSkuDao  shopdao.IShopSkuDao
	shopUserDao  userDao.IShopUserDao
	iShopCartDao dao.IApiShopCartDao
}

// NewIApiShopCartServiceImpl 购物车服务
func NewIApiShopCartServiceImpl(goodsDao shopdao.IShopGoodsDao,
	goodsSkuDao shopdao.IShopSkuDao,
	shopUserDao userDao.IShopUserDao,
	iShopCartDao dao.IApiShopCartDao) service.IApiShopCartService {
	return &IApiShopCartServiceImpl{
		goodsDao:     goodsDao,
		goodsSkuDao:  goodsSkuDao,
		iShopCartDao: iShopCartDao,
		shopUserDao:  shopUserDao,
	}
}

// GenCart 生成购物车，购物车下单
func (s *IApiShopCartServiceImpl) GenCart(c *gin.Context, req *api.CartSetDataReq) (*api.CartDto, error) {
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

	userInfo, err := s.shopUserDao.GetByUserID(c, userId)
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
		GoodsID:     info.ID,
		SkuID:       skuInfo.ID,
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

// List 查询用户购物车列表。
func (s *IApiShopCartServiceImpl) List(c *gin.Context) ([]*api.CartDto, error) {
	userID := baizeContext.GetUserId(c)
	//var userID int64 = 486436412988592128
	//if userID == 0 {
	//	return nil, errors.New("用户ID不能为空")
	//}
	ret, err := s.iShopCartDao.List(c, userID)
	if err != nil {
		return nil, err
	}

	skuIDs := make([]int64, 0, len(ret))
	skuIDSet := make(map[int64]struct{}, len(ret))
	for _, item := range ret {
		if _, ok := skuIDSet[item.SkuID]; ok {
			continue
		}
		skuIDSet[item.SkuID] = struct{}{}
		skuIDs = append(skuIDs, item.SkuID)
	}

	skuQuantityMap := make(map[int64]int64, len(skuIDs))
	if len(skuIDs) > 0 {
		skuList, skuErr := s.goodsSkuDao.ListByIDs(c, skuIDs)
		if skuErr != nil {
			zap.L().Error("batch get sku list failed", zap.Error(skuErr))
			return nil, errors.New("读取商品库存失败")
		}
		for _, skuInfo := range skuList {
			skuQuantityMap[int64(skuInfo.ID)] = skuInfo.Quantity
		}
	}

	for k, v := range ret {
		ret[k].ImageURL = fileUtils.BuildAbsoluteURL(c, v.ImageURL)
		skuQuantity, ok := skuQuantityMap[v.SkuID]
		if !ok {
			ret[k].IsStockEnough = false
			continue
		}
		ret[k].IsStockEnough = skuQuantity >= v.Quantity
	}
	return ret, nil
}

func (s *IApiShopCartServiceImpl) Remove(c *gin.Context, ids []string) error {
	return s.iShopCartDao.Remove(c, ids)
}
