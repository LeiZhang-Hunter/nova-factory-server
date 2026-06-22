package embedder

import (
	"github.com/gin-gonic/gin"
)

type Embedder interface {
	EmbeddingWithLLM(c *gin.Context) (EmbedderLlm, error)
}

type EmptyEmbedder struct{}

func NewEmptyEmbedder() *EmptyEmbedder {
	return &EmptyEmbedder{}
}

func (e *EmptyEmbedder) EmbeddingWithLLM(c *gin.Context) (EmbedderLlm, error) {
	return &EmbedderLlmData{}, nil
}
