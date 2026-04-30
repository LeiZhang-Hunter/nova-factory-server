package controller

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
)

// SeckillConfig 商品秒杀配置控制器。
type SeckillConfig struct {
	service service.IShopSeckillConfigService
}

// NewSeckillConfig 创建商品秒杀配置控制器。
func NewSeckillConfig(service service.IShopSeckillConfigService) *SeckillConfig {
	return &SeckillConfig{service: service}
}

// PrivateRoutes 注册商品秒杀配置私有路由。
func (s *SeckillConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/activity/seckill_config")
	group.GET("/list", middlewares.HasPermission("shop:activity:seckill_config:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:activity:seckill_config:info"), s.Info)
	group.POST("/set", middlewares.HasPermission("shop:activity:seckill_config:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:activity:seckill_config:remove"), s.Delete)
}

// List 获取商品秒杀配置列表
// @Summary 获取商品秒杀配置列表
// @Description 分页查询商品秒杀配置列表
// @Tags 商城/活动管理/秒杀配置
// @Param object query models.SeckillConfigQuery true "商品秒杀配置查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/seckill_config/list [get]
func (s *SeckillConfig) List(ctx *gin.Context) {
	req := new(models.SeckillConfigQuery)
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

// Info 获取商品秒杀配置详情
// @Summary 获取商品秒杀配置详情
// @Description 根据ID获取商品秒杀配置详情
// @Tags 商城/活动管理/秒杀配置
// @Param id path int true "商品秒杀配置ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/seckill_config/info/{id} [get]
func (s *SeckillConfig) Info(ctx *gin.Context) {
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

// Set 保存商品秒杀配置
// @Summary 保存商品秒杀配置
// @Description 新增或修改商品秒杀配置
// @Tags 商城/活动管理/秒杀配置
// @Param body body models.SeckillConfigSet true "商品秒杀配置保存参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /shop/activity/seckill_config/set [post]
func (s *SeckillConfig) Set(ctx *gin.Context) {
	req := new(models.SeckillConfigSet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	// 开始时间与持续时间的总和不能超过 24 小时，避免配置跨天。
	if req.BeginClock+req.ContinueClock > 24 {
		baizeContext.Waring(ctx, "开启时间和持续时间总和不能超过24小时")
		return
	}
	var err error
	if req.Images != "" {
		req.Images, err = fileUtils.NormalizeResourcePath(req.Images)
		if err != nil {
			baizeContext.Waring(ctx, err.Error())
			zap.L().Error("NormalizeResourcePath error", zap.Any("err", err))
			return
		}
	}
	data, err := s.service.Set(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Delete 删除商品秒杀配置
// @Summary 删除商品秒杀配置
// @Description 根据ID删除商品秒杀配置
// @Tags 商城/活动管理/秒杀配置
// @Param ids path string true "商品秒杀配置ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/activity/seckill_config/remove/{ids} [delete]
func (s *SeckillConfig) Delete(ctx *gin.Context) {
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
