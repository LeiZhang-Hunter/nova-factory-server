package sessionCache

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/errgo.v2/errors"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/converts"
	"nova-factory-server/app/utils/stringUtils"
	"time"
)

var (
	SessionKey         = `session_key`
	ErrSessionNotFound = errors.New("session:id 对应的session不存在")
)

type Store struct {
	expiration time.Duration
	cache      cache.Cache
}

func NewStore(cache cache.Cache) *Store {
	return &Store{
		expiration: time.Duration(viper.GetInt("token.expire_time")) * time.Minute,
		cache:      cache,
	}
}

func (s *Store) Generate(ctx context.Context, userId int64) (*Session, error) {
	sId := sessionId(userId)
	s.cache.HSet(ctx, redisKey(sId), sessionStatus.UserId, userId)
	return NewSession(sId, s.cache), nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	ok := s.cache.Expire(ctx, redisKey(id), s.expiration)

	if !ok {
		return ErrSessionNotFound
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.cache.Del(ctx, redisKey(id))
	return nil

}

func (s *Store) Get(ctx context.Context, id string) (*Session, error) {
	cnt := s.cache.Exists(ctx, redisKey(id))
	if cnt != 1 {
		return nil, ErrSessionNotFound
	}
	return NewSession(id, s.cache), nil
}

func NewSession(id string, cache cache.Cache) *Session {
	return &Session{
		id:     id,
		values: make(map[string]string),
		cache:  cache,
	}
}

type Session struct {
	id     string
	values map[string]string
	cache  cache.Cache
}

func (s *Session) Get(ctx context.Context, key string) string {
	val := s.values[key]
	if val != "" {
		return val
	}
	result := s.cache.HGet(ctx, redisKey(s.id), key)
	s.values[key] = result
	return result

}

func (s *Session) Set(ctx context.Context, key string, val any) {
	gs := converts.String(val)
	s.values[key] = gs
	s.cache.JudgmentAndHSet(ctx, redisKey(s.id), key, gs)
}

func (s *Session) Id() string {
	return s.id
}

func sessionId(userId int64) string {
	return fmt.Sprintf("%d:%s", userId, stringUtils.GetUUID())
}

func redisKey(id string) string {
	return fmt.Sprintf("%s:%s", SessionKey, id)
}
