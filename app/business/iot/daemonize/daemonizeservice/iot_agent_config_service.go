package daemonizeservice

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"

	"github.com/gin-gonic/gin"
)

type IotAgentConfigService interface {
	Create(ctx *gin.Context, config *daemonizemodels.SysIotAgentConfigSetReq) (*daemonizemodels.SysIotAgentConfig, error)
	Update(ctx *gin.Context, config *daemonizemodels.SysIotAgentConfigSetReq) (*daemonizemodels.SysIotAgentConfig, error)
	List(c *gin.Context, req *daemonizemodels.SysIotAgentConfigListReq) (*daemonizemodels.SysIotAgentConfigListData, error)
	GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizemodels.SysIotAgentConfig, error)
	Remove(ctx context.Context, ids []string) error
}
