package shopdaoimpl

import "nova-factory-server/app/business/shop/product/shopdao"

type IShopOrderDaoImpl struct {
}

func NewIShopOrderDaoImpl() shopdao.IShopOrderDao {
	return &IShopOrderDaoImpl{}
}
