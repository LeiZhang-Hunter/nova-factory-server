package redis

import (
	"context"
	"errors"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"nova-factory-server/app/datasource/cache/cacheError"
)

var (
	lua = `if redis.call("exists", KEYS[1])
							then 
								return redis.call("hset", KEYS[1], ARGV[1],ARGV[2])
							else
								return -1
							end`
)

const (
	appMaxRetries     = 3
	appRetryBaseDelay = 200 * time.Millisecond
)

// isRetryable 判断错误是否为连接层错误（可重试），redis.Nil 不可重试
func isRetryable(err error) bool {
	if err == nil || errors.Is(err, redis.Nil) {
		return false
	}
	// 网络级别错误可重试
	if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	// 其他 redis 协议错误（如 READONLY、WRONGTYPE）不重试
	// 剩余未知错误一律重试，宁可多试一次也不要直接失败
	return true
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

//func (r *RedisCache) Publish(ctx context.Context, channel string, message interface{}) error {
//	return r.client.Publish(ctx, channel, message).Err()
//}
//
//func (r *RedisCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
//	return r.client.Subscribe(ctx, channels...).Channel()
//}

func (r *RedisCache) Set(ctx context.Context, key string, val string, expiration time.Duration) {
	err := retry(func() error {
		_, e := r.client.Set(ctx, key, val, expiration).Result()
		return e
	})
	if err != nil {
		log.Printf("[redis] Set key=%s failed after retries: %v", key, err)
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	var result string
	err := retry(func() error {
		s, e := r.client.Get(ctx, key).Result()
		if e != nil {
			return e
		}
		result = s
		return nil
	})
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", cacheError.Nil
		}
		log.Printf("[redis] Get key=%s failed after retries: %v", key, err)
		return "", err
	}
	return result, nil
}

func (r *RedisCache) Del(ctx context.Context, keys ...string) {
	err := retry(func() error {
		_, e := r.client.Del(ctx, keys...).Result()
		return e
	})
	if err != nil {
		log.Printf("[redis] Del keys=%v failed after retries: %v", keys, err)
	}
}

func (r *RedisCache) HSet(ctx context.Context, key string, values ...any) {
	err := retry(func() error {
		_, e := r.client.HSet(ctx, key, values...).Result()
		return e
	})
	if err != nil {
		log.Printf("[redis] HSet key=%s failed after retries: %v", key, err)
	}
}

func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) bool {
	var ok bool
	err := retry(func() error {
		res, e := r.client.Expire(ctx, key, expiration).Result()
		if e != nil {
			return e
		}
		ok = res
		return nil
	})
	if err != nil {
		log.Printf("[redis] Expire key=%s failed after retries: %v", key, err)
		return false
	}
	return ok
}

func (r *RedisCache) Exists(ctx context.Context, keys ...string) int64 {
	var cnt int64
	err := retry(func() error {
		res, e := r.client.Exists(ctx, keys...).Result()
		if e != nil {
			return e
		}
		cnt = res
		return nil
	})
	if err != nil {
		log.Printf("[redis] Exists keys=%v failed after retries: %v", keys, err)
		return 0
	}
	return cnt
}

func (r *RedisCache) HGet(ctx context.Context, key, field string) string {
	var result string
	err := retry(func() error {
		s, e := r.client.HGet(ctx, key, field).Result()
		if e != nil {
			return e
		}
		result = s
		return nil
	})
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ""
		}
		log.Printf("[redis] HGet key=%s field=%s failed after retries: %v", key, field, err)
		return ""
	}
	return result
}

func (r *RedisCache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64) {
	var (
		keys []string
		c    uint64
	)
	err := retry(func() error {
		k, nc, e := r.client.Scan(ctx, cursor, match, count).Result()
		if e != nil {
			return e
		}
		keys = k
		c = nc
		return nil
	})
	if err != nil {
		log.Printf("[redis] Scan cursor=%d match=%s failed after retries: %v", cursor, match, err)
		return nil, 0
	}
	return keys, c
}

func (r *RedisCache) JudgmentAndHSet(ctx context.Context, ids, key string, gs any) {
	err := retry(func() error {
		_, e := r.client.Eval(ctx, lua, []string{ids}, key, gs).Int()
		return e
	})
	if err != nil {
		log.Printf("[redis] JudgmentAndHSet ids=%s key=%s failed after retries: %v", ids, key, err)
	}
}

func (r *RedisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) bool {
	var ok bool
	err := retry(func() error {
		res, e := r.client.SetNX(ctx, key, value, expiration).Result()
		if e != nil {
			return e
		}
		ok = res
		return nil
	})
	if err != nil {
		log.Printf("[redis] SetNX key=%s failed after retries: %v", key, err)
		return false
	}
	return ok
}

func (r *RedisCache) Publish(ctx context.Context, channel string, message interface{}) {
	err := retry(func() error {
		return r.client.Publish(ctx, channel, message).Err()
	})
	if err != nil {
		log.Printf("[redis] Publish channel=%s failed after retries: %v", channel, err)
	}
}

func (r *RedisCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.client.Subscribe(ctx, channels...)
}

func (r *RedisCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.client.ZRangeByScore(ctx, key, opt)
}

func (r *RedisCache) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAdd(ctx, key, members...)
}

func (r *RedisCache) ZRem(ctx context.Context, key string, members ...any) *redis.IntCmd {
	return r.client.ZRem(ctx, key, members...)
}

func (r *RedisCache) MGet(ctx context.Context, keys []string) *redis.SliceCmd {
	result := r.client.MGet(ctx, keys...)
	return result
}

// retry 应用层重试，区分 redis.Nil（不重试）和网络连接错误（重试）
func retry(fn func() error) error {
	var lastErr error
	for i := 0; i < appMaxRetries; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if errors.Is(lastErr, redis.Nil) {
			return lastErr // 空值不重试
		}
		if !isRetryable(lastErr) {
			return lastErr // 协议级错误不重试
		}
		time.Sleep(appRetryBaseDelay << i) // 200ms, 400ms, 800ms 指数退避
	}
	return lastErr
}
