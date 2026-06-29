package dao

import (
	"nova-factory-server/app/business/shop/logistics/models"

	"github.com/gin-gonic/gin"
)

// ITrackingDao 物流轨迹数据访问接口
type ITrackingDao interface {
	// GetSignedRecord 查询已签收的轨迹记录（DB 永久缓存）
	GetSignedRecord(c *gin.Context, outsid, companyCode string) (*models.TrackingQueryResponse, error)
	// SaveSignedRecord 保存已签收的轨迹记录
	SaveSignedRecord(c *gin.Context, record *models.TrackingRecordSet) error
}
