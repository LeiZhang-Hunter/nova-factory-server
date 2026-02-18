package controller

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/home/homeModels"
	"nova-factory-server/app/business/home/homeService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Home struct {
	iHomeService homeService.HomeService
}

func NewHome(iHomeService homeService.HomeService) *Home {
	return &Home{
		iHomeService: iHomeService,
	}
}

func (h *Home) PrivateRoutes(r *gin.RouterGroup) {
	group := r.Group("/home")
	{
		group.GET("/stats", middlewares.HasPermission("home:stats"), h.GetHomeStats)
		group.GET("/app/stats", middlewares.HasPermission("home:app:stats"), h.GetHomeStats)
	}
}

// GetHomeStats 首页统计
// @Summary 首页统计
// @Description 首页统计
// @Tags 首页
// @Security BearerAuth
// @Produce application/json
// @Router /home/stats [get]
func (h *Home) GetHomeStats(c *gin.Context) {
	req := new(homeModels.HomeRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("读取首页统计参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	isApp := false
	if req.Platform == "h5" {
		isApp = true
	}
	stats, err := h.iHomeService.GetHomeStats(c, isApp)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, stats)
}
