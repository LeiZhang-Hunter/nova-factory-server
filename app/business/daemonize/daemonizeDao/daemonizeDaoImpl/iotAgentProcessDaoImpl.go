package daemonizeDaoImpl

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"time"
)

type IotAgentProcessDaoImpl struct {
	cache cache.Cache
}

func NewIotAgentProcessDaoImpl(cache cache.Cache) daemonizeDao.IotAgentProcess {
	return &IotAgentProcessDaoImpl{
		cache: cache,
	}
}

func (i *IotAgentProcessDaoImpl) RecordHeardBeat(ctx context.Context, objectId uint64, processes []*daemonizeModels.SysIotAgentProcess) error {
	if processes == nil {
		return nil
	}
	key := redis.MakeCacheKey(redis.AgentProcessCacheKey, "", objectId)
	content, err := json.Marshal(processes)
	if err != nil {
		return err
	}
	i.cache.Set(ctx, key, string(content), time.Duration(agent.CHECK_ONLINE_DURATION)*time.Second)
	return nil
}

func (i *IotAgentProcessDaoImpl) GetHeardBeatInfo(ctx context.Context, objectIds []uint64) map[uint64][]*daemonizeModels.SysIotAgentProcess {
	processes := make(map[uint64][]*daemonizeModels.SysIotAgentProcess)
	if objectIds == nil || len(objectIds) == 0 {
		return processes
	}
	var objectIdMap map[uint64]uint64 = make(map[uint64]uint64)
	var objectIdList []string
	for _, id := range objectIds {
		objectIdMap[id] = id
	}
	for _, id := range objectIdMap {
		objectIdList = append(objectIdList, redis.MakeCacheKey(redis.AgentProcessCacheKey, "", id))

	}
	slice := i.cache.MGet(ctx, objectIdList).Val()
	for k, v := range slice {
		var process []*daemonizeModels.SysIotAgentProcess
		str, ok := v.(string)
		if !ok {
			continue
		}
		if str == "" {
			continue
		}
		err := json.Unmarshal([]byte(str), &process)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
			continue
		}
		processes[objectIds[k]] = process
	}
	return processes
}
