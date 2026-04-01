package aidatasetservice

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"github.com/gin-gonic/gin"
)

type IChunkService interface {
	ChunkList(c *gin.Context, req *aidatasetmodels.ChunkListReq) (*aidatasetmodels.ChunkListResponse, error)
	AddChunk(c *gin.Context, req *aidatasetmodels.AddChunkReq) (*aidatasetmodels.AddChunkResponse, error)
	RemoveChunk(c *gin.Context, req *aidatasetmodels.RemoveChunkReq) error
	UpdateChunk(c *gin.Context, req *aidatasetmodels.UpdateChunkReq) (*aidatasetmodels.UpdateChunkResponse, error)
	RetrievalChunk(c *gin.Context, req *aidatasetmodels.RetrievalListReq) (*aidatasetmodels.RetrievalApiListResponse, error)
}
