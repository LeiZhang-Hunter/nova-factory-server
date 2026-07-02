package controller

import (
	"fmt"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SeckillActivity struct {
	service       service.IShopSeckillActivityService
	configService service.IShopSeckillConfigService
}

func NewSeckillActivity(service service.IShopSeckillActivityService, configService service.IShopSeckillConfigService) *SeckillActivity {
	return &SeckillActivity{
		service:       service,
		configService: configService,
	}
}

func (s *SeckillActivity) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/activity/seckill_activity")
	group.GET("/list", middlewares.HasPermission("shop:activity:seckill_activity:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:activity:seckill_activity:info"), s.Info)
	group.POST("/set", middlewares.HasPermission("shop:activity:seckill_activity:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:activity:seckill_activity:remove"), s.Delete)
}

func (s *SeckillActivity) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/shop/activity/seckill_activity/list", "shop:activity:seckill_activity:list")
	router.RegisterPermission("GET", "/shop/activity/seckill_activity/info/:id", "shop:activity:seckill_activity:info")
	router.RegisterPermission("POST", "/shop/activity/seckill_activity/set", "shop:activity:seckill_activity:set")
	router.RegisterPermission("DELETE", "/shop/activity/seckill_activity/remove/:ids", "shop:activity:seckill_activity:remove")
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

// Info 获取秒杀活动详情
// @Summary 获取秒杀活动详情
// @Description 根据ID获取秒杀活动详情
// @Tags 商城/活动管理/秒杀活动
// @Param id path int true "秒杀活动ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/seckill_activity/info/{id} [get]
func (s *SeckillActivity) Info(ctx *gin.Context) {
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
	if len(req.TimeIDs) == 0 {
		return fmt.Errorf("时间段ID不能为空")
	}
	configIDs, invalidTimeID := parseSeckillConfigIDs(req.TimeIDs)
	if invalidTimeID != "" {
		return fmt.Errorf("时间段ID格式不正确: %s", invalidTimeID)
	}
	configs, err := s.configService.GetByIDs(ctx, configIDs)
	if err != nil {
		return err
	}
	exists := make(map[int64]struct{}, len(configs))
	for _, config := range configs {
		if config == nil {
			continue
		}
		exists[config.ID] = struct{}{}
	}
	for _, configID := range configIDs {
		if _, ok := exists[configID]; !ok {
			return fmt.Errorf("时间段ID不存在: %d", configID)
		}
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
	if len(req.ProductInfos) == 0 {
		return fmt.Errorf("请选择参与活动的商品")
	}
	if err := validateSeckillActivityProducts(req.ProductInfos); err != nil {
		return err
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

// normalizeSeckillActivityTimeIDs 清洗时间段ID数组，移除空值和多余空格。
func normalizeSeckillActivityTimeIDs(raw []string) []string {
	result := make([]string, 0, len(raw))
	for _, part := range raw {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	return result
}

// parseSeckillConfigIDs 将时间段ID数组转换为去重后的主键数组。
func parseSeckillConfigIDs(raw []string) ([]int64, string) {
	result := make([]int64, 0, len(raw))
	seen := make(map[int64]struct{}, len(raw))
	for _, item := range raw {
		configID, err := strconv.ParseInt(item, 10, 64)
		if err != nil || configID <= 0 {
			return nil, item
		}
		if _, ok := seen[configID]; ok {
			continue
		}
		seen[configID] = struct{}{}
		result = append(result, configID)
	}
	return result, ""
}

// validateSeckillActivityProducts 校验活动商品列表。
func validateSeckillActivityProducts(productInfos []*models.SeckillActivityProductInfo) error {
	seen := make(map[int64]struct{}, len(productInfos))
	for _, productInfo := range productInfos {
		if productInfo == nil || productInfo.ID <= 0 {
			return fmt.Errorf("活动商品不能为空")
		}
		if _, ok := seen[productInfo.ID]; ok {
			return fmt.Errorf("活动商品重复: %d", productInfo.ID)
		}
		seen[productInfo.ID] = struct{}{}
		if productInfo.Status < 0 {
			return fmt.Errorf("商品状态不正确: %d", productInfo.ID)
		}
		if productInfo.Sort < 0 {
			return fmt.Errorf("商品排序不能小于0: %d", productInfo.ID)
		}
		if productInfo.IsHot < 0 {
			return fmt.Errorf("商品热门状态不正确: %d", productInfo.ID)
		}
		for _, attr := range productInfo.Attrs {
			if attr == nil {
				continue
			}
			if (attr.SkuID) == 0 {
				return fmt.Errorf("商品规格不能为空: %d", productInfo.ID)
			}
			if attr.Status < 0 {
				return fmt.Errorf("商品规格状态不正确: %s", attr.SkuID)
			}
			if attr.Price < 0 {
				return fmt.Errorf("商品规格活动价不能小于0: %s", attr.SkuID)
			}
			if attr.Cost < 0 {
				return fmt.Errorf("商品规格成本价不能小于0: %s", attr.SkuID)
			}
			if attr.OtPrice < 0 {
				return fmt.Errorf("商品规格原价不能小于0: %s", attr.SkuID)
			}
			if attr.Quota < 0 {
				return fmt.Errorf("商品规格限量不能小于0: %s", attr.SkuID)
			}
		}
	}
	return nil
}
