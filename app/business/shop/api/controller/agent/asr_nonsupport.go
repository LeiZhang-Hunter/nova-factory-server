//go:build !ai

package agent

import "github.com/gin-gonic/gin"

type ASR struct{}

func NewASR() *ASR {
	return &ASR{}
}

func (*ASR) PrivateRoutes(_ *gin.RouterGroup) {}
