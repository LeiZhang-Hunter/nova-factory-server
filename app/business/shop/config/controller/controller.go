package shopcontroller

import "github.com/google/wire"

type Controller struct {
	WechatConfig             *WechatConfig
	Logistics                *Logistics
	ShopErpIntegrationConfig *ShopErpIntegrationConfig
}

var ProviderSet = wire.NewSet(NewWechatConfig, NewLogistics, NewShopErpIntegrationConfig, wire.Struct(new(Controller), "*"))
