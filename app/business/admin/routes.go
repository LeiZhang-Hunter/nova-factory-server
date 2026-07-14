package admin

import (
	"context"

	basicscontroller "nova-factory-server/app/business/admin/basics/controller"
	"nova-factory-server/app/business/admin/monitor/monitorcontroller"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/business/admin/system/systemcontroller"
	"nova-factory-server/app/business/admin/tool/toolcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/objectFile/localhostObject"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

type Admin struct{}

func NewGinEngine(app *routes.App,
	cache cache.Cache,
	bc *basicscontroller.Basics,
	sc *systemcontroller.System,
	mc *monitorcontroller.Monitor,
	gc *toolcontroller.Tool,
	requestLogSvc monitorservice.IRequestLogService) *Admin {
	r := app.Engine
	group := r.Group("")
	{
		if viper.GetString("upload_file.type") == "local" {
			group.Static(localhostObject.ResourcePrefix, viper.GetString("upload_file.localhost.public_path"))
		}
		sc.Login.PublicRoutes(group)
		sc.Sse.PublicRoutes(group)
	}

	// 启动请求日志定时清理
	go requestLogSvc.ScheduleClean(context.Background())

	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		sc.Profile.PrivateRoutes(group)
		sc.Login.PrivateRoutes(group)
		sc.User.PrivateRoutes(group)
		sc.Dept.PrivateRoutes(group)
		sc.DictType.PrivateRoutes(group)
		sc.DictData.PrivateRoutes(group)
		sc.Role.PrivateRoutes(group)
		sc.Post.PrivateRoutes(group)
		sc.Permission.PrivateRoutes(group)
		sc.Config.PrivateRoutes(group)
		sc.File.PrivateRoutes(group)
		sc.Notice.PrivateRoutes(group)
		sc.SelectBox.PrivateRoutes(group)
		sc.Shift.PrivateRoutes(group)
		bc.CompanyInfo.PrivateRoutes(group)
		mc.Server.PrivateRoutes(group)
		mc.Oper.PrivateRoutes(group)
		mc.RequestLog.PrivateRoutes(group)
		mc.UserOnline.PrivateRoutes(group)
		mc.Logfor.PrivateRoutes(group)
		mc.Job.PrivateRoutes(group)
		gc.GenTable.PrivateRoutes(group)
	}

	{
		sc.User.PrivateMcpRoutes(app.McpServer)
		sc.Dept.PrivateMcpRoutes(app.McpServer)
		sc.DictType.PrivateMcpRoutes(app.McpServer)
		sc.DictData.PrivateMcpRoutes(app.McpServer)
		sc.Role.PrivateMcpRoutes(app.McpServer)
		sc.Post.PrivateMcpRoutes(app.McpServer)
		sc.Permission.PrivateMcpRoutes(app.McpServer)
		sc.Config.PrivateMcpRoutes(app.McpServer)
		sc.Notice.PrivateMcpRoutes(app.McpServer)
		sc.SelectBox.PrivateMcpRoutes(app.McpServer)
		sc.Shift.PrivateMcpRoutes(app.McpServer)
		bc.CompanyInfo.PrivateMcpRoutes(app.McpServer)
		mc.Server.PrivateMcpRoutes(app.McpServer)
		mc.Oper.PrivateMcpRoutes(app.McpServer)
		mc.RequestLog.PrivateMcpRoutes(app.McpServer)
		mc.UserOnline.PrivateMcpRoutes(app.McpServer)
		mc.Logfor.PrivateMcpRoutes(app.McpServer)
		mc.Job.PrivateMcpRoutes(app.McpServer)
		gc.GenTable.PrivateMcpRoutes(app.McpServer)
	}

	return &Admin{}
}
