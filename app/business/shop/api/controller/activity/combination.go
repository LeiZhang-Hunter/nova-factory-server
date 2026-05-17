package activity

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type CombinationController struct {
	CombinationService service.IApiShopCombinationService
	PinkService        service.IApiShopPinkService
}

func NewCombinationController(
	combinationService service.IApiShopCombinationService,
	pinkService service.IApiShopPinkService,
) *CombinationController {
	return &CombinationController{
		CombinationService: combinationService,
		PinkService:        pinkService,
	}
}

// PrivateRoutes 注册拼团路由
func (c *CombinationController) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/combination")
	group.GET("/list", c.List)
	group.GET("/:id", c.GetByID)
	group.GET("/:id/pink", c.GetPinkList)
	group.GET("/pink/:pinkId", c.GetPinkDetail)
}

// List 获取拼团商品列表
func (c *CombinationController) List(ctx *gin.Context) {
	req := new(models.CombinationQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.CombinationService.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetByID 获取拼团商品详情
func (c *CombinationController) GetByID(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.CombinationService.GetByID(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetPinkList 获取正在进行中的团列表
func (c *CombinationController) GetPinkList(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.CombinationService.GetPinkList(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetPinkDetail 获取团详情（含当前人数）
func (c *CombinationController) GetPinkDetail(ctx *gin.Context) {
	pinkID := baizeContext.ParamInt64(ctx, "pinkId")
	if pinkID == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	pink, err := c.PinkService.GetDetail(ctx, pinkID)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	if pink == nil {
		baizeContext.Waring(ctx, "拼团记录不存在")
		return
	}
	// 获取当前团内人数
	count, _ := c.PinkService.GetPinkMemberCount(ctx, pinkID)
	type PinkDetail struct {
		*models.Pink
		CurrentPeople int64 `json:"currentPeople"`
	}
	baizeContext.SuccessData(ctx, PinkDetail{
		Pink:          pink,
		CurrentPeople: count,
	})
}
