package daemonizeDao

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentDao interface {
	GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizeModels.SysIotAgent, err error)
	UpdateHeartBeat(ctx context.Context, data *daemonizeModels.SysIotAgent) error
	Update(ctx context.Context, data *daemonizeModels.SysIotAgent) (*daemonizeModels.SysIotAgent, error)
	Create(ctx context.Context, doAgent *daemonizeModels.SysIotAgent) (*daemonizeModels.SysIotAgent, error)
	GetAgentList(ctx *gin.Context, req *daemonizeModels.SysIotAgentListReq) (*daemonizeModels.SysIotAgentListData, error)
	Remove(c *gin.Context, ids []string) error
	UpdateConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
	UpdateLastConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error)
}
