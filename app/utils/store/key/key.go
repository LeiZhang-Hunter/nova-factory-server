package key

import "github.com/gin-gonic/gin"

type Info struct {
	UserID int64
	DeptID int64
	Tools  []string
}

type keys interface {
	GetUserId(key string) int64
	GetTool(c *gin.Context, key string) ([]string, error)
	GetInfo(c *gin.Context, key string) (*Info, error)
}
