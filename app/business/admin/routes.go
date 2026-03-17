package admin

import (
	"nova-factory-server/app/business/admin/monitor/monitorController"
	"nova-factory-server/app/business/admin/product/productController"
	"nova-factory-server/app/business/admin/system/systemController"
	"nova-factory-server/app/business/admin/tool/toolController"
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
	sc *systemController.System,
	mc *monitorController.Monitor,
	gc *toolController.Tool,
	product *productController.Product) *Admin {
	r := app.Engine
	group := r.Group("")
	//不做鉴权的
	{
		if viper.GetString("upload_file.type") == "local" {
			group.Static(localhostObject.ResourcePrefix, viper.GetString("upload_file.localhost.public_path"))
		}
		sc.Login.PublicRoutes(group) //登录
		sc.Sse.PublicRoutes(group)   //SSE链接

		product.Laboratory.PublicRoutes(group)
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
		sc.Shift.PrivateRoutes(group)
		mc.Server.PrivateRoutes(group)     //服务器详情
		mc.Oper.PrivateRoutes(group)       //操作日志
		mc.UserOnline.PrivateRoutes(group) //在线用户
		mc.Logfor.PrivateRoutes(group)     //登录日志
		mc.Job.PrivateRoutes(group)        //定时任务
		gc.GenTable.PrivateRoutes(group)   //代码生成
		product.Laboratory.PrivateRoutes(group)
	}

	return &Admin{}
}
