package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	Prefix   string
}

func (cfg Config) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func NewClient(cfg Config) Redis {
	return &standaloneRedis{
		prefix: cfg.Prefix,
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr(),
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}
}

type Redis interface {
	AppendPrefix(key string) string
	AppendPrefixSlice(keys []string) []string
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	Publish(ctx context.Context, channel string, message interface{}) (int64, error)
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	Ping(ctx context.Context) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	RPopCount(ctx context.Context, key string, count int) ([]string, error)
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	Pipeline() redis.Pipeliner

	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SPopN(ctx context.Context, key string, count int64) ([]string, error)
	Close() error
}
