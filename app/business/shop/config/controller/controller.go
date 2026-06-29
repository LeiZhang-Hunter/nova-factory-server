package shopcontroller

import "github.com/google/wire"

type Controller struct {
	WechatConfig             *WechatConfig
	Logistics                *Logistics
	ShopErpIntegrationConfig *ShopErpIntegrationConfig
	LogisticsConfig          *LogisticsConfig
}

var ProviderSet = wire.NewSet(
	NewWechatConfig, NewLogistics, NewShopErpIntegrationConfig, NewLogisticsConfig,
	wire.Struct(new(Controller), "*"),
)
