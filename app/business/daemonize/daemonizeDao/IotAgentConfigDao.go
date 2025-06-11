package daemonizeDao

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentConfigDao interface {
	GetByUuid(ctx context.Context, uuid string) (*daemonizeModels.SysIotAgentConfig, error)
	GetLastedConfig(ctx context.Context) (*daemonizeModels.SysIotAgentConfig, error)
	GetLastedConfigList(ctx context.Context, count int) ([]*daemonizeModels.SysIotAgentConfig, error)
	GetVersionListByUuidList(ctx context.Context, uuidList []string) (versionMap map[uint64]string, err error)
	GetByVersion(ctx context.Context, configVersion string) (config *daemonizeModels.SysIotAgentConfig, err error)
	Create(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (*daemonizeModels.SysIotAgentConfig, error)
	Update(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (*daemonizeModels.SysIotAgentConfig, error)
	List(c *gin.Context, req *daemonizeModels.SysIotAgentConfigListReq) (*daemonizeModels.SysIotAgentConfigListData, error)
}
