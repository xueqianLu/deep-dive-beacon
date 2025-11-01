package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var releaseScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
`)

func Init(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

// DistributedLock represents a distributed lock.
type DistributedLock struct {
	client *redis.Client
	key    string
	value  string
}

// NewDistributedLock creates a new DistributedLock.
func NewDistributedLock(client *redis.Client, key string) *DistributedLock {
	// Generate a random value for this lock instance to prevent other instances from unlocking it.
	randomBytes := make([]byte, 16)
	// crypto/rand.Read is a cryptographically secure random number generator.
	// It's better than math/rand for this purpose.
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to a less random value if crypto/rand fails
		return &DistributedLock{
			client: client,
			key:    key,
			value:  fmt.Sprintf("%d", time.Now().UnixNano()),
		}
	}
	value := hex.EncodeToString(randomBytes)

	return &DistributedLock{
		client: client,
		key:    key,
		value:  value,
	}
}

// Acquire tries to acquire the lock with a timeout.
func (l *DistributedLock) Acquire(ctx context.Context, expiration time.Duration) (bool, error) {
	ok, err := l.client.SetNX(ctx, l.key, l.value, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// BlockingAcquire tries to acquire the lock with a timeout, blocking and retrying until it succeeds or the context is cancelled.
func (l *DistributedLock) BlockingAcquire(ctx context.Context, expiration time.Duration, waitTimeout time.Duration) (bool, error) {
	retryTicker := time.NewTicker(100 * time.Millisecond) // Retry every 100ms
	defer retryTicker.Stop()

	ctx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-retryTicker.C:
			ok, err := l.Acquire(ctx, expiration)
			if err != nil {
				// You might want to log this error but continue retrying
				continue
			}
			if ok {
				return true, nil
			}
			// Continue loop to retry
		}
	}
}

// Release releases the lock. It uses a Lua script to ensure atomicity.
func (l *DistributedLock) Release(ctx context.Context) error {
	// Use a Lua script to make the release operation atomic.
	// This prevents a client from releasing a lock that was acquired by another client.
	res, err := releaseScript.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return errors.New("lock release failed: lock not held or value mismatch")
	}

	return nil
}
