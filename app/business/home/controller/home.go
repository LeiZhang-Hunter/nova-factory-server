package controller

import (
	"nova-factory-server/app/business/home/homeService"
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
		group.GET("/stats", h.GetHomeStats)
	}
}

func (h *Home) GetHomeStats(c *gin.Context) {
	stats, err := h.iHomeService.GetHomeStats(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, stats)
}
