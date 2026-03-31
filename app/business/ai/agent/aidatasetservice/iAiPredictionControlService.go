package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAiPredictionControlService interface {
	Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionControl) (*aidatasetmodels.SysAiPredictionControl, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aidatasetmodels.SysAiPredictionControlListReq) (*aidatasetmodels.SysAiPredictionControlList, error)
}
