package monitorService

import (
	"nova-factory-server/app/business/admin/monitor/monitorModels"

	"github.com/gin-gonic/gin"
)

type IUserOnlineService interface {
	SelectUserOnlineList(c *gin.Context, ol *monitorModels.SysUserOnlineDQL) (list []*monitorModels.SysUserOnline, total int64)
	ForceLogout(c *gin.Context, tokenId string)
}
