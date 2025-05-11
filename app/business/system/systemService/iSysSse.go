package systemService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemModels"
)

type ISseService interface {
	BuildNotificationChannel(c *gin.Context)
	SendNotification(c context.Context, ss *systemModels.Sse)
}
