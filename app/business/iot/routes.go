//go:build iot
// +build iot

package iot

import (
	"github.com/google/wire"
	"nova-factory-server/app/business/iot/alert/alertcontroller"
	"nova-factory-server/app/business/iot/asset/building/buildingcontroller"
	"nova-factory-server/app/business/iot/asset/camera/cameracontroller"
	"nova-factory-server/app/business/iot/asset/device/devicecontroller"
	"nova-factory-server/app/business/iot/asset/material/materialcontroller"
	"nova-factory-server/app/business/iot/asset/resource/resourcecontroller"
	"nova-factory-server/app/business/iot/configuration/configurationcontroller"
	"nova-factory-server/app/business/iot/craft/craftroutecontroller"
	"nova-factory-server/app/business/iot/daemonize/daemonizecontroller"
	"nova-factory-server/app/business/iot/dashboard/dashboardcontroller"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorcontroller"
	homeController "nova-factory-server/app/business/iot/home/controller"
	"nova-factory-server/app/business/iot/metric/device/metriccontroller"
	iotSystemControllerImpl "nova-factory-server/app/business/iot/system/controller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	materialC *materialcontroller.Material,
	craft *craftroutecontroller.CraftRoute,
	metric *metriccontroller.MetricServer,
	controller *daemonizecontroller.DaemonizeServer,
	deviceMonitor *devicemonitorcontroller.DeviceMonitorController,
	alert *alertcontroller.Controller,
	building buildingcontroller.Controller,
	dashboard dashboardcontroller.Controller,
	resource *resourcecontroller.ResourceController,
	home *homeController.Home,
	configuration *configurationcontroller.Controller,
	iotSystem iotSystemControllerImpl.System,
	camera cameracontroller.CameraController,
	dc *devicecontroller.Device,
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
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
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

	}

	// mpc
	{
		iotSystem.Electric.PrivateMcpRoutes(app.McpServer)
		dc.Info.PrivateMcpRoutes(app.McpServer)  //资产管理---设备模块
		dc.Group.PrivateMcpRoutes(app.McpServer) //资产管理---设备分组
		dc.Template.PrivateMcpRoutes(app.McpServer)
		dc.TemplateData.PrivateMcpRoutes(app.McpServer)
		dc.DeviceSubject.PrivateMcpRoutes(app.McpServer)
		dc.DeviceCheckPlan.PrivateMcpRoutes(app.McpServer)
		dc.DeviceCheckMachinery.PrivateMcpRoutes(app.McpServer)
		dc.DeviceCheckSubject.PrivateMcpRoutes(app.McpServer)

		materialC.Material.PrivateMcpRoutes(app.McpServer)    //资产管理---物料管理
		resource.ResourceFile.PrivateMcpRoutes(app.McpServer) //资产管理---资料管理

		craft.CraftRoute.PrivateMcpRoutes(app.McpServer)     //工艺路线
		craft.Process.PrivateMcpRoutes(app.McpServer)        //工序设置
		craft.ProcessContext.PrivateMcpRoutes(app.McpServer) //工序内容
		craft.RouteProcess.PrivateMcpRoutes(app.McpServer)   //工艺组成
		craft.Schedule.PrivateMcpRoutes(app.McpServer)

		controller.IotAgent.PrivateMcpRoutes(app.McpServer)
		controller.Config.PrivateMcpRoutes(app.McpServer)

		deviceMonitor.DeviceMonitor.PrivateMcpRoutes(app.McpServer)
		deviceMonitor.DeviceReport.PrivateMcpRoutes(app.McpServer)
		deviceMonitor.DeviceUtilization.PrivateMcpRoutes(app.McpServer)
		deviceMonitor.ControlLog.PrivateMcpRoutes(app.McpServer)

		alert.AlertTemplate.PrivateMcpRoutes(app.McpServer)
		alert.Alert.PrivateMcpRoutes(app.McpServer)
		alert.AlertLog.PrivateMcpRoutes(app.McpServer)
		alert.AlertAction.PrivateMcpRoutes(app.McpServer)
		alert.AlertAiReason.PrivateMcpRoutes(app.McpServer)

		building.Building.PrivateMcpRoutes(app.McpServer)
		building.Floor.PrivateMcpRoutes(app.McpServer)

		dashboard.Dashboard.PrivateMcpRoutes(app.McpServer)
		dashboard.Data.PrivateMcpRoutes(app.McpServer)

		home.PrivateMcpRoutes(app.McpServer)

		// 生产管理
		camera.Camera.PrivateMcpRoutes(app.McpServer)

		// 组态管理
		configuration.Configuration.PrivateMcpRoutes(app.McpServer)
	}

	//grpc
	s := app.GrpcServer
	controller.Daemonize.PrivateRoutes(s)
	controller.Config.PrivateGrpcRoutes(s)
	deviceMonitor.DeviceControl.PrivateRoutes(s)
	deviceMonitor.CameraGrpc.PrivateRoutes(s)

	return &Iot{}
}
