package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/constant/dataScopeAspect"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/baizeContext"
	keyStore "nova-factory-server/app/utils/store/key"
	permissionStore "nova-factory-server/app/utils/store/permissions"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
		if builder.authenticateByAPIKey(c) {
			return
		}

		// Use the current route domain to choose admin or shop session storage/profile.
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

// authenticateByAPIKey 将外部 MCP API Key 调用接入现有业务 session 鉴权流程。
func (builder *SessionAuthBuilder) authenticateByAPIKey(c *gin.Context) bool {
	apiKey := strings.TrimSpace(c.GetHeader(agent.AuthorizationApiKey))
	if apiKey == "" {
		return false
	}
	if err := builder.initAPIKeySession(c, apiKey); err != nil {
		c.Set(sessionStatus.MsgKey, err.Error())
		builder.handleFailure(c, http.StatusUnauthorized)
		return true
	}
	c.Next()
	return true
}

func (builder *SessionAuthBuilder) initAPIKeySession(c *gin.Context, apiKey string) error {
	info, err := keyStore.GetStore().GetInfo(c, apiKey)
	if err != nil {
		return err
	}
	if info == nil || info.UserID == 0 {
		return errors.New("api key无效")
	}

	// MCP 适配层会注入解析后的工具名，直接调用工具时也要校验授权。
	toolName := strings.TrimSpace(c.GetHeader(agent.MCPToolNameHeader))
	if toolName == "" {
		return errors.New("缺少MCP工具标识")
	}
	if !containsString(info.Tools, toolName) {
		return fmt.Errorf("api key无权调用MCP工具: %s", toolName)
	}

	manager := builder.manager()
	if sess, ok := builder.cachedAPIKeySession(c, manager, apiKey); ok {
		c.Set(sessionStatus.SessionKey, sess)
		_ = manager.Refresh(c, sess.Id())
		builder.refreshAPIKeySessionMapping(c, apiKey, sess.Id())
		return nil
	}

	permissions := permissionStore.GetStore().GetPermission(c, info.UserID)
	sess, err := currentOrNewSession(c, manager, info.UserID)
	if err != nil {
		return err
	}
	builder.fillAPIKeySessionData(c, sess, info, permissions)
	builder.refreshAPIKeySessionMapping(c, apiKey, sess.Id())
	return nil
}

// cachedAPIKeySession 复用短生命周期业务 session，避免高频 MCP 调用反复创建 session。
func (builder *SessionAuthBuilder) cachedAPIKeySession(c *gin.Context, manager *session.Manager, apiKey string) (session.Session, bool) {
	sessionID, err := builder.cache.Get(c, builder.apiKeySessionMappingKey(apiKey))
	if err != nil || sessionID == "" {
		return nil, false
	}
	sess, err := manager.Get(c, sessionID)
	if err != nil {
		return nil, false
	}
	if sess.Get(c, sessionStatus.SessionType) != builder.expectedSessionType() {
		return nil, false
	}
	return sess, true
}

// fillAPIKeySessionData 写入后续 handler 通过 baizeContext 读取的最小必要字段。
func (builder *SessionAuthBuilder) fillAPIKeySessionData(c *gin.Context, sess session.Session, info *keyStore.Info, permissions []string) {
	sess.Set(c, sessionStatus.SessionType, builder.expectedSessionType())
	sess.Set(c, sessionStatus.UserId, info.UserID)
	sess.Set(c, sessionStatus.DeptId, info.DeptID)
	sess.Set(c, sessionStatus.Permission, permissions)
	sess.Set(c, sessionStatus.RolePerms, []string{})
	sess.Set(c, sessionStatus.IpAddr, c.ClientIP())
	sess.Set(c, sessionStatus.LoginTime, time.Now().Unix())
	if sess.Get(c, sessionStatus.UserName) == "" {
		sess.Set(c, sessionStatus.UserName, fmt.Sprintf("api-key:%s:%d", builder.expectedSessionType(), info.UserID))
	}
	if sess.Get(c, sessionStatus.Avatar) == "" {
		sess.Set(c, sessionStatus.Avatar, "")
	}
	if sess.Get(c, sessionStatus.Role) == "" {
		sess.Set(c, sessionStatus.Role, []int64{})
	}
	if sess.Get(c, sessionStatus.DataScopeAspect) == "" {
		sess.Set(c, sessionStatus.DataScopeAspect, apiKeyDataScope(c, sess, info))
	}
}

func (builder *SessionAuthBuilder) refreshAPIKeySessionMapping(c *gin.Context, apiKey string, sessionID string) {
	expire := time.Duration(viper.GetInt("token.expire_time")) * time.Minute
	builder.cache.Set(c, builder.apiKeySessionMappingKey(apiKey), sessionID, expire)
}

func currentOrNewSession(c *gin.Context, manager *session.Manager, userID int64) (session.Session, error) {
	if value, ok := c.Get(sessionStatus.SessionKey); ok {
		if sess, ok := value.(session.Session); ok {
			return sess, nil
		}
	}
	return manager.InitSession(c, userID)
}

// apiKeySessionMappingKey 带上 session 类型，避免 admin/shop 共用同一个缓存 session。
func (builder *SessionAuthBuilder) apiKeySessionMappingKey(apiKey string) string {
	sum := sha256.Sum256([]byte(apiKey))
	return fmt.Sprintf("mcp:api_key_session:%s:%s", builder.expectedSessionType(), hex.EncodeToString(sum[:]))
}

func apiKeyDataScope(c *gin.Context, sess session.Session, info *keyStore.Info) string {
	if sessionHasRole(c, sess, 1) {
		return dataScopeAspect.DataScopeAll
	}
	if info.DeptID != 0 {
		return dataScopeAspect.DataScopeDept
	}
	return dataScopeAspect.NoDataScope
}

func sessionHasRole(c *gin.Context, sess session.Session, roleID int64) bool {
	rolesText := sess.Get(c, sessionStatus.Role)
	if rolesText == "" {
		return false
	}
	roles := make([]int64, 0)
	if err := json.Unmarshal([]byte(rolesText), &roles); err != nil {
		return false
	}
	for _, item := range roles {
		if item == roleID {
			return true
		}
	}
	return false
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
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
