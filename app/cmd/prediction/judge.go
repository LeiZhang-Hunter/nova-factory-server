package main

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/cmd/prediction/condition"
	"nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
)

type judge struct {
}

func newJudge() *judge {
	return &judge{}
}

// judge 判断告警条件
func (a *judge) judge(config *aiDataSetModels.SysAiPrediction, advanced *logalert.Advanced, values []metricModels.MetricQueryValue) bool {

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

func (a *judge) judgeException(config *aiDataSetModels.SysAiPredictionException, values []metricModels.MetricQueryValue) bool {

	if len(values) < int(config.Threshold) {
		return false
	}

	var ruleRet bool = true

	//var threshold int64 = 0
	//for _, value := range values {
	//
	//	operatorFun, ok := condition.OperatorMap[group.Operator]
	//	if !ok {
	//		continue
	//	}
	//
	//	threshold++
	//	if config.Threshold <= threshold {
	//		if rule.MatchType == condition.MatchTypeAny {
	//			ruleRet = true
	//			break
	//		}
	//	} else {
	//		if rule.MatchType == condition.MatchTypeAny {
	//			ruleRet = false
	//		} else {
	//			return false
	//		}
	//	}
	//
	//}

	return ruleRet
}
