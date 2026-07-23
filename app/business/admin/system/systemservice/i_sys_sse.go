package systemservice

import (
	"context"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"

	"github.com/gin-gonic/gin"
)

type ISseService interface {
	BuildNotificationChannel(c *gin.Context)
	SendNotification(c context.Context, ss *modelentity.Sse)
}
