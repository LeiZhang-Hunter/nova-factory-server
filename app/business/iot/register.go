//go:build iot
// +build iot

package iot

import (
	"nova-factory-server/app/business/iot/alert/alertcontroller"
	"nova-factory-server/app/business/iot/alert/alertdao/alertdaoimpl"
	"nova-factory-server/app/business/iot/alert/alertservice/alertserviceimpl"
	"nova-factory-server/app/business/iot/asset/building/buildingcontroller"
	"nova-factory-server/app/business/iot/asset/building/buildingdao/buildingdaoimpl"
	"nova-factory-server/app/business/iot/asset/building/buildingservice/buildingserviceimpl"
	"nova-factory-server/app/business/iot/asset/camera/cameracontroller"
	"nova-factory-server/app/business/iot/asset/camera/cameradao/cameraDaoImpl"
	"nova-factory-server/app/business/iot/asset/camera/cameraservice/cameraServiceImpl"
	"nova-factory-server/app/business/iot/asset/device/devicecontroller"
	"nova-factory-server/app/business/iot/asset/device/devicedao/devicedaoImpl"
	"nova-factory-server/app/business/iot/asset/device/deviceservice/deviceserviceimpl"
	"nova-factory-server/app/business/iot/asset/material/materialcontroller"
	"nova-factory-server/app/business/iot/asset/material/materialdao/materialdaoimpl"
	"nova-factory-server/app/business/iot/asset/material/materialservice/materialserviceimpl"
	"nova-factory-server/app/business/iot/asset/resource/resourcecontroller"
	"nova-factory-server/app/business/iot/asset/resource/resourcedao/resourcedaoimpl"
	"nova-factory-server/app/business/iot/asset/resource/resourceservice/resourceserviceimpl"
	"nova-factory-server/app/business/iot/configuration/configurationcontroller"
	"nova-factory-server/app/business/iot/configuration/configurationdao/configurationDaoImpl"
	"nova-factory-server/app/business/iot/configuration/configurationservice/configurationServiceImpl"
	"nova-factory-server/app/business/iot/craft/craftroutecontroller"
	"nova-factory-server/app/business/iot/craft/craftroutedao/craftroutedaoimpl"
	"nova-factory-server/app/business/iot/craft/craftrouteservice/craftrouteserviceimpl"
	"nova-factory-server/app/business/iot/daemonize/daemonizecontroller"
	"nova-factory-server/app/business/iot/daemonize/daemonizedao/daemonizedaoimpl"
	"nova-factory-server/app/business/iot/daemonize/daemonizeservice/daemonizeserviceimpl"
	"nova-factory-server/app/business/iot/dashboard/dashboardcontroller"
	"nova-factory-server/app/business/iot/dashboard/dashboarddao/dashboarddaoimpl"
	"nova-factory-server/app/business/iot/dashboard/dashboardservice/dashboardserviceimpl"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorcontroller"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao/deviceMonitorDaoImpl"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorservice/deviceMonitorServiceImpl"
	homeController "nova-factory-server/app/business/iot/home/controller"
	"nova-factory-server/app/business/iot/home/homeservice/homeserviceimpl"
	"nova-factory-server/app/business/iot/metric/device/metriccontroller"
	"nova-factory-server/app/business/iot/metric/device/metricdao/metricdaoimpl"
	"nova-factory-server/app/business/iot/metric/device/metricservice/metricserviceimpl"
	iotSystemControllerImpl "nova-factory-server/app/business/iot/system/controller"
	iotSystemDaoImpl "nova-factory-server/app/business/iot/system/dao/systemdaoimpl"
	iotSystemServiceImpl "nova-factory-server/app/business/iot/system/service/systemserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	devicedaoImpl.ProviderSet,
	deviceserviceimpl.ProviderSet,
	devicecontroller.ProviderSet,

	materialdaoimpl.ProviderSet,
	materialserviceimpl.ProviderSet,
	materialcontroller.ProviderSet,

	craftroutedaoimpl.ProviderSet,
	craftrouteserviceimpl.ProviderSet,
	craftroutecontroller.ProviderSet,

	metricdaoimpl.ProviderSet,
	metricserviceimpl.ProviderSet,
	metriccontroller.ProviderSet,

	daemonizedaoimpl.ProviderSet,
	daemonizeserviceimpl.ProviderSet,
	daemonizecontroller.ProviderSet,

	deviceMonitorDaoImpl.ProviderSet,
	deviceMonitorServiceImpl.ProviderSet,
	devicemonitorcontroller.ProviderSet,

	alertdaoimpl.ProviderSet,
	alertserviceimpl.ProviderSet,
	alertcontroller.ProviderSet,

	buildingdaoimpl.ProviderSet,
	buildingserviceimpl.ProviderSet,
	buildingcontroller.ProviderSet,

	dashboarddaoimpl.ProviderSet,
	dashboardserviceimpl.ProviderSet,
	dashboardcontroller.ProviderSet,

	resourcecontroller.ProviderSet,
	resourceserviceimpl.ProviderSet,
	resourcedaoimpl.ProviderSet,

	homeserviceimpl.ProviderSet,
	homeController.ProviderSet,

	configurationcontroller.ProviderSet,
	configurationServiceImpl.ProviderSet,
	configurationDaoImpl.ProviderSet,

	iotSystemControllerImpl.ProviderSet,
	iotSystemServiceImpl.ProviderSet,
	iotSystemDaoImpl.ProviderSet,

	cameracontroller.ProviderSet,
	cameraServiceImpl.ProviderSet,
	cameraDaoImpl.ProviderSet,

	GinProviderSet,
)
