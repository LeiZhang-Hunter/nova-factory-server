package monitorServiceImpl

import (
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserOnlineService struct {
	cache cache.Cache
}

func NewUserOnlineService(cache cache.Cache) monitorservice.IUserOnlineService {
	return &UserOnlineService{cache: cache}
}

func (userOnlineService *UserOnlineService) SelectUserOnlineList(c *gin.Context, ol *monitormodels.SysUserOnlineDQL) (list []*monitormodels.SysUserOnline, total int64) {
	manager := session.NewAdminManager(userOnlineService.cache)
	sessions := manager.ScanSessions(c)

	list = make([]*monitormodels.SysUserOnline, 0, len(sessions))
	for _, newSession := range sessions {
		oui := new(monitormodels.SysUserOnline)
		oui.TokenId = newSession.Id()
		oui.UserName = newSession.Get(c, sessionStatus.UserName)
		oui.Browser = newSession.Get(c, sessionStatus.Browser)
		oui.Ipaddr = newSession.Get(c, sessionStatus.IpAddr)
		oui.Os = newSession.Get(c, sessionStatus.Os)
		oui.LoginTime, _ = strconv.ParseInt(newSession.Get(c, sessionStatus.LoginTime), 10, 64)
		if ol.UserName != "" && !strings.Contains(oui.UserName, ol.UserName) {
			continue
		}
		if ol.Ipaddr != "" && !strings.Contains(oui.Ipaddr, ol.Ipaddr) {
			continue
		}
		list = append(list, oui)
	}

	total = int64(len(list))
	return
}

func (userOnlineService *UserOnlineService) ForceLogout(c *gin.Context, tokenId string) {
	_ = session.NewAdminManager(userOnlineService.cache).Remove(c, tokenId)
}
