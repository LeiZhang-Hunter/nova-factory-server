package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/utils/store/embedder"
)

type IAiModelProviderService interface {
	ListWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error)
	EmbeddingWithLLM(c *gin.Context) (embedder.EmbedderLlm, error)
}
