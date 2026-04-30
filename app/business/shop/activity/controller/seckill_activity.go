package controller

import (
	"fmt"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strings"

	"github.com/gin-gonic/gin"
)

type SeckillActivity struct {
	service service.IShopSeckillActivityService
}

func NewSeckillActivity(service service.IShopSeckillActivityService) *SeckillActivity {
	return &SeckillActivity{service: service}
}

func (s *SeckillActivity) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/activity/seckill_activity")
	group.GET("/list", middlewares.HasPermission("shop:activity:seckill_activity:list"), s.List)
	group.POST("/set", middlewares.HasPermission("shop:activity:seckill_activity:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:activity:seckill_activity:remove"), s.Delete)
}

// List 获取秒杀活动列表
// @Summary 获取秒杀活动列表
// @Description 分页查询秒杀活动列表
// @Tags 商城/活动管理/秒杀活动
// @Param object query models.SeckillActivityQuery true "秒杀活动查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/seckill_activity/list [get]
func (s *SeckillActivity) List(ctx *gin.Context) {
	req := new(models.SeckillActivityQuery)
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

// Set 保存秒杀活动
// @Summary 保存秒杀活动
// @Description 新增或修改秒杀活动
// @Tags 商城/活动管理/秒杀活动
// @Param body body models.SeckillActivitySet true "秒杀活动保存参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /shop/activity/seckill_activity/set [post]
func (s *SeckillActivity) Set(ctx *gin.Context) {
	req := new(models.SeckillActivitySet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	if err := s.validateSet(ctx, req); err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	data, err := s.service.Set(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Delete 删除秒杀活动
// @Summary 删除秒杀活动
// @Description 根据ID删除秒杀活动
// @Tags 商城/活动管理/秒杀活动
// @Param ids path string true "秒杀活动ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/activity/seckill_activity/remove/{ids} [delete]
func (s *SeckillActivity) Delete(ctx *gin.Context) {
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

func (s *SeckillActivity) validateSet(ctx *gin.Context, req *models.SeckillActivitySet) error {
	req.Title = strings.TrimSpace(req.Title)
	req.TimeIDs = normalizeSeckillActivityTimeIDs(req.TimeIDs)

	if req.Type == 0 {
		req.Type = 1
	}
	if req.Type < 0 {
		return fmt.Errorf("活动类型不正确")
	}
	if req.Title == "" {
		return fmt.Errorf("活动名称不能为空")
	}
	if req.StartDay <= 0 {
		return fmt.Errorf("开始日期不能为空")
	}
	if req.EndDay <= 0 {
		return fmt.Errorf("结束日期不能为空")
	}
	if req.EndDay < req.StartDay {
		return fmt.Errorf("结束日期不能小于开始日期")
	}
	if req.TimeIDs == "" {
		return fmt.Errorf("时间段ID不能为空")
	}
	if req.OnceNum < 0 {
		return fmt.Errorf("每日购买数量不能小于0")
	}
	if req.Num < 0 {
		return fmt.Errorf("活动总购买数不能小于0")
	}
	if req.Num > 0 && req.OnceNum > 0 && req.OnceNum > req.Num {
		return fmt.Errorf("每日购买数量不能大于活动总购买数")
	}
	if req.IsCommission < 0 {
		return fmt.Errorf("分佣状态不正确")
	}
	if req.Status < 0 {
		return fmt.Errorf("显示状态不正确")
	}
	if req.LinkID < 0 {
		return fmt.Errorf("关联ID不能小于0")
	}

	if req.ID > 0 {
		current, err := s.service.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if current == nil {
			return fmt.Errorf("数据不存在")
		}
	}
	return nil
}

func normalizeSeckillActivityTimeIDs(raw string) string {
	parts := strings.Split(strings.TrimSpace(raw), ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	return strings.Join(result, ",")
}
