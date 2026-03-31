package aidatasetdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAiPredictionListDao interface {
	Set(c *gin.Context, data *aidatasetmodels.SetSysAiPrediction) (*aidatasetmodels.SysAiPrediction, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aidatasetmodels.SysAiPredictionListReq) (*aidatasetmodels.SysAiPredictionList, error)
	// All 循环查找告警策略
	All(c *gin.Context) ([]*aidatasetmodels.SysAiPrediction, error)
}
