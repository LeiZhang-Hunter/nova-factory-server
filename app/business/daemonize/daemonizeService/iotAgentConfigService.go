package daemonizeService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentConfigService interface {
	Create(ctx *gin.Context, config *daemonizeModels.SysIotAgentConfigSetReq) (*daemonizeModels.SysIotAgentConfig, error)
	Update(ctx *gin.Context, config *daemonizeModels.SysIotAgentConfigSetReq) (*daemonizeModels.SysIotAgentConfig, error)
	List(c *gin.Context, req *daemonizeModels.SysIotAgentConfigListReq) (*daemonizeModels.SysIotAgentConfigListData, error)
	GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizeModels.SysIotAgentConfig, error)
}
