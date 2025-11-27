package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"time"
)

type exception struct {
	metricCDao   metricDao.IMetricDao
	deviceMapDao deviceMonitorDao.IDeviceDataReportDao
	judge        *judge
}

func newException(metricCDao metricDao.IMetricDao, deviceMapDao deviceMonitorDao.IDeviceDataReportDao) *exception {
	return &exception{
		metricCDao:   metricCDao,
		deviceMapDao: deviceMapDao,
		judge:        newJudge(),
	}
}

func (e *exception) predict(config *aiDataSetModels.SysAiPredictionException) {
	if config == nil {
		return
	}
	//config.Dev
	var devList []string
	err := json.Unmarshal([]byte(config.Dev), &devList)
	if err != nil {
		zap.L().Error("json unmarshal error", zap.Error(err))
		return
	}

	// 没有测点
	if len(devList) == 0 {
		return
	}

	var name string = ""
	for k, dev := range devList {
		if k == len(devList)-1 {
			name += dev
		} else {
			name += dev + ","
		}
	}

	end := time.Now().UnixMilli()
	start := end - config.Interval*1000
	var exceptionStr string = fmt.Sprintf("%s(value)", config.AggFunction)

	var level int = 0
	var ctx gin.Context
	result, err := e.metricCDao.Query(&ctx, &metricModels.MetricDataQueryReq{
		Type:       "line",
		Name:       name,
		Start:      uint64(start),
		End:        uint64(end),
		Step:       0,
		Interval:   0,
		Field:      " ",
		Level:      &level,
		Expression: exceptionStr,
		Predict: metricModels.Predict{
			Model:  config.Model,
			Enable: true,
		},
	})
	if err != nil {
		zap.L().Error("query metric error", zap.Error(err))
		return
	}

	// 读取测点列表
	list, err := e.deviceMapDao.GetDevList(&ctx, devList)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
		return
	}

	var devMap map[string]deviceMonitorModel.SysIotDbDevMapData = make(map[string]deviceMonitorModel.SysIotDbDevMapData)
	for _, dev := range list {
		devMap[dev.DevName] = dev
	}

	if len(result.Values) > 0 {
		if e.judge.judgeException(config, result.Values) {
			// 触发告警信息
			fmt.Println(111)
		}
		return
	}

	if len(result.MultiValues) > 0 {
		for _, v := range result.MultiValues {
			if e.judge.judgeException(config, v) {
				// 触发告警信息
				fmt.Println(222)
			}
		}
		return
	}

	return
}
