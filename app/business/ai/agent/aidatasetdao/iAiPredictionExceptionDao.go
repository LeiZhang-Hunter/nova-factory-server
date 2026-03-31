package aidatasetdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAiPredictionExceptionDao interface {
	Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionException) (*aidatasetmodels.SysAiPredictionException, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aidatasetmodels.SysAiPredictionExceptionListReq) (*aidatasetmodels.SysAiPredictionExceptionList, error)
	All(c *gin.Context) ([]*aidatasetmodels.SysAiPredictionException, error)
}
