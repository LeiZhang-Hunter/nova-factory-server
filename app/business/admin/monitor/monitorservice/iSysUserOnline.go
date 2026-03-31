package monitorservice

import (
	"nova-factory-server/app/business/admin/monitor/monitormodels"

	"github.com/gin-gonic/gin"
)

type IUserOnlineService interface {
	SelectUserOnlineList(c *gin.Context, ol *monitormodels.SysUserOnlineDQL) (list []*monitormodels.SysUserOnline, total int64)
	ForceLogout(c *gin.Context, tokenId string)
}
