package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAiModelProviderService interface {
	ListWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error)
}
