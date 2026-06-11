package sessionCache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/converts"
	"nova-factory-server/app/utils/stringUtils"

	"github.com/spf13/viper"
	"gopkg.in/errgo.v2/errors"
)

type SessionPrefix string

type SessionProfile struct {
	Name        string
	Prefix      SessionPrefix
	SessionType string
}

const (
	AdminSessionPrefix SessionPrefix = "session:admin"
	ShopSessionPrefix  SessionPrefix = "session:shop"
)

var (
	AdminProfile = SessionProfile{
		Name:        "admin",
		Prefix:      AdminSessionPrefix,
		SessionType: sessionStatus.SessionTypeAdmin,
	}
	ShopProfile = SessionProfile{
		Name:        "shop",
		Prefix:      ShopSessionPrefix,
		SessionType: sessionStatus.SessionTypeShopUser,
	}
	ErrSessionNotFound = errors.New("session id not found")
)

type Store struct {
	expiration time.Duration
	cache      cache.Cache
	profile    SessionProfile
}

func NewStore(cache cache.Cache, profile SessionProfile) *Store {
	return &Store{
		expiration: time.Duration(viper.GetInt("token.expire_time")) * time.Minute,
		cache:      cache,
		profile:    normalizeProfile(profile),
	}
}

func (s *Store) Profile() SessionProfile {
	return s.profile
}

func (s *Store) Generate(ctx context.Context, userId int64) (*Session, error) {
	sId := sessionId(userId)
	s.cache.HSet(ctx, s.redisKey(sId), sessionStatus.UserId, userId)
	return NewSession(sId, s.cache, s.profile), nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	ok := s.cache.Expire(ctx, s.redisKey(id), s.expiration)
	if !ok {
		return ErrSessionNotFound
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.cache.Del(ctx, s.redisKey(id))
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (*Session, error) {
	cnt := s.cache.Exists(ctx, s.redisKey(id))
	if cnt != 1 {
		return nil, ErrSessionNotFound
	}
	return NewSession(id, s.cache, s.profile), nil
}

func (s *Store) ScanSessions(ctx context.Context) []*Session {
	var cursor uint64
	sessions := make([]*Session, 0, 16)
	for {
		keys, newCursor := s.cache.Scan(ctx, cursor, s.scanPattern(), 10)
		for _, key := range keys {
			id := s.trimRedisKey(key)
			if id != "" {
				sessions = append(sessions, NewSession(id, s.cache, s.profile))
			}
		}
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	return sessions
}

func (s *Store) redisKey(id string) string {
	return redisKey(s.profile, id)
}

func (s *Store) scanPattern() string {
	return fmt.Sprintf("%s:*", s.profile.Prefix)
}

func (s *Store) trimRedisKey(key string) string {
	return strings.TrimPrefix(key, fmt.Sprintf("%s:", s.profile.Prefix))
}

func NewSession(id string, cache cache.Cache, profiles ...SessionProfile) *Session {
	return &Session{
		id:      id,
		values:  make(map[string]string),
		cache:   cache,
		profile: resolveProfile(profiles...),
	}
}

type Session struct {
	id      string
	values  map[string]string
	cache   cache.Cache
	profile SessionProfile
}

func (s *Session) Get(ctx context.Context, key string) string {
	val := s.values[key]
	if val != "" {
		return val
	}
	result := s.cache.HGet(ctx, s.redisKey(), key)
	s.values[key] = result
	return result
}

func (s *Session) Set(ctx context.Context, key string, val any) {
	gs := converts.String(val)
	s.values[key] = gs
	s.cache.JudgmentAndHSet(ctx, s.redisKey(), key, gs)
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) Profile() SessionProfile {
	return s.profile
}

func (s *Session) redisKey() string {
	return redisKey(s.profile, s.id)
}

func sessionId(userId int64) string {
	return fmt.Sprintf("%d:%s", userId, stringUtils.GetUUID())
}

func redisKey(profile SessionProfile, id string) string {
	return fmt.Sprintf("%s:%s", normalizeProfile(profile).Prefix, id)
}

func resolveProfile(profiles ...SessionProfile) SessionProfile {
	if len(profiles) == 0 {
		return AdminProfile
	}
	return normalizeProfile(profiles[0])
}

func normalizeProfile(profile SessionProfile) SessionProfile {
	if profile.Prefix == "" {
		return AdminProfile
	}
	return profile
}
