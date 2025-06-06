package daemonizeServiceImpl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
	"go.uber.org/zap"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"sync"
	"time"
)

type DaemonizeServiceImpl struct {
	dao        daemonizeDao.IotAgentDao
	processDao daemonizeDao.IotAgentProcess
	configDao  daemonizeDao.IotAgentConfigDao
	manager    *ManagerServiceImpl
}

func NewDaemonizeServiceImpl(dao daemonizeDao.IotAgentDao, processDao daemonizeDao.IotAgentProcess, configDao daemonizeDao.IotAgentConfigDao) daemonizeService.DaemonizeService {
	return &DaemonizeServiceImpl{
		dao:        dao,
		configDao:  configDao,
		processDao: processDao,
		manager:    NewManagerServiceImpl(),
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

func (d *DaemonizeServiceImpl) AgentGetConfig(ctx context.Context, req *v1.AgentGetConfigReq) (*v1.AgentGetConfigRes, error) {
	res := &v1.AgentGetConfigRes{}
	config, err := d.configDao.GetByUuid(ctx, req.GetConfigUuid())
	if err != nil {
		return nil, err
	}
	if config != nil {

		// 转化配置
		var watcherMenConfig v1.ManagerConfig
		err = json.Unmarshal([]byte(config.Content), &watcherMenConfig)
		if err != nil {
			zap.L().Error("configDao.GetByUuid error:", zap.Error(err))
			return nil, err
		}
		content, err := json.Marshal(&watcherMenConfig)
		if err != nil {
			zap.L().Error("configDao.GetByUuid error:", zap.Error(err))
			return nil, err
		}
		res.Content = string(content)
	}
	return res, nil
}

// AgentHeartbeat grpc接口 agent心跳
func (d *DaemonizeServiceImpl) AgentHeartbeat(ctx context.Context, req *v1.AgentHeartbeatReq) (res *v1.AgentHeartbeatRes, err error) {
	agent, err := d.dao.GetByObjectId(ctx, req.ObjectId)
	if err != nil {
		return
	}
	configUuid := agent.ConfigUUID
	agent = &daemonizeModels.SysIotAgent{
		ObjectID: int64(req.GetObjectId()),
		Name:     req.GetName(),
		Version:  req.GetVersion(),
	}
	agentProcessList := make([]*daemonizeModels.SysIotAgentProcess, 0, len(req.AgentProcessInfo.ProcessList))
	for _, processInfo := range req.AgentProcessInfo.ProcessList {
		agentProcessList = append(agentProcessList, &daemonizeModels.SysIotAgentProcess{
			AgentObjectID: agent.ObjectID,
			Status:        int32(processInfo.State),
			Name:          processInfo.Name,
			Version:       processInfo.Version,
			StartTime:     processInfo.StartTime.AsTime(),
		})
	}

	err = d.processDao.RecordHeardBeat(ctx, req.ObjectId, agentProcessList)
	if err != nil {
		zap.L().Error("recordHeardBeat error:", zap.Error(err))
		return nil, err
	}
	res = &v1.AgentHeartbeatRes{}
	// 是否要更新配置
	if configUuid != req.ConfigUuid {
		res.ConfigUuid = configUuid
	}
	return
}

func (d *DaemonizeServiceImpl) AgentOperate(ctx context.Context, req *v1.AgentOperateReq, stream v1.AgentControllerService_AgentOperateServer) (err error) {
	clientId := req.GetObjectId()
	if clientId == 0 {
		zap.L().Error("AgentOperate error: agent id is 0")
		return gerror.NewCode(gcode.CodeInvalidParameter, "agent id cannot be empty")
	}

	g.Log().Debugf(stream.Context(), "agent add client stream, id: %v", clientId)
	d.manager.AddClient(clientId, stream)
	defer d.manager.DeleteClient(clientId)
	for {
		if d.manager.IsStopped() {
			g.Log().Debugf(stream.Context(), "server stopped, release client stream, id: %v", clientId)
			return nil
		}

		err = stream.Context().Err()
		if errors.Is(err, context.Canceled) {
			return nil
		} else if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}

// AgentOperateProcess grpc接口 操作agent进程
func (d *DaemonizeServiceImpl) AgentOperateProcess(ctx context.Context, cmd v1.AgentCmd, processOperateInfoList []*v1.ProcessOperateInfo) {
	wg := sync.WaitGroup{}
	wg.Add(len(processOperateInfoList))
	for _, info := range processOperateInfoList {
		go func() {
			defer wg.Done()
			stream := d.manager.getClient(info.AgentObjectId)
			if stream == nil {
				g.Log().Warningf(ctx, "get client empty:%v, object_id:%d", cmd, info.AgentObjectId)
				return
			}
			err := stream.Send(&v1.AgentOperateRes{
				Cmd:   cmd,
				Names: info.Names,
			})
			if err != nil {
				g.Log().Errorf(ctx, "send cmd[%v] to agent[%v] error, err:%v", cmd, info.AgentObjectId, err)
			}
			g.Log().Infof(ctx, "send cmd[%v] to agent[%v] error", cmd, info.AgentObjectId)
			return
		}()
	}
	wg.Wait()
	return
}
