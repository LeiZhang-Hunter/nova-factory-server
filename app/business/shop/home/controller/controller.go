package controller

import "github.com/google/wire"

// ProviderSet 首页模块控制器 Provider。
var ProviderSet = wire.NewSet(NewHomeModule, NewHomeModuleItem, wire.Struct(new(Controller), "*"))

// Controller 首页模块控制器聚合。
type Controller struct {
	HomeModule     *HomeModule
	HomeModuleItem *HomeModuleItem
}
