package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type standaloneRedis struct {
	client *redis.Client
	prefix string
}

func (r *standaloneRedis) AppendPrefix(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}
func (r *standaloneRedis) AppendPrefixSlice(keys []string) []string {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.AppendPrefix(key)
	}
	return prefixedKeys
}

func (r *standaloneRedis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.AppendPrefix(key)).Result()
}

func (r *standaloneRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client.Set(ctx, r.AppendPrefix(key), value, expiration).Result()
}

func (r *standaloneRedis) Ping(ctx context.Context) (string, error) {
	return r.client.Ping(ctx).Result()
}

func (r *standaloneRedis) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	return r.client.Publish(ctx, r.AppendPrefix(channel), message).Result()
}

func (r *standaloneRedis) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.client.Subscribe(ctx, r.AppendPrefixSlice(channels)...)
}

func (r *standaloneRedis) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, r.AppendPrefix(key)).Result()
}

func (r *standaloneRedis) RPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return r.client.RPopCount(ctx, r.AppendPrefix(key), count).Result()
}

func (r *standaloneRedis) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.LPush(ctx, r.AppendPrefix(key), values...).Result()
}

func (r *standaloneRedis) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

func (r *standaloneRedis) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(ctx, r.AppendPrefix(key), members...).Result()
}

func (r *standaloneRedis) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.SRem(ctx, r.AppendPrefix(key), members...).Result()
}

func (r *standaloneRedis) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return r.client.SPopN(ctx, r.AppendPrefix(key), count).Result()
}

func (r *standaloneRedis) Close() error {
	return r.client.Close()
}
