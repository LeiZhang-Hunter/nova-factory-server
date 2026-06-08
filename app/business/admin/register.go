package admin

import (
	basicscontroller "nova-factory-server/app/business/admin/basics/controller"
	basicsdaoimpl "nova-factory-server/app/business/admin/basics/dao/impl"
	basicsserviceimpl "nova-factory-server/app/business/admin/basics/service/impl"
	"nova-factory-server/app/business/admin/monitor/monitorcontroller"
	"nova-factory-server/app/business/admin/monitor/monitordao/monitorDaoImpl"
	"nova-factory-server/app/business/admin/monitor/monitorservice/monitorServiceImpl"
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

	basicsdaoimpl.ProviderSet,
	basicsserviceimpl.ProviderSet,
	basicscontroller.ProviderSet,

	systemdaoimpl.ProviderSet,
	systemServiceImpl.ProviderSet,
	systemcontroller.ProviderSet,

	monitorDaoImpl.ProviderSet,
	monitorServiceImpl.ProviderSet,
	monitorcontroller.ProviderSet,

	GinProviderSet,
)
