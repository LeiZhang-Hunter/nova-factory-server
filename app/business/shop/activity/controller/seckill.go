package controller

import (
	"fmt"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/fileUtils"
	"nova-factory-server/app/utils/gin_mcp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Seckill struct {
	service service.IShopSeckillService
}

func NewSeckill(service service.IShopSeckillService) *Seckill {
	return &Seckill{service: service}
}

func (s *Seckill) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/activity/seckill")
	group.GET("/list", middlewares.HasPermission("shop:activity:seckill:list"), s.List)
	group.POST("/set", middlewares.HasPermission("shop:activity:seckill:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:activity:seckill:remove"), s.Delete)
}

func (s *Seckill) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/shop/activity/seckill/list", "shop:activity:seckill:list")
	router.RegisterPermission("POST", "/shop/activity/seckill/set", "shop:activity:seckill:set")
	router.RegisterPermission("DELETE", "/shop/activity/seckill/remove/:ids", "shop:activity:seckill:remove")
}

// List 获取秒杀商品列表
// @Summary 获取秒杀商品列表
// @Description 分页查询秒杀商品列表
// @Tags 商城/活动管理/秒杀商品
// @Param object query models.SeckillQuery true "秒杀商品查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/seckill/list [get]
func (s *Seckill) List(ctx *gin.Context) {
	req := new(models.SeckillQuery)
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

// Set 保存秒杀商品
// @Summary 保存秒杀商品
// @Description 新增或修改秒杀商品
// @Tags 商城/活动管理/秒杀商品
// @Param body body models.SeckillSet true "秒杀商品保存参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /shop/activity/seckill/set [post]
func (s *Seckill) Set(ctx *gin.Context) {
	req := new(models.SeckillSet)
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

// Delete 删除秒杀商品
// @Summary 删除秒杀商品
// @Description 根据ID删除秒杀商品
// @Tags 商城/活动管理/秒杀商品
// @Param ids path string true "秒杀商品ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/activity/seckill/remove/{ids} [delete]
func (s *Seckill) Delete(ctx *gin.Context) {
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

func (s *Seckill) validateSet(ctx *gin.Context, req *models.SeckillSet) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Info = strings.TrimSpace(req.Info)
	req.UnitName = strings.TrimSpace(req.UnitName)
	req.Image = strings.TrimSpace(req.Image)
	req.Images = strings.TrimSpace(req.Images)
	req.StartTime = strings.TrimSpace(req.StartTime)
	req.StopTime = strings.TrimSpace(req.StopTime)
	req.TimeID = strings.TrimSpace(req.TimeID)
	req.Logistics = strings.TrimSpace(req.Logistics)
	req.CustomForm = strings.TrimSpace(req.CustomForm)

	if req.ProductID <= 0 {
		return fmt.Errorf("商品ID不能为空")
	}
	if req.Title == "" {
		return fmt.Errorf("活动标题不能为空")
	}
	if req.Image == "" {
		return fmt.Errorf("推荐图不能为空")
	}
	if req.Price <= 0 {
		return fmt.Errorf("秒杀价格必须大于0")
	}
	if req.OtPrice < 0 {
		return fmt.Errorf("原价不能小于0")
	}
	if req.Cost < 0 {
		return fmt.Errorf("成本不能小于0")
	}
	if req.GiveIntegral < 0 {
		return fmt.Errorf("返积分不能小于0")
	}
	if req.Stock < 0 {
		return fmt.Errorf("库存不能小于0")
	}
	if req.Sort < 0 {
		return fmt.Errorf("排序值不能小于0")
	}
	if req.Postage < 0 {
		return fmt.Errorf("邮费不能小于0")
	}
	if req.Num <= 0 {
		return fmt.Errorf("最多秒杀数量必须大于0")
	}
	if req.OnceNum < 0 {
		return fmt.Errorf("单次购买数量不能小于0")
	}
	if req.OnceNum > 0 && int64(req.OnceNum) > int64(req.Num) {
		return fmt.Errorf("单次购买数量不能大于最多秒杀数量")
	}
	if req.Weight < 0 {
		return fmt.Errorf("商品重量不能小于0")
	}
	if req.Volume < 0 {
		return fmt.Errorf("商品体积不能小于0")
	}
	if req.Quota < 0 {
		return fmt.Errorf("限购总数不能小于0")
	}
	if req.QuotaShow < 0 {
		return fmt.Errorf("限购总数显示不能小于0")
	}
	if req.StartTime == "" {
		return fmt.Errorf("开始时间不能为空")
	}
	if req.StopTime == "" {
		return fmt.Errorf("结束时间不能为空")
	}

	startAt, err := parseSeckillTime(req.StartTime)
	if err != nil {
		return fmt.Errorf("开始时间格式不正确")
	}
	stopAt, err := parseSeckillTime(req.StopTime)
	if err != nil {
		return fmt.Errorf("结束时间格式不正确")
	}
	if !stopAt.After(startAt) {
		return fmt.Errorf("活动结束时间必须大于开始时间")
	}
	if stopAt.Before(time.Now()) {
		return fmt.Errorf("活动结束时间不能小于当前时间")
	}

	if req.ID > 0 {
		current, err := s.service.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if current == nil {
			return fmt.Errorf("数据不存在")
		}
		currentStopAt, err := parseSeckillTime(current.StopTime)
		if err == nil && currentStopAt.Before(time.Now()) {
			return fmt.Errorf("活动已结束,请重新添加")
		}
	}

	if req.Image != "" {
		req.Image, err = fileUtils.NormalizeResourcePath(req.Image)
		if err != nil {
			return err
		}
	}
	if req.Images != "" {
		imagesArr := strings.Split(req.Images, ",")
		images := make([]string, 0, len(imagesArr))
		for _, img := range imagesArr {
			img = strings.TrimSpace(img)
			if img == "" {
				continue
			}
			imgPath, normalizeErr := fileUtils.NormalizeResourcePath(img)
			if normalizeErr != nil {
				return normalizeErr
			}
			images = append(images, imgPath)
		}
		req.Images = strings.Join(images, ",")
	}
	return nil
}

func parseSeckillTime(raw string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		time.RFC3339,
	}
	raw = strings.TrimSpace(raw)
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid seckill time")
}
