package alertService

import "github.com/gin-gonic/gin"

type RunnerService interface {
	Load(ctx *gin.Context, agentId string) (string, error)
}
