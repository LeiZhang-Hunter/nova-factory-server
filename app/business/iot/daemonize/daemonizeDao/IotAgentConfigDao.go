package daemonizeDao

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizeModels"

	"github.com/gin-gonic/gin"
)

type IotAgentConfigDao interface {
	GetByUuid(ctx context.Context, uuid string) (*daemonizeModels.SysIotAgentConfig, error)
	GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizeModels.SysIotAgentConfig, error)
	GetLastedConfigList(ctx context.Context, count int) ([]*daemonizeModels.SysIotAgentConfig, error)
	GetVersionListByUuidList(ctx context.Context, uuidList []string) (versionMap map[uint64]string, err error)
	GetByVersion(ctx context.Context, configVersion string) (config *daemonizeModels.SysIotAgentConfig, err error)
	Create(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (*daemonizeModels.SysIotAgentConfig, error)
	Update(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (*daemonizeModels.SysIotAgentConfig, error)
	List(c *gin.Context, req *daemonizeModels.SysIotAgentConfigListReq) (*daemonizeModels.SysIotAgentConfigListData, error)
	Remove(ctx context.Context, ids []string) error
}
