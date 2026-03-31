package daemonizeservice

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"

	"github.com/gin-gonic/gin"
)

type IotAgentService interface {
	Add(ctx *gin.Context, req *daemonizemodels.SysIotAgentSetReq) (*daemonizemodels.SysIotAgent, error)
	List(ctx *gin.Context, req *daemonizemodels.SysIotAgentListReq) (*daemonizemodels.SysIotAgentListData, error)
	Update(ctx *gin.Context, req *daemonizemodels.SysIotAgentSetReq) (*daemonizemodels.SysIotAgent, error)
	Remove(c *gin.Context, ids []string) error
	GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizemodels.SysIotAgent, err error)
	UpdateConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
	UpdateLastConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
	Info(ctx *gin.Context, objectId uint64) (*daemonizemodels.SysIotAgent, error)
}
