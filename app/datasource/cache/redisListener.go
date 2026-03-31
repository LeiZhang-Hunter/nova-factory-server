package cache

import (
	"context"
	"encoding/json"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/business/admin/system/systemmodels"
	"nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/setting"

	"go.uber.org/zap"
)

func NewRedisSubscribe(cache Cache, ss systemservice.ISseService, js monitorservice.IJobService) *RedisSubscribe {
	return &RedisSubscribe{cache: cache, ss: ss, js: js}
}

type RedisSubscribe struct {
	cache Cache
	ss    systemservice.ISseService
	js    monitorservice.IJobService
}

func (r *RedisSubscribe) Run() {
	if setting.Conf.Cluster {
		go r.SubscribeNotification()
		go r.SubscribeJob()
	}

}

func (r *RedisSubscribe) SubscribeNotification() {
	background := context.Background()
	subscribe := r.cache.Subscribe(background, "notification")
	defer subscribe.Close()
	ch := subscribe.Channel()
	for msg := range ch {
		var sse systemmodels.Sse
		err := json.Unmarshal([]byte(msg.Payload), &sse)
		if err != nil {
			zap.L().Error("sse unmarshal error", zap.Error(err))
			continue
		}
		r.ss.SendNotification(background, &sse)
	}

}

func (r *RedisSubscribe) SubscribeJob() {
	background := context.Background()
	subscribe := r.cache.Subscribe(background, "job")
	defer subscribe.Close()
	ch := subscribe.Channel()
	for msg := range ch {
		var jb monitormodels.JobRedis
		err := json.Unmarshal([]byte(msg.Payload), &jb)
		if err != nil {
			zap.L().Error("sse unmarshal error", zap.Error(err))
			continue
		}
		if jb.Type == 0 {
			r.js.StartRunCron(background, &jb)
		} else {
			r.js.DeleteRunCron(background, &jb)
		}
	}

}
