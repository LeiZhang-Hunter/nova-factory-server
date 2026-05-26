package activity

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type SeckillController struct {
	Service service.IApiShopSeckillService
}

func NewSeckillController(service service.IApiShopSeckillService) *SeckillController {
	return &SeckillController{Service: service}
}

// PrivateRoutes 注册秒杀路由
func (c *SeckillController) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/seckill")
	group.GET("/config/list", c.ListConfigs)
	group.GET("/goods/list", c.ListGoods)
	group.GET("/goods/:id", c.GetGoodsDetail)
}

// ListConfigs 获取秒杀时间段列表
func (c *SeckillController) ListConfigs(ctx *gin.Context) {
	data, err := c.Service.ListConfigs(ctx)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// ListGoods 获取秒杀商品列表
func (c *SeckillController) ListGoods(ctx *gin.Context) {
	req := &models.SeckillQuery{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.Service.ListGoods(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetGoodsDetail 获取秒杀商品详情
func (c *SeckillController) GetGoodsDetail(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	timeID := baizeContext.QueryInt64(ctx, "timeId")
	if timeID == 0 {
		timeID = baizeContext.QueryInt64(ctx, "time_id")
	}
	data, err := c.Service.GetGoodsDetail(ctx, id, timeID)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}
