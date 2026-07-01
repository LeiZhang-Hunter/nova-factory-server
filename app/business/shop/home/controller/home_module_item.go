package controller

import (
	"nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/business/shop/home/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// HomeModuleItem 首页模块明细控制器。
type HomeModuleItem struct {
	service service.IShopHomeModuleItemService
}

// NewHomeModuleItem 创建首页模块明细控制器。
func NewHomeModuleItem(service service.IShopHomeModuleItemService) *HomeModuleItem {
	return &HomeModuleItem{service: service}
}

// PrivateRoutes 注册首页模块明细私有路由。
func (s *HomeModuleItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/home/module_item")
	group.GET("/list", middlewares.HasPermission("shop:home:module_item:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:home:module_item:info"), s.Info)
	group.POST("/set", middlewares.HasPermission("shop:home:module_item:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:home:module_item:remove"), s.Delete)
}

func (s *HomeModuleItem) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/shop/home/module_item/list", "shop:home:module_item:list")
	router.RegisterPermission("GET", "/shop/home/module_item/info/:id", "shop:home:module_item:info")
	router.RegisterPermission("POST", "/shop/home/module_item/set", "shop:home:module_item:set")
	router.RegisterPermission("DELETE", "/shop/home/module_item/remove/:ids", "shop:home:module_item:remove")
}

// List 获取首页模块明细列表。
// @Summary 获取首页模块明细列表
// @Description 分页查询首页模块明细列表
// @Tags 商城/首页装修/模块明细管理
// @Param object query models.HomeModuleItemQuery true "首页模块明细查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/home/module_item/list [get]
func (s *HomeModuleItem) List(ctx *gin.Context) {
	req := new(models.HomeModuleItemQuery)
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

// Info 获取首页模块明细详情。
// @Summary 获取首页模块明细详情
// @Description 根据ID获取首页模块明细详情
// @Tags 商城/首页装修/模块明细管理
// @Param id path int true "首页模块明细ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/home/module_item/info/{id} [get]
func (s *HomeModuleItem) Info(ctx *gin.Context) {
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

// Set 保存首页模块明细。
// @Summary 保存首页模块明细
// @Description 新增或修改首页模块明细
// @Tags 商城/首页装修/模块明细管理
// @Param body body models.HomeModuleItemSet true "首页模块明细保存参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /shop/home/module_item/set [post]
func (s *HomeModuleItem) Set(ctx *gin.Context) {
	req := new(models.HomeModuleItemSet)
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

// Delete 删除首页模块明细。
// @Summary 删除首页模块明细
// @Description 根据ID删除首页模块明细
// @Tags 商城/首页装修/模块明细管理
// @Param ids path string true "首页模块明细ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/home/module_item/remove/{ids} [delete]
func (s *HomeModuleItem) Delete(ctx *gin.Context) {
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
