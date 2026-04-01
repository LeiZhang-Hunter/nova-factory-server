package daemonizedao

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"

	"github.com/gin-gonic/gin"
)

type IotAgentConfigDao interface {
	GetByUuid(ctx context.Context, uuid string) (*daemonizemodels.SysIotAgentConfig, error)
	GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizemodels.SysIotAgentConfig, error)
	GetLastedConfigList(ctx context.Context, count int) ([]*daemonizemodels.SysIotAgentConfig, error)
	GetVersionListByUuidList(ctx context.Context, uuidList []string) (versionMap map[uint64]string, err error)
	GetByVersion(ctx context.Context, configVersion string) (config *daemonizemodels.SysIotAgentConfig, err error)
	Create(ctx context.Context, config *daemonizemodels.SysIotAgentConfig) (*daemonizemodels.SysIotAgentConfig, error)
	Update(ctx context.Context, config *daemonizemodels.SysIotAgentConfig) (*daemonizemodels.SysIotAgentConfig, error)
	List(c *gin.Context, req *daemonizemodels.SysIotAgentConfigListReq) (*daemonizemodels.SysIotAgentConfigListData, error)
	Remove(ctx context.Context, ids []string) error
}
