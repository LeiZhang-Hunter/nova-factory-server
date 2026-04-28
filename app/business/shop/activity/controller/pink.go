package controller

import (
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Pink struct {
	service service.IShopPinkService
}

func NewPink(service service.IShopPinkService) *Pink {
	return &Pink{service: service}
}

func (p *Pink) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/activity/pink")
	group.GET("/list", middlewares.HasPermission("shop:activity:pink:list"), p.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:activity:pink:info"), p.GetByID)
}

// List 获取拼团列表
// @Summary 获取拼团列表
// @Description 分页查询拼团开团记录
// @Tags 商城/活动管理/拼团记录
// @Param object query models.PinkQuery true "拼团记录查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/pink/list [get]
func (p *Pink) List(ctx *gin.Context) {
	req := new(models.PinkQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := p.service.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetByID 获取拼团详情
// @Summary 获取拼团详情
// @Description 根据ID获取拼团开团记录详情
// @Tags 商城/活动管理/拼团记录
// @Param id path int true "拼团记录ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/pink/info/{id} [get]
func (p *Pink) GetByID(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := p.service.GetByID(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	if data == nil {
		baizeContext.Waring(ctx, "拼团记录不存在")
		return
	}
	baizeContext.SuccessData(ctx, data)
}
