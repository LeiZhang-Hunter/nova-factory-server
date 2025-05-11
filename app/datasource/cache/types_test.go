package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	redis2 "nova-factory-server/app/datasource/cache/redis"
	"testing"
)

func TestRedisCache(t *testing.T) {
	ctx := context.Background()
	redisDb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		DB:   1,
	})
	err := redisDb.Set(ctx, "1", "1", 0).Err()
	fmt.Println(err)
	rc := redis2.NewRedisCache(redisDb)
	rc.Set(ctx, "2", "2", 0)

}
