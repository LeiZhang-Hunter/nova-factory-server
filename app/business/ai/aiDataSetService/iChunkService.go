package aiDataSetService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IChunkService interface {
	ChunkList(c *gin.Context, req *aiDataSetModels.ChunkListReq) (*aiDataSetModels.ChunkListResponse, error)
	AddChunk(c *gin.Context, req *aiDataSetModels.AddChunkReq) (*aiDataSetModels.AddChunkResponse, error)
	RemoveChunk(c *gin.Context, req *aiDataSetModels.RemoveChunkReq) error
	UpdateChunk(c *gin.Context, req *aiDataSetModels.UpdateChunkReq) (*aiDataSetModels.UpdateChunkResponse, error)
	RetrievalChunk(c *gin.Context, req *aiDataSetModels.RetrievalListReq) (*aiDataSetModels.RetrievalApiListResponse, error)
}
