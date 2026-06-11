package session

import (
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session/sessionCache"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Manager struct {
	Propagator
	Store
}
type SessionData struct {
	SessionType     string   `json:"session_type"`
	Admin           string   `json:"admin"`
	ShopUser        string   `json:"shop_user"`
	UserId          int64    `json:"user_id"`
	DeptId          int64    `json:"dept_id"`
	Avatar          string   `json:"avatar"`
	UserName        string   `json:"user_name"`
	Role            []int64  `json:"role"`
	RolePerms       []string `json:"role_perms"`
	Permission      []string `json:"permission"`
	IpAddr          string   `json:"ip_addr"`
	LoginTime       int64    `json:"login_time"`
	Os              string   `json:"os"`
	Browser         string   `json:"browser"`
	DataScopeAspect string   `json:"data_scope_aspect"`
}

func NewManager(cache cache.Cache, profile sessionCache.SessionProfile) *Manager {
	return &Manager{
		Propagator: sessionCache.NewPropagator(),
		Store:      sessionCache.NewStore(cache, profile),
	}
}

func NewAdminManager(cache cache.Cache) *Manager {
	return NewManager(cache, sessionCache.AdminProfile)
}

func NewShopManager(cache cache.Cache) *Manager {
	return NewManager(cache, sessionCache.ShopProfile)
}

func (m *Manager) GetSession(ctx *gin.Context) (Session, error) {
	val, ok := ctx.Get(sessionStatus.SessionKey)
	if ok {
		return val.(*sessionCache.Session), nil
	}
	sessId, err := m.Extract(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx, sessId)
	if err != nil {
		return nil, err
	}
	ctx.Set(sessionStatus.SessionKey, sess)
	return sess, nil
}
func (m *Manager) InitSession(ctx *gin.Context, userId int64) (Session, error) {
	sess, err := m.Generate(ctx, userId)
	if err != nil {
		return nil, err
	}
	err = m.Refresh(ctx, sess.Id())
	if err != nil {
		return nil, err
	}
	ctx.Set(sessionStatus.SessionKey, sess)
	return sess, err
}
func (m *Manager) InitSessionWithData(ctx *gin.Context, userId int64, data *SessionData) (Session, error) {
	sess, err := m.InitSession(ctx, userId)
	if err != nil {
		return nil, err
	}
	m.SetSessionData(ctx, sess, data)
	return sess, nil
}
func (m *Manager) SetSessionData(ctx *gin.Context, sess Session, data *SessionData) {
	if data == nil {
		data = &SessionData{}
	}
	if data.UserId == 0 {
		data.UserId = sessionUserId(sess, ctx)
	}
	if data.IpAddr == "" {
		data.IpAddr = ctx.ClientIP()
	}
	if data.LoginTime == 0 {
		data.LoginTime = time.Now().Unix()
	}
	if data.Role == nil {
		data.Role = make([]int64, 0)
	}
	if data.RolePerms == nil {
		data.RolePerms = make([]string, 0)
	}
	if data.Permission == nil {
		data.Permission = make([]string, 0)
	}

	sess.Set(ctx, sessionStatus.SessionType, data.SessionType)
	sess.Set(ctx, sessionStatus.UserId, data.UserId)
	sess.Set(ctx, sessionStatus.DeptId, data.DeptId)
	sess.Set(ctx, sessionStatus.Avatar, data.Avatar)
	sess.Set(ctx, sessionStatus.UserName, data.UserName)
	sess.Set(ctx, sessionStatus.Role, data.Role)
	sess.Set(ctx, sessionStatus.RolePerms, data.RolePerms)
	sess.Set(ctx, sessionStatus.Permission, data.Permission)
	sess.Set(ctx, sessionStatus.IpAddr, data.IpAddr)
	sess.Set(ctx, sessionStatus.LoginTime, data.LoginTime)
	sess.Set(ctx, sessionStatus.Os, data.Os)
	sess.Set(ctx, sessionStatus.Browser, data.Browser)
	sess.Set(ctx, sessionStatus.DataScopeAspect, data.DataScopeAspect)
}
func (m *Manager) InitAppSession(ctx *gin.Context, userId int64) (Session, error) {
	sess, err := m.Generate(ctx, userId)
	if err != nil {
		return nil, err
	}
	ctx.Set(sessionStatus.SessionKey, sess)
	return sess, err
}
func (m *Manager) RemoveSession(ctx *gin.Context) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return
	}
	_ = m.Store.Remove(ctx, sess.Id())

}
func (m *Manager) RefreshSession(ctx *gin.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx, sess.Id())
}

func sessionUserId(sess Session, ctx *gin.Context) int64 {
	userId, err := strconv.ParseInt(sess.Get(ctx, sessionStatus.UserId), 10, 64)
	if err != nil {
		return 0
	}
	return userId
}
