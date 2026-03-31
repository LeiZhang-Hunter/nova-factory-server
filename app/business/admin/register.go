package admin

import (
	"nova-factory-server/app/business/admin/monitor/monitorcontroller"
	"nova-factory-server/app/business/admin/monitor/monitordao/monitorDaoImpl"
	"nova-factory-server/app/business/admin/monitor/monitorservice/monitorServiceImpl"
	"nova-factory-server/app/business/admin/product/productcontroller"
	"nova-factory-server/app/business/admin/product/productdao/productDaoImpl"
	"nova-factory-server/app/business/admin/product/productservice/productServiceImpl"
	"nova-factory-server/app/business/admin/system/systemcontroller"
	"nova-factory-server/app/business/admin/system/systemdao/systemdaoimpl"
	"nova-factory-server/app/business/admin/system/systemservice/systemServiceImpl"
	"nova-factory-server/app/business/admin/tool/toolcontroller"
	"nova-factory-server/app/business/admin/tool/tooldao/tooldaoimpl"
	"nova-factory-server/app/business/admin/tool/toolservice/toolserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	tooldaoimpl.ProviderSet,
	toolserviceimpl.ProviderSet,
	toolcontroller.ProviderSet,

	systemdaoimpl.ProviderSet,
	systemServiceImpl.ProviderSet,
	systemcontroller.ProviderSet,

	monitorDaoImpl.ProviderSet,
	monitorServiceImpl.ProviderSet,
	monitorcontroller.ProviderSet,

	productDaoImpl.ProviderSet,
	productServiceImpl.ProviderSet,
	productcontroller.ProviderSet,
	GinProviderSet,
)
