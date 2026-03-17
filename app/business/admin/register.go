package admin

import (
	"nova-factory-server/app/business/admin/monitor/monitorController"
	"nova-factory-server/app/business/admin/monitor/monitorDao/monitorDaoImpl"
	"nova-factory-server/app/business/admin/monitor/monitorService/monitorServiceImpl"
	"nova-factory-server/app/business/admin/product/productController"
	"nova-factory-server/app/business/admin/product/productDao/productDaoImpl"
	"nova-factory-server/app/business/admin/product/productService/productServiceImpl"
	"nova-factory-server/app/business/admin/system/systemController"
	"nova-factory-server/app/business/admin/system/systemDao/systemDaoImpl"
	"nova-factory-server/app/business/admin/system/systemService/systemServiceImpl"
	"nova-factory-server/app/business/admin/tool/toolController"
	"nova-factory-server/app/business/admin/tool/toolDao/toolDaoImpl"
	"nova-factory-server/app/business/admin/tool/toolService/toolServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	toolDaoImpl.ProviderSet,
	toolServiceImpl.ProviderSet,
	toolController.ProviderSet,

	systemDaoImpl.ProviderSet,
	systemServiceImpl.ProviderSet,
	systemController.ProviderSet,

	monitorDaoImpl.ProviderSet,
	monitorServiceImpl.ProviderSet,
	monitorController.ProviderSet,

	productDaoImpl.ProviderSet,
	productServiceImpl.ProviderSet,
	productController.ProviderSet,
)
