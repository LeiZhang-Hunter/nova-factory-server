package daemonizeServiceImpl

import (
	"context"
	"errors"
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
	"go.uber.org/zap"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeService"
)

type DaemonizeServiceImpl struct {
	dao daemonizeDao.IotAgentDao
}

func NewDaemonizeServiceImpl(dao daemonizeDao.IotAgentDao) daemonizeService.DaemonizeService {
	return &DaemonizeServiceImpl{
		dao: dao,
	}
}

// AgentRegister agent注册
func (d *DaemonizeServiceImpl) AgentRegister(ctx context.Context, req *v1.AgentRegisterReq) (res *v1.AgentRegisterRes, err error) {
	res = &v1.AgentRegisterRes{}
	agent, err := d.dao.GetByObjectId(ctx, req.GetObjectId())
	if err != nil {
		zap.L().Error("dao.GetByObjectId error:", zap.Error(err))
		return
	}

	if agent == nil {
		return nil, errors.New("agent not exist")
	}

	agent.Version = req.Version
	agent.Ipv4 = req.Ipv4
	agent.Ipv6 = req.Ipv6
	err = d.dao.Update(ctx, agent)
	if err != nil {
		return nil, err
	}
	//agent.ConfigUUID = req.
	// 更新心跳包
	err = d.dao.UpdateHeartBeat(ctx, agent)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DaemonizeServiceImpl) AgentGetConfig(context.Context, *v1.AgentGetConfigReq) (*v1.AgentGetConfigRes, error) {
	return nil, nil
}
