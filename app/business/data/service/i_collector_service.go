package service

import (
	"nova-factory-server/app/business/data/models/dto"

	"github.com/gin-gonic/gin"
)

// ICollectorService 采集器服务接口
type ICollectorService interface {
	Create(c *gin.Context, req *dto.CreateCollectorRequest) (*dto.CollectorResponse, error)
	Get(c *gin.Context, id string) (*dto.CollectorResponse, error)
	List(c *gin.Context, offset, limit int, name, deviceID, status string) ([]dto.CollectorResponse, int64, error)
	ListOnline(c *gin.Context) ([]dto.CollectorResponse, error)
	Update(c *gin.Context, id string, req *dto.UpdateCollectorRequest) (*dto.CollectorResponse, error)
	Delete(c *gin.Context, id string) error
}
