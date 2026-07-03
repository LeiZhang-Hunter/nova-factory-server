package key

import "github.com/gin-gonic/gin"

type keys interface {
	GetUserId(key string) int64
	GetTool(c *gin.Context, key string) ([]string, error)
}
