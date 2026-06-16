package integration

import (
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/config"

	"github.com/gin-gonic/gin"
)

// Integration 管家婆或者金蝶的配置
type Integration interface {
	GetService(c *gin.Context) (api.Service, config.Config, error)
}

type EmptyIntegrationStore struct{}

func NewEmptyIntegrationStore() Integration {
	return &EmptyIntegrationStore{}
}

func (*EmptyIntegrationStore) GetService(c *gin.Context) (api.Service, config.Config, error) {
	return nil, nil, nil
}
