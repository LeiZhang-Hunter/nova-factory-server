package daemonizeserviceimpl

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizedao"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"
	"nova-factory-server/app/business/iot/daemonize/daemonizeservice"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type IotAgentConfigServiceImpl struct {
	dao daemonizedao.IotAgentConfigDao
}

func NewIotAgentConfigServiceImpl(dao daemonizedao.IotAgentConfigDao) daemonizeservice.IotAgentConfigService {
	return &IotAgentConfigServiceImpl{
		dao: dao,
	}
}

func (i *IotAgentConfigServiceImpl) Create(ctx *gin.Context, config *daemonizemodels.SysIotAgentConfigSetReq) (*daemonizemodels.SysIotAgentConfig, error) {
	data := daemonizemodels.OfSysIotAgentConfig(config)
	data.ID = uint64(snowflake.GenID())
	data.DeptID = baizeContext.GetDeptId(ctx)
	data.SetCreateBy(baizeContext.GetUserId(ctx))
	return i.dao.Create(ctx, data)
}

func (i *IotAgentConfigServiceImpl) Update(ctx *gin.Context, config *daemonizemodels.SysIotAgentConfigSetReq) (*daemonizemodels.SysIotAgentConfig, error) {
	data := daemonizemodels.OfSysIotAgentConfig(config)
	data.SetUpdateBy(baizeContext.GetUserId(ctx))
	return i.dao.Update(ctx, data)
}

func (i *IotAgentConfigServiceImpl) List(c *gin.Context, req *daemonizemodels.SysIotAgentConfigListReq) (*daemonizemodels.SysIotAgentConfigListData, error) {
	return i.dao.List(c, req)
}

func (i *IotAgentConfigServiceImpl) GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizemodels.SysIotAgentConfig, error) {
	return i.dao.GetLastedConfig(ctx, agentId)
}

func (i *IotAgentConfigServiceImpl) Remove(ctx context.Context, ids []string) error {
	return i.dao.Remove(ctx, ids)
}
