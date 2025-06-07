package daemonizeServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type iotAgentServiceImpl struct {
	dao daemonizeDao.IotAgentDao
}

func NewIotAgentServiceImpl(dao daemonizeDao.IotAgentDao) daemonizeService.IotAgentService {
	return &iotAgentServiceImpl{
		dao: dao,
	}
}

func (i *iotAgentServiceImpl) Add(ctx *gin.Context, req *daemonizeModels.SysIotAgentSetReq) (*daemonizeModels.SysIotAgent, error) {
	//i.dao.Update()
	data := daemonizeModels.ToSysIotAgent(req)
	data.ObjectID = uint64(snowflake.GenID())
	data.SetCreateBy(baizeContext.GetUserId(ctx))
	ret, err := i.dao.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (i *iotAgentServiceImpl) List(ctx *gin.Context, req *daemonizeModels.SysIotAgentListReq) (*daemonizeModels.SysIotAgentListData, error) {
	return i.dao.GetAgentList(ctx, req)
}

func (i *iotAgentServiceImpl) Update(ctx *gin.Context, req *daemonizeModels.SysIotAgentSetReq) (*daemonizeModels.SysIotAgent, error) {
	data := daemonizeModels.ToSysIotAgent(req)
	data.SetUpdateBy(baizeContext.GetUserId(ctx))
	return i.dao.Update(ctx, data)
}

func (i *iotAgentServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
