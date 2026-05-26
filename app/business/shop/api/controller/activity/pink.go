package activity

import (
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type PinkController struct {
	Service service.IApiShopPinkService
}

func NewPinkController(service service.IApiShopPinkService) *PinkController {
	return &PinkController{Service: service}
}

// PrivateRoutes 注册拼团路由
func (c *PinkController) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/pink")
	group.POST("", c.Create)             // 开团或参团
	group.GET("/:id", c.GetDetail)       // 获取拼团详情
	group.GET("/list/my", c.ListMyPinks) // 获取我的拼团记录
}

// Create 开团或参团
// @Summary 开团或参团
// @Description 创建拼团订单时调用，开团 pink_id=0，参团 pink_id>0
// @Tags 商城/拼团
// @Param body body PinkCreateReq true "拼团参数"
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "操作成功"
// @Router /api/v1/app/shop/pink [post]
func (c *PinkController) Create(ctx *gin.Context) {
	req := new(PinkCreateReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	// 获取当前用户ID（从鉴权中间件获取）
	userID := baizeContext.GetUserId(ctx)
	if userID == 0 {
		baizeContext.InvalidToken(ctx)
		return
	}

	data, err := c.Service.Create(ctx, userID, req.CombinationID, req.PinkID, req.OrderID)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// GetDetail 获取拼团详情
func (c *PinkController) GetDetail(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}

	data, err := c.Service.GetDetail(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// ListMyPinks 获取我的拼团记录
func (c *PinkController) ListMyPinks(ctx *gin.Context) {
	userID := baizeContext.GetUserId(ctx)
	if userID == 0 {
		baizeContext.InvalidToken(ctx)
		return
	}

	data, err := c.Service.ListMyPinks(ctx, userID)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// PinkCreateReq 拼团创建请求
type PinkCreateReq struct {
	CombinationID int64 `json:"combinationId,string" binding:"required"` // 拼团商品ID
	PinkID        int64 `json:"pinkId,string"`                           // 拼团记录ID（0=开团，>0=参团）
	OrderID       int64 `json:"orderId,string"`                          // 订单ID
}
