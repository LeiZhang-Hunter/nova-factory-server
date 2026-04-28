package controller

import (
	"fmt"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Combination struct {
	service service.IShopCombinationService
}

func NewCombination(service service.IShopCombinationService) *Combination {
	return &Combination{service: service}
}

func (c *Combination) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/activity/combination")
	group.GET("/list", middlewares.HasPermission("shop:activity:combination:list"), c.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:activity:combination:info"), c.GetByID)
	group.POST("/set", middlewares.HasPermission("shop:activity:combination:set"), c.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:activity:combination:remove"), c.Delete)
}

// List 获取拼团商品列表
// @Summary 获取拼团商品列表
// @Description 分页查询拼团商品列表
// @Tags 商城/活动管理/拼团商品
// @Param object query models.CombinationQuery true "拼团商品查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/combination/list [get]
func (c *Combination) List(ctx *gin.Context) {
	req := new(models.CombinationQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.service.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetByID 获取拼团商品详情
// @Summary 获取拼团商品详情
// @Description 根据ID获取拼团商品详情
// @Tags 商城/活动管理/拼团商品
// @Param id path int true "拼团商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/activity/combination/info/{id} [get]
func (c *Combination) GetByID(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.service.GetByID(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	if data == nil {
		baizeContext.Waring(ctx, "拼团商品不存在")
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Set 保存拼团商品
// @Summary 保存拼团商品
// @Description 新增或修改拼团商品
// @Tags 商城/活动管理/拼团商品
// @Param body body models.CombinationSet true "拼团商品保存参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /shop/activity/combination/set [post]
func (c *Combination) Set(ctx *gin.Context) {
	req := new(models.CombinationSet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	if err := c.validateSet(ctx, req); err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	data, err := c.service.Set(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Delete 删除拼团商品
// @Summary 删除拼团商品
// @Description 根据ID删除拼团商品
// @Tags 商城/活动管理/拼团商品
// @Param ids path string true "拼团商品ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/activity/combination/remove/{ids} [delete]
func (c *Combination) Delete(ctx *gin.Context) {
	ids := baizeContext.ParamInt64Array(ctx, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	if err := c.service.DeleteByIDs(ctx, ids); err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.Success(ctx)
}

func (c *Combination) validateSet(ctx *gin.Context, req *models.CombinationSet) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Info = strings.TrimSpace(req.Info)
	req.UnitName = strings.TrimSpace(req.UnitName)
	req.Image = strings.TrimSpace(req.Image)
	req.Images = strings.TrimSpace(req.Images)
	req.Attr = strings.TrimSpace(req.Attr)

	if req.ProductID <= 0 {
		return fmt.Errorf("商品ID不能为空")
	}
	if req.Title == "" {
		return fmt.Errorf("活动标题不能为空")
	}
	if req.Image == "" {
		return fmt.Errorf("推荐图不能为空")
	}
	if req.People <= 0 {
		return fmt.Errorf("参团人数必须大于0")
	}
	if req.Price <= 0 {
		return fmt.Errorf("拼团价格必须大于0")
	}
	if req.Num <= 0 {
		return fmt.Errorf("单次购买数量必须大于0")
	}
	if req.OnceNum <= 0 {
		return fmt.Errorf("每个订单可购买数量必须大于0")
	}
	if req.Num < req.OnceNum {
		return fmt.Errorf("限制单次购买数量不能大于总购买数量")
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
	if req.EffectiveTime < 0 {
		return fmt.Errorf("拼团订单有效时间不能小于0")
	}
	if req.StartTime > 0 && req.StopTime > 0 {
		if req.StopTime <= req.StartTime {
			return fmt.Errorf("活动结束时间必须大于开始时间")
		}
		if req.StopTime < time.Now().Unix() {
			return fmt.Errorf("活动结束时间不能小于当前时间")
		}
	}
	if req.ID > 0 {
		current, err := c.service.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if current == nil {
			return fmt.Errorf("数据不存在")
		}
		if current.StopTime > 0 && current.StopTime < time.Now().Unix() {
			return fmt.Errorf("活动已结束,请重新添加")
		}
	}
	return nil
}
