package cache

import (
	"context"
	"encoding/json"
	"time"

	"golang.org/x/sync/singleflight"

	"godan/internal/pkg/redis"
)

var g singleflight.Group

// GetOrSet: cache-aside with singleflight protection.
// If key exists in Redis, unmarshal to dest and return true.
// If not, call loader(), store result, and return false (caller handles).
func Get(ctx context.Context, key string, dest interface{}) bool {
	val, err := redis.Get(ctx, key)
	if err != nil || val == "" {
		return false
	}
	return json.Unmarshal([]byte(val), dest) == nil
}

// Set stores value in Redis with TTL.
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, _ := json.Marshal(value)
	return redis.Set(ctx, key, string(data), ttl)
}

// Del removes a key.
func Del(ctx context.Context, key string) {
	redis.Del(ctx, key)
}

// GetOrLoad: singleflight + Redis cache.
// Concurrent calls for the same key will only trigger loader once.
func GetOrLoad(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	// try cache first
	if Get(ctx, key, dest) {
		return nil
	}

	// singleflight to prevent cache stampede
	v, err, _ := g.Do(key, func() (interface{}, error) {
		// double-check cache
		if Get(ctx, key, dest) {
			return nil, nil
		}

		result, err := loader()
		if err != nil {
			return nil, err
		}

		Set(ctx, key, result, ttl)
		return result, nil
	})

	if err != nil {
		return err
	}

	if v != nil {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, dest)
	}

	return nil
}

// WithTTL returns a random TTL in [base, base+spread] to avoid cache avalanche.
func WithTTL(base time.Duration, spread time.Duration) time.Duration {
	return base + time.Duration(time.Now().UnixNano()%int64(spread))
}
