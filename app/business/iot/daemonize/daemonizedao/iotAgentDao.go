package daemonizedao

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"

	"github.com/gin-gonic/gin"
)

type IotAgentDao interface {
	GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizemodels.SysIotAgent, err error)
	UpdateHeartBeat(ctx context.Context, data *daemonizemodels.SysIotAgent) error
	Update(ctx context.Context, data *daemonizemodels.SysIotAgent) (*daemonizemodels.SysIotAgent, error)
	Create(ctx context.Context, doAgent *daemonizemodels.SysIotAgent) (*daemonizemodels.SysIotAgent, error)
	GetAgentList(ctx *gin.Context, req *daemonizemodels.SysIotAgentListReq) (*daemonizemodels.SysIotAgentListData, error)
	Remove(c *gin.Context, ids []string) error
	UpdateConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
	UpdateLastConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
}
