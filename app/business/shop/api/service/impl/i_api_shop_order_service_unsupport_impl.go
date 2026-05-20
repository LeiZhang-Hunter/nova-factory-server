//go:build !erp
// +build !erp

package impl

import (
	"gorm.io/gorm"
	activityDao "nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/service"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/datasource/cache"
)

// IApiShopOrderServiceImpl 没有erp模块的时候加载erp
type IApiShopOrderServiceImpl struct {
	cache           cache.Cache
	db              *gorm.DB
	orderDao        dao.IApiShopOrderDao
	userDao         dao.IApiShopWechatUserDao
	addressDao      dao.IApiShopAddressDao
	cartDao         dao.IApiShopCartDao
	seckillDao      activityDao.IShopSeckillDao
	combDao         activityDao.IShopCombinationDao
	goodsDao        dao.IApiShopGoodsDao
	skuDao          dao.IApiShopSkuDao
	discountService discountservice.IDiscountCalculateService
}

// NewIApiShopOrderServiceImpl 创建订单服务实现。
func NewIApiShopOrderServiceImpl(
	cache cache.Cache,
	db *gorm.DB,
	orderDao dao.IApiShopOrderDao,
	userDao dao.IApiShopWechatUserDao,
	addressDao dao.IApiShopAddressDao,
	cartDao dao.IApiShopCartDao,
	seckillDao activityDao.IShopSeckillDao,
	combDao activityDao.IShopCombinationDao,
	goodsDao dao.IApiShopGoodsDao,
	skuDao dao.IApiShopSkuDao,
	discountService discountservice.IDiscountCalculateService,
) service.IApiShopOrderService {
	return &IApiShopOrderServiceImpl{
		cache:           cache,
		db:              db,
		orderDao:        orderDao,
		userDao:         userDao,
		addressDao:      addressDao,
		cartDao:         cartDao,
		seckillDao:      seckillDao,
		combDao:         combDao,
		goodsDao:        goodsDao,
		skuDao:          skuDao,
		discountService: discountService,
	}
}
