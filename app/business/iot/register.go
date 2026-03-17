//go:build iot
// +build iot

package iot

import (
	"nova-factory-server/app/business/iot/alert/alertController"
	"nova-factory-server/app/business/iot/alert/alertDao/alertDaoImpl"
	"nova-factory-server/app/business/iot/alert/alertService/alertServiceImpl"
	"nova-factory-server/app/business/iot/asset/building/buildingController"
	"nova-factory-server/app/business/iot/asset/building/buildingDao/buildingDaoImpl"
	"nova-factory-server/app/business/iot/asset/building/buildingService/buildingServiceImpl"
	"nova-factory-server/app/business/iot/asset/camera/cameraController"
	"nova-factory-server/app/business/iot/asset/camera/cameraDao/cameraDaoImpl"
	"nova-factory-server/app/business/iot/asset/camera/cameraService/cameraServiceImpl"
	"nova-factory-server/app/business/iot/asset/device/deviceController"
	"nova-factory-server/app/business/iot/asset/device/deviceDao/deviceDaoImpl"
	"nova-factory-server/app/business/iot/asset/device/deviceService/deviceServiceImpl"
	"nova-factory-server/app/business/iot/asset/material/materialController"
	"nova-factory-server/app/business/iot/asset/material/materialDao/materialDaoImpl"
	"nova-factory-server/app/business/iot/asset/material/materialService/materialServiceImpl"
	"nova-factory-server/app/business/iot/asset/resource/resourceController"
	"nova-factory-server/app/business/iot/asset/resource/resourceDao/resourceDaoImpl"
	"nova-factory-server/app/business/iot/asset/resource/resourceService/resourceServiceImpl"
	"nova-factory-server/app/business/iot/configuration/configurationController"
	"nova-factory-server/app/business/iot/configuration/configurationDao/configurationDaoImpl"
	"nova-factory-server/app/business/iot/configuration/configurationService/configurationServiceImpl"
	"nova-factory-server/app/business/iot/craft/craftRouteController"
	"nova-factory-server/app/business/iot/craft/craftRouteDao/craftRouteDaoImpl"
	"nova-factory-server/app/business/iot/craft/craftRouteService/craftRouteServiceImpl"
	"nova-factory-server/app/business/iot/daemonize/daemonizeController"
	"nova-factory-server/app/business/iot/daemonize/daemonizeDao/daemonizeDaoImpl"
	"nova-factory-server/app/business/iot/daemonize/daemonizeService/daemonizeServiceImpl"
	"nova-factory-server/app/business/iot/dashboard/dashboardController"
	"nova-factory-server/app/business/iot/dashboard/dashboardDao/dashboardDaoImpl"
	"nova-factory-server/app/business/iot/dashboard/dashboardService/dashboardServiceImpl"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorController"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorDao/deviceMonitorDaoImpl"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorService/deviceMonitorServiceImpl"
	homeController "nova-factory-server/app/business/iot/home/controller"
	"nova-factory-server/app/business/iot/home/homeService/homeServiceImpl"
	"nova-factory-server/app/business/iot/metric/device/metricController"
	"nova-factory-server/app/business/iot/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/business/iot/metric/device/metricService/metricServiceImpl"
	iotSystemControllerImpl "nova-factory-server/app/business/iot/system/controller"
	iotSystemDaoImpl "nova-factory-server/app/business/iot/system/dao/systemDaoImpl"
	iotSystemServiceImpl "nova-factory-server/app/business/iot/system/service/systemServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	deviceDaoImpl.ProviderSet,
	deviceServiceImpl.ProviderSet,
	deviceController.ProviderSet,

	materialDaoImpl.ProviderSet,
	materialServiceImpl.ProviderSet,
	materialController.ProviderSet,

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

	resourceController.ProviderSet,
	resourceServiceImpl.ProviderSet,
	resourceDaoImpl.ProviderSet,

	homeServiceImpl.ProviderSet,
	homeController.ProviderSet,

	configurationController.ProviderSet,
	configurationServiceImpl.ProviderSet,
	configurationDaoImpl.ProviderSet,

	iotSystemControllerImpl.ProviderSet,
	iotSystemServiceImpl.ProviderSet,
	iotSystemDaoImpl.ProviderSet,

	cameraController.ProviderSet,
	cameraServiceImpl.ProviderSet,
	cameraDaoImpl.ProviderSet,

	NewGinEngine,
)
