package impl

import (
	"gorm.io/gorm"
	activityDao "nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/service"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	orderDao "nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/datasource/cache"
)

// IApiShopOrderServiceImpl 提供订单相关的业务实现。
type IApiShopOrderServiceImpl struct {
	cache           cache.Cache
	db              *gorm.DB
	apiOrderDao     dao.IApiShopOrderDao
	orderDao        orderDao.IOrderDao
	userDao         dao.IApiShopWechatUserDao
	addressDao      dao.IApiShopAddressDao
	cartDao         dao.IApiShopCartDao
	seckillDao      activityDao.IShopSeckillDao
	combDao         activityDao.IShopCombinationDao
	goodsDao        dao.IApiShopGoodsDao
	skuDao          dao.IApiShopSkuDao
	configDao       dao.IApiShopSysConfigDao
	discountService discountservice.IDiscountCalculateService
	orderSync       *shopOrderSyncService
}

// NewIApiShopOrderServiceImpl 创建订单服务实现。
func NewIApiShopOrderServiceImpl(
	cache cache.Cache,
	db *gorm.DB,
	apiOrderDao dao.IApiShopOrderDao,
	orderDao orderDao.IOrderDao,
	userDao dao.IApiShopWechatUserDao,
	addressDao dao.IApiShopAddressDao,
	cartDao dao.IApiShopCartDao,
	seckillDao activityDao.IShopSeckillDao,
	combDao activityDao.IShopCombinationDao,
	goodsDao dao.IApiShopGoodsDao,
	skuDao dao.IApiShopSkuDao,
	configDao dao.IApiShopSysConfigDao,
	discountService discountservice.IDiscountCalculateService,
	orderSync *shopOrderSyncService,
) service.IApiShopOrderService {
	return &IApiShopOrderServiceImpl{
		cache:           cache,
		db:              db,
		apiOrderDao:     apiOrderDao,
		orderDao:        orderDao,
		userDao:         userDao,
		addressDao:      addressDao,
		cartDao:         cartDao,
		seckillDao:      seckillDao,
		combDao:         combDao,
		goodsDao:        goodsDao,
		skuDao:          skuDao,
		configDao:       configDao,
		discountService: discountService,
		orderSync:       orderSync,
	}
}
