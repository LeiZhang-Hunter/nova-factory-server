package session

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/middlewares/session/sessionCache"
)

type Store interface {
	Generate(ctx context.Context, userId int64) (*sessionCache.Session, error)
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*sessionCache.Session, error)
}

type Session interface {
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, val any)
	Id() string
}

type Propagator interface {
	Extract(c *gin.Context) (string, error)
}
