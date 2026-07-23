package systemdao

import (
	"context"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IConfigDao interface {
	SelectConfigList(ctx context.Context, config *modelquery.SysConfigDQL) (sysConfigList []*modelresponse.SysConfigVo, total int64)
	SelectConfigListAll(ctx context.Context, config *modelquery.SysConfigDQL) (list []*modelresponse.SysConfigVo)
	SelectConfigById(ctx context.Context, configId int64) (Config *modelresponse.SysConfigVo)
	InsertConfig(ctx context.Context, config *modelrequest.SysConfigDML)
	UpdateConfig(ctx context.Context, config *modelrequest.SysConfigDML)
	DeleteConfigById(ctx context.Context, configId int64)
	SelectConfigIdByConfigKey(ctx context.Context, configKey string) int64
	SelectConfigValueByConfigKey(ctx context.Context, configKey string) string
}
