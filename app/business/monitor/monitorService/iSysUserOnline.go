package monitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorModels"
)

type IUserOnlineService interface {
	SelectUserOnlineList(c *gin.Context, ol *monitorModels.SysUserOnlineDQL) (list []*monitorModels.SysUserOnline, total int64)
	ForceLogout(c *gin.Context, tokenId string)
}
