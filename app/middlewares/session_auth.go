package middlewares

import (
	"net/http"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type AuthDomain int

const (
	AdminDomain AuthDomain = iota
	ShopDomain
)

type AuthMode int

const (
	RequiredAuth AuthMode = iota
	OptionalAuth
)

type AuthTransport int

const (
	HTTPTransport AuthTransport = iota
	WebSocketTransport
)

type AuthOption func(*SessionAuthBuilder)

type SessionAuthBuilder struct {
	cache       cache.Cache
	domain      AuthDomain
	mode        AuthMode
	transport   AuthTransport
	refresh     bool
	ignorePaths baize.Set[string]
}

func NewSessionAuth(cache cache.Cache, domain AuthDomain, opts ...AuthOption) *SessionAuthBuilder {
	builder := &SessionAuthBuilder{
		cache:       cache,
		domain:      domain,
		mode:        RequiredAuth,
		transport:   HTTPTransport,
		ignorePaths: baize.Set[string]{},
	}
	for _, opt := range opts {
		opt(builder)
	}
	return builder
}

func Optional() AuthOption {
	return func(builder *SessionAuthBuilder) {
		builder.mode = OptionalAuth
	}
}

func WebSocket() AuthOption {
	return func(builder *SessionAuthBuilder) {
		builder.transport = WebSocketTransport
	}
}

func Refresh() AuthOption {
	return func(builder *SessionAuthBuilder) {
		builder.refresh = true
	}
}

func IgnorePaths(paths ...string) AuthOption {
	return func(builder *SessionAuthBuilder) {
		builder.IgnorePaths(paths...)
	}
}

func (builder *SessionAuthBuilder) Optional() *SessionAuthBuilder {
	builder.mode = OptionalAuth
	return builder
}

func (builder *SessionAuthBuilder) ForWebSocket() *SessionAuthBuilder {
	builder.transport = WebSocketTransport
	return builder
}

func (builder *SessionAuthBuilder) WithRefresh() *SessionAuthBuilder {
	builder.refresh = true
	return builder
}

func (builder *SessionAuthBuilder) IgnorePaths(paths ...string) *SessionAuthBuilder {
	for _, path := range paths {
		builder.ignorePaths.Add(path)
	}
	return builder
}

func (builder *SessionAuthBuilder) BuildForWebSocket() gin.HandlerFunc {
	builder.transport = WebSocketTransport
	return builder.Build()
}

func (builder *SessionAuthBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		manager := builder.manager()

		sessionID, err := manager.Extract(c)
		if err != nil {
			builder.handleFailure(c, http.StatusUnauthorized)
			return
		}

		currentSession, err := manager.Get(c, sessionID)
		if err != nil {
			builder.handleFailure(c, http.StatusUnauthorized)
			return
		}

		if currentSession.Get(c, sessionStatus.SessionType) != builder.expectedSessionType() {
			builder.handleFailure(c, http.StatusForbidden)
			return
		}

		c.Set(sessionStatus.SessionKey, currentSession)
		if builder.shouldRefresh(c) {
			_ = manager.Refresh(c, currentSession.Id())
		}
		c.Next()
	}
}

func (builder *SessionAuthBuilder) manager() *session.Manager {
	if builder.domain == ShopDomain {
		return session.NewShopManager(builder.cache)
	}
	return session.NewAdminManager(builder.cache)
}

func (builder *SessionAuthBuilder) expectedSessionType() string {
	if builder.domain == ShopDomain {
		return sessionStatus.SessionTypeShopUser
	}
	return sessionStatus.SessionTypeAdmin
}

func (builder *SessionAuthBuilder) shouldRefresh(c *gin.Context) bool {
	return builder.refresh && !builder.ignorePaths.Contains(c.Request.RequestURI)
}

func (builder *SessionAuthBuilder) handleFailure(c *gin.Context, statusCode int) {
	if builder.mode == OptionalAuth {
		c.Next()
		return
	}
	if builder.transport == WebSocketTransport {
		c.AbortWithStatus(statusCode)
		return
	}
	baizeContext.InvalidToken(c)
	c.Abort()
}
