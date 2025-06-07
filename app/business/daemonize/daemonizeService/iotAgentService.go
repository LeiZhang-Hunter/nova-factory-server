package daemonizeService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentService interface {
	Add(ctx *gin.Context, req *daemonizeModels.SysIotAgentSetReq) (*daemonizeModels.SysIotAgent, error)
	List(ctx *gin.Context, req *daemonizeModels.SysIotAgentListReq) (*daemonizeModels.SysIotAgentListData, error)
	Update(ctx *gin.Context, req *daemonizeModels.SysIotAgentSetReq) (*daemonizeModels.SysIotAgent, error)
	Remove(c *gin.Context, ids []string) error
}
