package homeservice

import (
	"nova-factory-server/app/business/iot/home/homemodels"

	"github.com/gin-gonic/gin"
)

type HomeService interface {
	GetHomeStats(c *gin.Context, isMobile bool) (*homemodels.HomeStats, error)
}
