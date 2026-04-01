package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type IConfigDao interface {
	SelectConfigList(ctx context.Context, config *systemmodels.SysConfigDQL) (sysConfigList []*systemmodels.SysConfigVo, total int64)
	SelectConfigListAll(ctx context.Context, config *systemmodels.SysConfigDQL) (list []*systemmodels.SysConfigVo)
	SelectConfigById(ctx context.Context, configId int64) (Config *systemmodels.SysConfigVo)
	InsertConfig(ctx context.Context, config *systemmodels.SysConfigVo)
	UpdateConfig(ctx context.Context, config *systemmodels.SysConfigVo)
	DeleteConfigById(ctx context.Context, configId int64)
	SelectConfigIdByConfigKey(ctx context.Context, configKey string) int64
	SelectConfigValueByConfigKey(ctx context.Context, configKey string) string
}
