package middlewares

import (
	"net/http"
	"strings"

	"nova-factory-server/app/business/data/dao"

	"github.com/gin-gonic/gin"
)

const (
	// HeaderDeviceID 设备拉取接口的设备ID请求头
	HeaderDeviceID = "X-Device-Id"
	// HeaderDeviceToken 设备拉取接口的设备Token请求头
	HeaderDeviceToken = "X-Device-Token"
	// ContextCollectorID 认证通过后写入上下文的主键ID
	ContextCollectorID = "collector_id"
)

// NewDeviceAuthMiddleware 校验设备 ID + Token，供采集器设备接口使用
func NewDeviceAuthMiddleware(collectorDao dao.ICollectorDAO) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := strings.TrimSpace(c.GetHeader(HeaderDeviceID))
		token := strings.TrimSpace(c.GetHeader(HeaderDeviceToken))
		if deviceID == "" || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "设备认证失败"})
			return
		}
		col, err := collectorDao.GetByDeviceID(c, deviceID)
		if err != nil || col == nil || col.Token != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "设备认证失败"})
			return
		}
		c.Set(ContextCollectorID, col.ID)
		c.Next()
	}
}
