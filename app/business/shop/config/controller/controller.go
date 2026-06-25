package shopcontroller

import "github.com/google/wire"

type Controller struct {
	WechatConfig *WechatConfig
	Logistics    *Logistics
}

var ProviderSet = wire.NewSet(NewWechatConfig, NewLogistics, wire.Struct(new(Controller), "*"))
