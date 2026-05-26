//go:build !ai

package agent

import "github.com/gin-gonic/gin"

type Voice struct{}

func NewVoice() *Voice {
	return &Voice{}
}

func (*Voice) PublicRoutes(_ *gin.RouterGroup) {}

func (*Voice) WsRegister(_ *gin.RouterGroup) {}
