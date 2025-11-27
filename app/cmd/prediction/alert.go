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
	"nova-factory-server/app/cmd/prediction/condition"
	"nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
	"time"
)

type alert struct {
	metricCDao   metricDao.IMetricDao
	deviceMapDao deviceMonitorDao.IDeviceDataReportDao
}

func newAlert(metricCDao metricDao.IMetricDao, deviceMapDao deviceMonitorDao.IDeviceDataReportDao) *alert {
	return &alert{
		metricCDao:   metricCDao,
		deviceMapDao: deviceMapDao,
	}
}

func (a *alert) predict(config *aiDataSetModels.SysAiPrediction) {
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
	result, err := a.metricCDao.Query(&ctx, &metricModels.MetricDataQueryReq{
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
	list, err := a.deviceMapDao.GetDevList(&ctx, devList)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
		return
	}

	var devMap map[string]deviceMonitorModel.SysIotDbDevMapData = make(map[string]deviceMonitorModel.SysIotDbDevMapData)
	for _, dev := range list {
		devMap[dev.DevName] = dev
	}

	var advanced logalert.Advanced
	err = json.Unmarshal([]byte(config.Advanced), &advanced)
	if err != nil {
		zap.L().Error("json Unmarshal fail", zap.Error(err))
	}

	if len(result.Values) > 0 {
		if a.judge(config, &advanced, result.Values) {
			// 触发告警信息
			fmt.Println(111)
		}
		return
	}

	if len(result.MultiValues) > 0 {
		for _, v := range result.MultiValues {
			if a.judge(config, &advanced, v) {
				// 触发告警信息
				fmt.Println(222)
			}
		}
		return
	}

	return
}

// judge 判断告警条件
func (a *alert) judge(config *aiDataSetModels.SysAiPrediction, advanced *logalert.Advanced, values []metricModels.MetricQueryValue) bool {

	if len(values) < int(config.Threshold) {
		return false
	}

	if len(advanced.Rules) == 0 {
		return false
	}

	var ruleRet bool = true

	for _, rule := range advanced.Rules {
		var threshold int64 = 0
		for _, value := range values {

			for _, group := range rule.Groups {
				operatorFun, ok := condition.OperatorMap[group.Operator]
				if !ok {
					continue
				}
				ret, err := operatorFun(value.Value, group.Value)
				if err != nil {
					zap.L().Error("规则比较失败", zap.Error(err))
					continue
				}

				if !ret {
					continue
				}
				threshold++
				if config.Threshold <= threshold {
					if rule.MatchType == condition.MatchTypeAny {
						ruleRet = true
						break
					}
				} else {
					if rule.MatchType == condition.MatchTypeAny {
						ruleRet = false
					} else {
						return false
					}
				}

			}
		}
	}

	return ruleRet
}
