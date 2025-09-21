package routes

import (
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/business/alert/alertController"
	"nova-factory-server/app/business/asset/building/buildingController"
	"nova-factory-server/app/business/asset/device/deviceController"
	"nova-factory-server/app/business/asset/material/materialController"
	"nova-factory-server/app/business/craft/craftRouteController"
	"nova-factory-server/app/business/daemonize/daemonizeController"
	"nova-factory-server/app/business/dashboard/dashboardController"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorController"
	"nova-factory-server/app/business/metric/device/metricController"
	"nova-factory-server/app/business/monitor/monitorController"
	"nova-factory-server/app/business/system/systemController"
	"nova-factory-server/app/business/tool/toolController"
	"nova-factory-server/app/daemonize"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/objectFile/localhostObject"
	"nova-factory-server/app/utils/logger"
	"time"

	"nova-factory-server/app/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/google/wire"
	"net/http"
	"nova-factory-server/app/docs"
	"nova-factory-server/app/setting"
	"strings"

	"github.com/gin-gonic/gin"
)

var ProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	cache cache.Cache,
	sc *systemController.System,
	mc *monitorController.Monitor,
	gc *toolController.Tool,
	dc *deviceController.Device,
	materialC *materialController.Material,
	ai *aiDataSetController.AiDataSet,
	craft *craftRouteController.CraftRoute,
	metric *metricController.MetricServer,
	controller *daemonizeController.DaemonizeServer,
	deviceMonitor *deviceMonitorController.DeviceMonitorController,
	alert *alertController.Controller,
	building buildingController.Controller,
	dashboard dashboardController.Controller,
) *gin.Engine {

	if setting.Conf.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.NewLoggerMiddlewareBuilder().
		IgnorePaths("/ping").Build())
	r.Use(newCors())
	group := r.Group("")
	if setting.Conf.Mode == "dev" {
		host := setting.Conf.Host
		docs.SwaggerInfo.Host = host[strings.Index(host, "//")+2:]
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		group.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}

	//不做鉴权的
	{

		if viper.GetString("upload_file.type") == "local" {
			group.Static(localhostObject.ResourcePrefix, viper.GetString("upload_file.localhost.public_path"))
		}
		sc.Login.PublicRoutes(group)          //登录
		sc.Sse.PublicRoutes(group)            //SSE链接
		controller.Config.PublicRoutes(group) //注册Agent公共配置接口
		alert.AlertLog.PublicRoutes(group)
		metric.Metric.PublicRoutes(group)
	}
	//做鉴权的
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{

		sc.Profile.PrivateRoutes(group)    //个人资料
		sc.Login.PrivateRoutes(group)      //登录
		sc.User.PrivateRoutes(group)       //用户
		sc.Dept.PrivateRoutes(group)       //部门
		sc.DictType.PrivateRoutes(group)   //字典类型
		sc.DictData.PrivateRoutes(group)   //地点数据
		sc.Role.PrivateRoutes(group)       //角色
		sc.Post.PrivateRoutes(group)       //岗位
		sc.Permission.PrivateRoutes(group) //岗位
		sc.Config.PrivateRoutes(group)     //配置
		sc.File.PrivateRoutes(group)       //文件
		sc.Notice.PrivateRoutes(group)     //消息
		sc.SelectBox.PrivateRoutes(group)
		mc.Server.PrivateRoutes(group)     //服务器详情
		mc.Oper.PrivateRoutes(group)       //操作日志
		mc.UserOnline.PrivateRoutes(group) //在线用户
		mc.Logfor.PrivateRoutes(group)     //登录日志
		mc.Job.PrivateRoutes(group)        //定时任务
		gc.GenTable.PrivateRoutes(group)   //代码生成

		dc.Info.PrivateRoutes(group)  //资产管理---设备模块
		dc.Group.PrivateRoutes(group) //资产管理---设备分组
		dc.Template.PrivateRoutes(group)
		dc.TemplateData.PrivateRoutes(group)

		materialC.Material.PrivateRoutes(group) //资产管理---物料管理

		ai.Dataset.PrivateRoutes(group)    // 工业智能体
		ai.Prediction.PrivateRoutes(group) // 工业智能体
		ai.Exception.PrivateRoutes(group)  // 工业智能体
		ai.Control.PrivateRoutes(group)

		craft.CraftRoute.PrivateRoutes(group)     //工艺路线
		craft.Process.PrivateRoutes(group)        //工序设置
		craft.ProcessContext.PrivateRoutes(group) //工序内容
		craft.RouteProcess.PrivateRoutes(group)   //工艺组成
		craft.RouteProduct.PrivateRoutes(group)
		craft.RouteProductBom.PrivateRoutes(group)
		craft.WorkOrder.PrivateRoutes(group) //生产工单
		craft.Task.PrivateRoutes(group)      //工单任务
		craft.Schedule.PrivateRoutes(group)
		craft.Schedule.PublicRoutes(group)

		metric.Metric.PrivateRoutes(group) //设备指标
		controller.IotAgent.PrivateRoutes(group)
		controller.Config.PrivateRoutes(group)

		deviceMonitor.DeviceMonitor.PrivateRoutes(group)
		deviceMonitor.DeviceReport.PrivateRoutes(group)

		alert.AlertTemplate.PrivateRoutes(group)
		alert.Alert.PrivateRoutes(group)
		alert.AlertLog.PrivateRoutes(group)
		alert.AlertAction.PrivateRoutes(group)
		alert.AlertAiReason.PrivateRoutes(group)

		building.Building.PrivateRoutes(group)

		dashboard.Dashboard.PrivateRoutes(group)
		dashboard.Data.PrivateRoutes(group)

		alert.Runner.Run()
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	s := daemonize.CreateGRpcServer()
	controller.Daemonize.PrivateRoutes(s)
	s.Start()
	return r

}

func newCors() gin.HandlerFunc {
	ss := []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma", "FooBar"}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} //允许访问的域名
	config.AllowMethods = []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = ss
	config.MaxAge = time.Hour
	config.AllowCredentials = false
	return cors.New(config)
}
