package middlewares

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
)

type optionalAuthMiddlewareBuilder struct {
	cache       cache.Cache
	sessionType string
}

// NewOptionalShopSessionAuthMiddlewareBuilder 创建可选认证中间件构建器。
// 当请求携带有效 session 时提取用户信息到 context，但不拦截未认证请求。
func NewOptionalShopSessionAuthMiddlewareBuilder(cache cache.Cache) *optionalAuthMiddlewareBuilder {
	return &optionalAuthMiddlewareBuilder{
		cache:       cache,
		sessionType: sessionStatus.SessionTypeShopUser,
	}
}

func (s *optionalAuthMiddlewareBuilder) Build() func(c *gin.Context) {
	return func(c *gin.Context) {
		manager := session.NewManger(s.cache)
		currentSession, err := manager.GetSession(c)
		if err != nil {
			c.Next()
			return
		}
		if s.sessionType != "" && currentSession.Get(c, sessionStatus.SessionType) != s.sessionType {
			c.Next()
			return
		}
		c.Next()
	}
}
