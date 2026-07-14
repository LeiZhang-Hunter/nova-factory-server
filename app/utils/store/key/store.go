package key

import "github.com/gin-gonic/gin"

type emptyKeys struct{}

func newEmptyPermissions() keys {
	return &emptyKeys{}
}

func (e *emptyKeys) GetUserId(key string) int64 {
	return 0
}

func (e *emptyKeys) GetTool(c *gin.Context, key string) ([]string, error) {
	return []string{}, nil
}

func (e *emptyKeys) GetInfo(c *gin.Context, key string) (*Info, error) {
	return nil, nil
}
