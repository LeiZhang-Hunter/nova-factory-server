package middlewares

import (
	"net/http"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type sessionAppAuthMiddlewareBuilder struct {
	cache       cache.Cache
	sessionType string
}

func NewShopSessionAppAuthMiddlewareBuilder(cache cache.Cache) *sessionAppAuthMiddlewareBuilder {
	return &sessionAppAuthMiddlewareBuilder{
		cache:       cache,
		sessionType: sessionStatus.SessionTypeShopUser,
	}
}

func (s *sessionAppAuthMiddlewareBuilder) Build() func(c *gin.Context) {
	return func(c *gin.Context) {
		manager := session.NewManger(s.cache)
		currentSession, err := manager.GetSession(c)
		if err != nil {
			baizeContext.InvalidToken(c)
			c.Abort()
			return
		}
		if s.sessionType != "" && currentSession.Get(c, sessionStatus.SessionType) != s.sessionType {
			baizeContext.InvalidToken(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func NewShopSessionAppWsAuthMiddlewareBuilder(cache cache.Cache) *sessionAppAuthMiddlewareBuilder {
	return &sessionAppAuthMiddlewareBuilder{
		cache:       cache,
		sessionType: sessionStatus.SessionTypeShopUser,
	}
}

func (s *sessionAppAuthMiddlewareBuilder) BuildForWebSocket() func(c *gin.Context) {
	return func(c *gin.Context) {
		manager := session.NewManger(s.cache)
		currentSession, err := manager.GetSession(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if s.sessionType != "" && currentSession.Get(c, sessionStatus.SessionType) != s.sessionType {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
