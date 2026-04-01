package systemservice

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type ISseService interface {
	BuildNotificationChannel(c *gin.Context)
	SendNotification(c context.Context, ss *systemmodels.Sse)
}
