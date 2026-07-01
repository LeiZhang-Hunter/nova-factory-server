package controller

import (
	"nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/business/shop/home/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// HomeModule 首页模块控制器。
type HomeModule struct {
	service service.IShopHomeModuleService
}

// NewHomeModule 创建首页模块控制器。
func NewHomeModule(service service.IShopHomeModuleService) *HomeModule {
	return &HomeModule{service: service}
}

// PrivateRoutes 注册首页模块私有路由。
func (s *HomeModule) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/home/module")
	group.GET("/list", middlewares.HasPermission("shop:home:module:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:home:module:info"), s.Info)
	group.POST("/set", middlewares.HasPermission("shop:home:module:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:home:module:remove"), s.Delete)
}

func (s *HomeModule) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/shop/home/module/list", "shop:home:module:list")
	router.RegisterPermission("GET", "/shop/home/module/info/:id", "shop:home:module:info")
	router.RegisterPermission("POST", "/shop/home/module/set", "shop:home:module:set")
	router.RegisterPermission("DELETE", "/shop/home/module/remove/:ids", "shop:home:module:remove")
}

// List 获取首页模块列表。
// @Summary 获取首页模块列表
// @Description 分页查询首页模块列表
// @Tags 商城/首页装修/模块管理
// @Param object query models.HomeModuleQuery true "首页模块查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/home/module/list [get]
func (s *HomeModule) List(ctx *gin.Context) {
	req := new(models.HomeModuleQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := s.service.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Info 获取首页模块详情。
// @Summary 获取首页模块详情
// @Description 根据ID获取首页模块详情
// @Tags 商城/首页装修/模块管理
// @Param id path int true "首页模块ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/home/module/info/{id} [get]
func (s *HomeModule) Info(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := s.service.GetByID(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Set 保存首页模块。
// @Summary 保存首页模块
// @Description 新增或修改首页模块
// @Tags 商城/首页装修/模块管理
// @Param body body models.HomeModuleSet true "首页模块保存参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /shop/home/module/set [post]
func (s *HomeModule) Set(ctx *gin.Context) {
	req := new(models.HomeModuleSet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := s.service.Set(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Delete 删除首页模块。
// @Summary 删除首页模块
// @Description 根据ID删除首页模块
// @Tags 商城/首页装修/模块管理
// @Param ids path string true "首页模块ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/home/module/remove/{ids} [delete]
func (s *HomeModule) Delete(ctx *gin.Context) {
	ids := baizeContext.ParamInt64Array(ctx, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	if err := s.service.DeleteByIDs(ctx, ids); err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.Success(ctx)
}
