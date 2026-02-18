package homeService

import (
	"nova-factory-server/app/business/home/homeModels"

	"github.com/gin-gonic/gin"
)

type HomeService interface {
	GetHomeStats(c *gin.Context, isMobile bool) (*homeModels.HomeStats, error)
}
