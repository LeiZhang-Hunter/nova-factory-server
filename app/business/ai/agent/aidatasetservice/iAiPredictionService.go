package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAiPredictionService interface {
	Set(c *gin.Context, data *aidatasetmodels.SetSysAiPrediction) (*aidatasetmodels.SysAiPrediction, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aidatasetmodels.SysAiPredictionListReq) (*aidatasetmodels.SysAiPredictionList, error)
}
