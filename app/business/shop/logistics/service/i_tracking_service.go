package service

import (
	"nova-factory-server/app/business/shop/logistics/models"

	"github.com/gin-gonic/gin"
)

// ITrackingService 物流轨迹查询服务接口
type ITrackingService interface {
	// Query 即时查询物流轨迹（缓存优先策略）
	Query(c *gin.Context, outsid, companyCode string) (*models.TrackingQueryResponse, error)
}
