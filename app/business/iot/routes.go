//go:build iot
// +build iot

package iot

import (
	"nova-factory-server/app/business/iot/alert/alertController"
	"nova-factory-server/app/business/iot/asset/building/buildingController"
	"nova-factory-server/app/business/iot/asset/camera/cameraController"
	"nova-factory-server/app/business/iot/asset/device/deviceController"
	"nova-factory-server/app/business/iot/asset/material/materialController"
	"nova-factory-server/app/business/iot/asset/resource/resourceController"
	"nova-factory-server/app/business/iot/configuration/configurationController"
	"nova-factory-server/app/business/iot/craft/craftRouteController"
	"nova-factory-server/app/business/iot/daemonize/daemonizeController"
	"nova-factory-server/app/business/iot/dashboard/dashboardController"
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorController"
	homeController "nova-factory-server/app/business/iot/home/controller"
	"nova-factory-server/app/business/iot/metric/device/metricController"
	iotSystemControllerImpl "nova-factory-server/app/business/iot/system/controller"
	"nova-factory-server/app/daemonize"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	materialC *materialController.Material,
	craft *craftRouteController.CraftRoute,
	metric *metricController.MetricServer,
	controller *daemonizeController.DaemonizeServer,
	deviceMonitor *deviceMonitorController.DeviceMonitorController,
	alert *alertController.Controller,
	building buildingController.Controller,
	dashboard dashboardController.Controller,
	resource *resourceController.ResourceController,
	home *homeController.Home,
	configuration *configurationController.Controller,
	iotSystem iotSystemControllerImpl.System,
	camera cameraController.CameraController,
	dc *deviceController.Device,
) *Iot {
	r := app.Engine
	group := r.Group("")
	//不做鉴权的
	{
		dc.Info.PublicRoutes(group)
		deviceMonitor.DeviceUtilization.PublicRoutes(group)
		resource.ResourceFile.PublicRoutes(group) //资产管理---资料管理
		building.Building.PublicRoutes(group)
		controller.Config.PublicRoutes(group) //注册Agent公共配置接口
		alert.AlertLog.PublicRoutes(group)
		metric.Metric.PublicRoutes(group)
		craft.Schedule.PublicRoutes(group)
		deviceMonitor.DeviceMonitor.PublicRoutes(group)
	}

	// 做鉴权的
	{
		iotSystem.Electric.PrivateRoutes(group)
		dc.Info.PrivateRoutes(group)  //资产管理---设备模块
		dc.Group.PrivateRoutes(group) //资产管理---设备分组
		dc.Template.PrivateRoutes(group)
		dc.TemplateData.PrivateRoutes(group)
		dc.DeviceSubject.PrivateRoutes(group)
		dc.DeviceCheckPlan.PrivateRoutes(group)
		dc.DeviceCheckMachinery.PrivateRoutes(group)
		dc.DeviceCheckSubject.PrivateRoutes(group)

		materialC.Material.PrivateRoutes(group)    //资产管理---物料管理
		resource.ResourceFile.PrivateRoutes(group) //资产管理---资料管理

		craft.CraftRoute.PrivateRoutes(group)     //工艺路线
		craft.Process.PrivateRoutes(group)        //工序设置
		craft.ProcessContext.PrivateRoutes(group) //工序内容
		craft.RouteProcess.PrivateRoutes(group)   //工艺组成
		craft.Schedule.PrivateRoutes(group)

		metric.Metric.PrivateRoutes(group) //设备指标
		controller.IotAgent.PrivateRoutes(group)
		controller.Config.PrivateRoutes(group)

		deviceMonitor.DeviceMonitor.PrivateRoutes(group)
		deviceMonitor.DeviceReport.PrivateRoutes(group)
		deviceMonitor.DeviceUtilization.PrivateRoutes(group)
		deviceMonitor.ControlLog.PrivateRoutes(group)

		alert.AlertTemplate.PrivateRoutes(group)
		alert.Alert.PrivateRoutes(group)
		alert.AlertLog.PrivateRoutes(group)
		alert.AlertAction.PrivateRoutes(group)
		alert.AlertAiReason.PrivateRoutes(group)

		building.Building.PrivateRoutes(group)
		building.Floor.PrivateRoutes(group)

		dashboard.Dashboard.PrivateRoutes(group)
		dashboard.Data.PrivateRoutes(group)

		home.PrivateRoutes(group)

		// 生产管理
		camera.Camera.PrivateRoutes(group)

		// 组态管理
		configuration.Configuration.PrivateRoutes(group)

		alert.Runner.Run()
	}

	// mpc
	{
		dc.Info.PrivateMcpRoutes(app.McpServer) //资产管理---设备模块
		deviceMonitor.DeviceMonitor.PrivateMcpRoutes(app.McpServer)
	}

	//grpc
	s := daemonize.CreateGRpcServer()
	controller.Daemonize.PrivateRoutes(s)
	deviceMonitor.DeviceControl.PrivateRoutes(s)
	deviceMonitor.CameraGrpc.PrivateRoutes(s)
	s.Start()
	return &Iot{}
}
