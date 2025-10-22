//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/business/ai/aiDataSetDao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/ai/aiDataSetService/aiDataSetServiceImpl"
	"nova-factory-server/app/business/alert/alertController"
	"nova-factory-server/app/business/alert/alertDao/alertDaoImpl"
	"nova-factory-server/app/business/alert/alertService/alertServiceImpl"
	"nova-factory-server/app/business/asset/building/buildingController"
	"nova-factory-server/app/business/asset/building/buildingDao/buildingDaoImpl"
	"nova-factory-server/app/business/asset/building/buildingService/buildingServiceImpl"
	"nova-factory-server/app/business/asset/device/deviceController"
	"nova-factory-server/app/business/asset/device/deviceDao/deviceDaoImpl"
	"nova-factory-server/app/business/asset/device/deviceService/deviceServiceImpl"
	"nova-factory-server/app/business/asset/material/materialController"
	"nova-factory-server/app/business/asset/material/materialDao/materialDaoImpl"
	"nova-factory-server/app/business/asset/material/materialService/materialServiceImpl"
	"nova-factory-server/app/business/craft/craftRouteController"
	"nova-factory-server/app/business/craft/craftRouteDao/craftRouteDaoImpl"
	"nova-factory-server/app/business/craft/craftRouteService/craftRouteServiceImpl"
	"nova-factory-server/app/business/daemonize/daemonizeController"
	"nova-factory-server/app/business/daemonize/daemonizeDao/daemonizeDaoImpl"
	"nova-factory-server/app/business/daemonize/daemonizeService/daemonizeServiceImpl"
	"nova-factory-server/app/business/dashboard/dashboardController"
	"nova-factory-server/app/business/dashboard/dashboardDao/dashboardDaoImpl"
	"nova-factory-server/app/business/dashboard/dashboardService/dashboardServiceImpl"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorController"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao/deviceMonitorDaoImpl"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService/deviceMonitorServiceImpl"
	"nova-factory-server/app/business/metric/device/metricController"
	"nova-factory-server/app/business/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/business/metric/device/metricService/metricServiceImpl"
	"nova-factory-server/app/business/monitor/monitorController"
	"nova-factory-server/app/business/monitor/monitorDao/monitorDaoImpl"
	"nova-factory-server/app/business/monitor/monitorService/monitorServiceImpl"
	"nova-factory-server/app/business/system/systemController"
	"nova-factory-server/app/business/system/systemDao/systemDaoImpl"
	"nova-factory-server/app/business/system/systemService/systemServiceImpl"
	"nova-factory-server/app/business/tool/toolController"
	"nova-factory-server/app/business/tool/toolDao/toolDaoImpl"
	"nova-factory-server/app/business/tool/toolService/toolServiceImpl"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"
)

func wireApp() (*gin.Engine, func(), error) {
	panic(wire.Build(
		toolDaoImpl.ProviderSet,
		toolServiceImpl.ProviderSet,
		toolController.ProviderSet,

		systemDaoImpl.ProviderSet,
		systemServiceImpl.ProviderSet,
		systemController.ProviderSet,

		monitorDaoImpl.ProviderSet,
		monitorServiceImpl.ProviderSet,
		monitorController.ProviderSet,

		deviceDaoImpl.ProviderSet,
		deviceServiceImpl.ProviderSet,
		deviceController.ProviderSet,

		materialDaoImpl.ProviderSet,
		materialServiceImpl.ProviderSet,
		materialController.ProviderSet,

		aiDataSetDaoImpl.ProviderSet,
		aiDataSetServiceImpl.ProviderSet,
		aiDataSetController.ProviderSet,

		craftRouteDaoImpl.ProviderSet,
		craftRouteServiceImpl.ProviderSet,
		craftRouteController.ProviderSet,

		metricDaoIMpl.ProviderSet,
		metricServiceImpl.ProviderSet,
		metricController.ProviderSet,

		daemonizeDaoImpl.ProviderSet,
		daemonizeServiceImpl.ProviderSet,
		daemonizeController.ProviderSet,

		deviceMonitorDaoImpl.ProviderSet,
		deviceMonitorServiceImpl.ProviderSet,
		deviceMonitorController.ProviderSet,

		alertDaoImpl.ProviderSet,
		alertServiceImpl.ProviderSet,
		alertController.ProviderSet,

		buildingDaoImpl.ProviderSet,
		buildingServiceImpl.ProviderSet,
		buildingController.ProviderSet,

		dashboardDaoImpl.ProviderSet,
		dashboardServiceImpl.ProviderSet,
		dashboardController.ProviderSet,

		datasource.ProviderSet,
		routes.ProviderSet,
	))
}
