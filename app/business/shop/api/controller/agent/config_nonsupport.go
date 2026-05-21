//go:build !ai

package agent

import "github.com/gin-gonic/gin"

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (*Config) PrivateRoutes(_ *gin.RouterGroup) {}
