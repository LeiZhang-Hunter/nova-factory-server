package middlewares

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/baizeContext"
)

type sessionAuthMiddlewareBuilder struct {
	paths       baize.Set[string]
	cache       cache.Cache
	sessionType string
}

func NewSessionAuthMiddlewareBuilder(cache cache.Cache) *sessionAuthMiddlewareBuilder {
	return &sessionAuthMiddlewareBuilder{
		cache:       cache,
		paths:       baize.Set[string]{},
		sessionType: sessionStatus.SessionTypeAdmin,
	}
}

func NewShopSessionAuthMiddlewareBuilder(cache cache.Cache) *sessionAuthMiddlewareBuilder {
	return &sessionAuthMiddlewareBuilder{
		cache:       cache,
		paths:       baize.Set[string]{},
		sessionType: sessionStatus.SessionTypeShopUser,
	}
}

func (s *sessionAuthMiddlewareBuilder) IgnorePaths(path string) *sessionAuthMiddlewareBuilder {
	s.paths.Add(path)
	return s
}

func (s *sessionAuthMiddlewareBuilder) Build() func(c *gin.Context) {
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
		if !s.paths.Contains(c.Request.RequestURI) {
			_ = manager.RefreshSession(c)
		}
		c.Next()
	}
}
