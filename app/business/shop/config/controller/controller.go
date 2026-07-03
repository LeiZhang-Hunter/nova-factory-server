package shopcontroller

import "github.com/google/wire"

type Controller struct {
	WechatConfig             *WechatConfig
	Logistics                *Logistics
	ShopErpIntegrationConfig *ShopErpIntegrationConfig
	LogisticsConfig          *LogisticsConfig
	EnterpriseAccount        *EnterpriseAccountConfig
}

var ProviderSet = wire.NewSet(
	NewWechatConfig, NewLogistics, NewShopErpIntegrationConfig, NewLogisticsConfig, NewEnterpriseAccountConfig,
	wire.Struct(new(Controller), "*"),
)
