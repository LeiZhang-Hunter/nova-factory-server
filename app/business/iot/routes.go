//go:build iot
// +build iot

package iot

import (
	"fmt"
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
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/daemonize"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"
	"time"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	iotdb *iotdb.IotDb,
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

	session, err := iotdb.GetSession()
	if err != nil {
		zap.L().Error("iotdb.GetSession()", zap.Error(err))
		panic(err)
	}
	defer iotdb.PutSession(session)
	for {
		statement, err := session.ExecuteStatement("count databases root.device")
		if err != nil {
			zap.L().Error("execute statement", zap.Error(err))
			time.Sleep(1 * time.Second)
			continue
		}
		hasDatabase, err := statement.Next()
		if err != nil {
			zap.L().Error("get hasDatabase", zap.Error(err))
			time.Sleep(1 * time.Second)
			continue
		}

		count := statement.GetInt32("count")
		if count < 1 {
			session.ExecuteStatement("create database root.device")
			time.Sleep(1 * time.Second)
			continue
		}

		statement, err = session.ExecuteStatement("count databases root.run_status_device")
		if err != nil {
			time.Sleep(1 * time.Second)
			return nil
		}
		hasDatabase, err = statement.Next()
		if err != nil {
			zap.L().Error("get hasDatabase", zap.Error(err))
			time.Sleep(1 * time.Second)
			continue
		}

		if !hasDatabase {
			time.Sleep(1 * time.Second)
			continue
		}

		count = statement.GetInt32("count")
		if count < 1 {
			session.ExecuteStatement("create database root.run_status_device")
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	// 创建设备数据采集模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (value DOUBLE)", iotdb2.NOVA_DEVICE_TEMPLATE))
	// 创建设备运行时间统计模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (duration INT64, status INT64)", iotdb2.NOVA_DEVICE_RUN_TEMPLATE))

	return &Iot{}
}
