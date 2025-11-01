package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestDistributedLock_Acquire(t *testing.T) {
	client, mock := redismock.NewClientMock()
	lockKey := "test-lock"

	t.Run("acquire new lock", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		expiration := 10 * time.Second

		mock.ExpectSetNX(lock.key, lock.value, expiration).SetVal(true)

		ok, err := lock.Acquire(context.Background(), expiration)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("acquire existing lock", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		expiration := 10 * time.Second

		mock.ExpectSetNX(lock.key, lock.value, expiration).SetVal(false)

		ok, err := lock.Acquire(context.Background(), expiration)
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("acquire with redis error", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		expiration := 10 * time.Second
		redisErr := errors.New("redis error")

		mock.ExpectSetNX(lock.key, lock.value, expiration).SetErr(redisErr)

		ok, err := lock.Acquire(context.Background(), expiration)
		assert.Error(t, err)
		assert.Equal(t, redisErr, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDistributedLock_BlockingAcquire(t *testing.T) {
	client, mock := redismock.NewClientMock()
	lockKey := "test-blocking-lock"

	t.Run("acquire immediately", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		expiration := 10 * time.Second
		waitTimeout := 1 * time.Second

		mock.ExpectSetNX(lock.key, lock.value, expiration).SetVal(true)

		ok, err := lock.BlockingAcquire(context.Background(), expiration, waitTimeout)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("acquire after retry", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		expiration := 10 * time.Second
		waitTimeout := 500 * time.Millisecond

		// First attempt fails
		mock.ExpectSetNX(lock.key, lock.value, expiration).SetVal(false)
		// Second attempt succeeds
		mock.ExpectSetNX(lock.key, lock.value, expiration).SetVal(true)

		ok, err := lock.BlockingAcquire(context.Background(), expiration, waitTimeout)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("timeout while waiting", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		expiration := 10 * time.Second
		waitTimeout := 250 * time.Millisecond // Shorter than one retry cycle + processing

		// Always fail to acquire
		mock.ExpectSetNX(lock.key, lock.value, expiration).SetVal(false)

		ok, err := lock.BlockingAcquire(context.Background(), expiration, waitTimeout)
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDistributedLock_Release(t *testing.T) {
	client, mock := redismock.NewClientMock()
	lockKey := "test-release-lock"

	t.Run("release held lock", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)

		mock.ExpectEvalSha(releaseScript.Hash(), []string{lock.key}, lock.value).SetVal(int64(1))

		err := lock.Release(context.Background())
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("release lock not held", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)

		mock.ExpectEvalSha(releaseScript.Hash(), []string{lock.key}, lock.value).SetVal(int64(0))

		err := lock.Release(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "lock release failed: lock not held or value mismatch", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("release with redis error", func(t *testing.T) {
		lock := NewDistributedLock(client, lockKey)
		redisErr := errors.New("redis error")

		mock.ExpectEvalSha(releaseScript.Hash(), []string{lock.key}, lock.value).SetErr(redisErr)

		err := lock.Release(context.Background())
		assert.Error(t, err)
		assert.Equal(t, redisErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewDistributedLock(t *testing.T) {
	client, _ := redismock.NewClientMock()
	lockKey := "test-new-lock"

	lock1 := NewDistributedLock(client, lockKey)
	lock2 := NewDistributedLock(client, lockKey)

	assert.Equal(t, lockKey, lock1.key)
	assert.Equal(t, lockKey, lock2.key)
	assert.NotEqual(t, lock1.value, lock2.value, "Values for different lock instances should be unique")
}
