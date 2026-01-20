package daemonizeServiceImpl

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IotAgentConfigServiceImpl struct {
	dao daemonizeDao.IotAgentConfigDao
}

func NewIotAgentConfigServiceImpl(dao daemonizeDao.IotAgentConfigDao) daemonizeService.IotAgentConfigService {
	return &IotAgentConfigServiceImpl{
		dao: dao,
	}
}

func (i *IotAgentConfigServiceImpl) Create(ctx *gin.Context, config *daemonizeModels.SysIotAgentConfigSetReq) (*daemonizeModels.SysIotAgentConfig, error) {
	data := daemonizeModels.OfSysIotAgentConfig(config)
	data.ID = uint64(snowflake.GenID())
	data.DeptID = baizeContext.GetDeptId(ctx)
	data.SetCreateBy(baizeContext.GetUserId(ctx))
	return i.dao.Create(ctx, data)
}

func (i *IotAgentConfigServiceImpl) Update(ctx *gin.Context, config *daemonizeModels.SysIotAgentConfigSetReq) (*daemonizeModels.SysIotAgentConfig, error) {
	data := daemonizeModels.OfSysIotAgentConfig(config)
	data.SetUpdateBy(baizeContext.GetUserId(ctx))
	return i.dao.Update(ctx, data)
}

func (i *IotAgentConfigServiceImpl) List(c *gin.Context, req *daemonizeModels.SysIotAgentConfigListReq) (*daemonizeModels.SysIotAgentConfigListData, error) {
	return i.dao.List(c, req)
}

func (i *IotAgentConfigServiceImpl) GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizeModels.SysIotAgentConfig, error) {
	return i.dao.GetLastedConfig(ctx, agentId)
}

func (i *IotAgentConfigServiceImpl) Remove(ctx context.Context, ids []string) error {
	return i.dao.Remove(ctx, ids)
}
