package daemonizedaoimpl

import (
	"context"
	"encoding/json"
	"nova-factory-server/app/business/iot/daemonize/daemonizedao"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"time"

	"go.uber.org/zap"
)

type IotAgentProcessDaoImpl struct {
	cache cache.Cache
}

func NewIotAgentProcessDaoImpl(cache cache.Cache) daemonizedao.IotAgentProcess {
	return &IotAgentProcessDaoImpl{
		cache: cache,
	}
}

func (i *IotAgentProcessDaoImpl) RecordHeardBeat(ctx context.Context, objectId uint64, processes []*daemonizemodels.SysIotAgentProcess) error {
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

func (i *IotAgentProcessDaoImpl) GetHeardBeatInfo(ctx context.Context, objectIds []uint64) map[uint64][]*daemonizemodels.SysIotAgentProcess {
	processes := make(map[uint64][]*daemonizemodels.SysIotAgentProcess)
	if objectIds == nil || len(objectIds) == 0 {
		return processes
	}
	var objectIdMap map[uint64]uint64 = make(map[uint64]uint64)
	var objectIdList []string
	var objectIdsArray []uint64 = make([]uint64, 0)
	for _, id := range objectIds {
		objectIdMap[id] = id
	}
	for _, id := range objectIdMap {
		objectIdList = append(objectIdList, redis.MakeCacheKey(redis.AgentProcessCacheKey, "", id))
		objectIdsArray = append(objectIdsArray, id)
	}
	slice := i.cache.MGet(ctx, objectIdList).Val()
	for k, v := range slice {
		var process []*daemonizemodels.SysIotAgentProcess
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
		processes[objectIdsArray[k]] = process
	}
	return processes
}
