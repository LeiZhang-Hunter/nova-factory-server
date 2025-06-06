package daemonizeDaoImpl

import (
	"context"
	"encoding/json"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
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
	i.cache.Set(ctx, key, string(content), 0)
	return nil
}
