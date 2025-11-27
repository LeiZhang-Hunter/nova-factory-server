package main

import (
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/metric/device/metricDao"
	"sync"
	"time"
)

type Runner struct {
	// predictionDao 预测策略数据源
	predictionDao aiDataSetDao.IAiPredictionListDao
	// predictionDao 预测异常策略数据源
	predictionExceptionDao aiDataSetDao.IAiPredictionExceptionDao
	metricCDao             metricDao.IMetricDao
	pool                   *ants.Pool
	alert                  *alert
	exception              *exception
	wait                   sync.WaitGroup
	done                   chan struct{}
}

func NewRunner(predictionDao aiDataSetDao.IAiPredictionListDao, predictionExceptionDao aiDataSetDao.IAiPredictionExceptionDao,
	metricCDao metricDao.IMetricDao, deviceMapDao deviceMonitorDao.IDeviceDataReportDao) *Runner {
	pool, err := ants.NewPool(10)
	if err != nil {
		panic(err)
		return nil
	}
	return &Runner{
		predictionDao:          predictionDao,
		predictionExceptionDao: predictionExceptionDao,
		pool:                   pool,
		metricCDao:             metricCDao,
		done:                   make(chan struct{}),
		alert:                  newAlert(metricCDao, deviceMapDao),
		exception:              newException(metricCDao, deviceMapDao),
	}
}

func (r *Runner) Run() {
	r.wait.Add(1)
	minute := 1 * time.Second
	go func() {
		defer r.wait.Done()
		timer := time.NewTimer(minute)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				{
					r.execute()
					timer.Reset(minute)
					break
				}
			case <-r.done:
				{
					return
				}

			}
		}
	}()
}

func (r *Runner) execute() {
	ctx := gin.Context{}
	//all, err := r.predictionDao.All(&ctx)
	//if err != nil {
	//	zap.L().Error("读取预测告警任务失败")
	//}
	//
	//if len(all) > 0 {
	//	for _, v := range all {
	//		r.alert.predict(v)
	//	}
	//}

	exceptions, err := r.predictionExceptionDao.All(&ctx)
	if err != nil {
		zap.L().Error("读取异常预测告警任务失败")
	}

	if len(exceptions) > 0 {
		for _, v := range exceptions {
			r.exception.predict(v)
		}
	}
}

func (r *Runner) stop() {
	close(r.done)
	r.wait.Wait()
}
