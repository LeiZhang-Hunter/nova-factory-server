package aidatasetdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAiPredictionControlDao interface {
	Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionControl) (*aidatasetmodels.SysAiPredictionControl, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aidatasetmodels.SysAiPredictionControlListReq) (*aidatasetmodels.SysAiPredictionControlList, error)
	Find(c *gin.Context) (*aidatasetmodels.SysAiPredictionControl, error)
}
