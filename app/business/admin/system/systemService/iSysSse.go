package systemService

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemModels"

	"github.com/gin-gonic/gin"
)

type ISseService interface {
	BuildNotificationChannel(c *gin.Context)
	SendNotification(c context.Context, ss *systemModels.Sse)
}
