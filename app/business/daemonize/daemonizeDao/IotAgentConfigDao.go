package daemonizeDao

import (
	"context"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentConfigDao interface {
	GetByUuid(ctx context.Context, uuid string) (*daemonizeModels.SysIotAgentConfig, error)
	GetLastedConfig(ctx context.Context) (*daemonizeModels.SysIotAgentConfig, error)
	GetLastedConfigList(ctx context.Context, count int) ([]*daemonizeModels.SysIotAgentConfig, error)
	GetVersionListByUuidList(ctx context.Context, uuidList []string) (versionMap map[string]string, err error)
	GetByVersion(ctx context.Context, configVersion string) (config *daemonizeModels.SysIotAgentConfig, err error)
	Create(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (err error)
}
