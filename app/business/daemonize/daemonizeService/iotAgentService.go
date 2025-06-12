package daemonizeService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentService interface {
	Add(ctx *gin.Context, req *daemonizeModels.SysIotAgentSetReq) (*daemonizeModels.SysIotAgent, error)
	List(ctx *gin.Context, req *daemonizeModels.SysIotAgentListReq) (*daemonizeModels.SysIotAgentListData, error)
	Update(ctx *gin.Context, req *daemonizeModels.SysIotAgentSetReq) (*daemonizeModels.SysIotAgent, error)
	Remove(c *gin.Context, ids []string) error
	GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizeModels.SysIotAgent, err error)
	UpdateConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
	UpdateLastConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
	Info(ctx *gin.Context, objectId uint64) (*daemonizeModels.SysIotAgent, error)
}
