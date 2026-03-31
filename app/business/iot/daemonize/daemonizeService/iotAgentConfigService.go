package daemonizeService

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizeModels"

	"github.com/gin-gonic/gin"
)

type IotAgentConfigService interface {
	Create(ctx *gin.Context, config *daemonizeModels.SysIotAgentConfigSetReq) (*daemonizeModels.SysIotAgentConfig, error)
	Update(ctx *gin.Context, config *daemonizeModels.SysIotAgentConfigSetReq) (*daemonizeModels.SysIotAgentConfig, error)
	List(c *gin.Context, req *daemonizeModels.SysIotAgentConfigListReq) (*daemonizeModels.SysIotAgentConfigListData, error)
	GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizeModels.SysIotAgentConfig, error)
	Remove(ctx context.Context, ids []string) error
}
