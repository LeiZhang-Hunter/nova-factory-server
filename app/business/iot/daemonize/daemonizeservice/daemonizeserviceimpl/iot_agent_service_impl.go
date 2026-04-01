package daemonizeserviceimpl

import (
	"context"
	daemonizeDao2 "nova-factory-server/app/business/iot/daemonize/daemonizedao"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"
	"nova-factory-server/app/business/iot/daemonize/daemonizeservice"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type iotAgentServiceImpl struct {
	dao        daemonizeDao2.IotAgentDao
	processDao daemonizeDao2.IotAgentProcess
}

func NewIotAgentServiceImpl(dao daemonizeDao2.IotAgentDao, processDao daemonizeDao2.IotAgentProcess) daemonizeservice.IotAgentService {
	return &iotAgentServiceImpl{
		dao:        dao,
		processDao: processDao,
	}
}

func (i *iotAgentServiceImpl) Add(ctx *gin.Context, req *daemonizemodels.SysIotAgentSetReq) (*daemonizemodels.SysIotAgent, error) {
	//i.dao.Update()
	data := daemonizemodels.ToSysIotAgent(req)
	data.ObjectID = uint64(snowflake.GenID())
	data.SetCreateBy(baizeContext.GetUserId(ctx))
	ret, err := i.dao.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (i *iotAgentServiceImpl) List(ctx *gin.Context, req *daemonizemodels.SysIotAgentListReq) (*daemonizemodels.SysIotAgentListData, error) {
	list, err := i.dao.GetAgentList(ctx, req)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	var objectIds []uint64
	for _, v := range list.Rows {
		objectIds = append(objectIds, v.ObjectID)
	}
	processes := i.processDao.GetHeardBeatInfo(ctx, objectIds)

	for _, v := range list.Rows {
		processList, ok := processes[v.ObjectID]
		if !ok {
			continue
		}
		v.Processes = processList
	}

	return list, err
}

func (i *iotAgentServiceImpl) Update(ctx *gin.Context, req *daemonizemodels.SysIotAgentSetReq) (*daemonizemodels.SysIotAgent, error) {
	data := daemonizemodels.ToSysIotAgent(req)
	data.SetUpdateBy(baizeContext.GetUserId(ctx))
	return i.dao.Update(ctx, data)
}

func (i *iotAgentServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}

func (i *iotAgentServiceImpl) GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizemodels.SysIotAgent, err error) {
	return i.dao.GetByObjectId(ctx, objectId)
}

func (i *iotAgentServiceImpl) UpdateConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error) {
	return i.dao.UpdateConfig(ctx, configId, objectIdList)
}
func (i *iotAgentServiceImpl) UpdateLastConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error) {
	return i.dao.UpdateLastConfig(ctx, configId, objectIdList)
}

func (i *iotAgentServiceImpl) Info(ctx *gin.Context, objectId uint64) (*daemonizemodels.SysIotAgent, error) {
	info, err := i.dao.GetByObjectId(ctx, objectId)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	processes := i.processDao.GetHeardBeatInfo(ctx, []uint64{objectId})
	processList, ok := processes[objectId]
	if !ok {
		return info, nil
	}
	info.Processes = processList

	return info, err
}
