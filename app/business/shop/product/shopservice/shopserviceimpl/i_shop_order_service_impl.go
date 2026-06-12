package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopservice"
)

type IShopOrderServiceImpl struct {
	orderDao shopdao.IShopOrderDao
}

func NewIShopOrderServiceImpl(orderDao shopdao.IShopOrderDao) shopservice.IShopOrderService {
	return &IShopOrderServiceImpl{
		orderDao: orderDao,
	}
}
