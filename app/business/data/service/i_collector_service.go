package service

import (
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"

	"github.com/gin-gonic/gin"
)

// ICollectorService 采集器服务接口
type ICollectorService interface {
	CreateCollector(c *gin.Context, req *dto.CreateCollectorRequest) (*entity.Collector, error)
	GetCollector(c *gin.Context, id string) (*entity.Collector, error)
	ListCollectors(c *gin.Context, offset, limit int, name, deviceID, status string) ([]entity.Collector, int64, error)
	ListAllOnline(c *gin.Context) ([]entity.Collector, error)
	UpdateCollector(c *gin.Context, id string, req *dto.UpdateCollectorRequest) (*entity.Collector, error)
	DeleteCollector(c *gin.Context, id string) error
}
