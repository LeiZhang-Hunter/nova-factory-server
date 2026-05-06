package shopcontroller

import "github.com/google/wire"

type Controller struct {
	WechatConfig *WechatConfig
}

var ProviderSet = wire.NewSet(NewWechatConfig, wire.Struct(new(Controller), "*"))
