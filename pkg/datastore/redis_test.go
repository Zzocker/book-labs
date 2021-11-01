package datastore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/pkg/errors"
)

func TestRedisKV(t *testing.T) {
	is := assert.New(t)

	kv := NewRedisKVStore(&RedisKVConfig{
		URL:      "localhost:6379",
		ExpiryS:  int(time.Minute.Seconds()), // 60 second
		Database: 5,
	})
	ctx := context.Background()
	const key = "redisKV_key"
	const value = "redisKV_value"

	t.Run("Set", func(t *testing.T) {
		err := kv.Set(ctx, key, []byte(value))
		is.NoError(err)
	})

	t.Run("Get", func(t *testing.T) {
		raw, err := kv.Get(ctx, key)
		is.NoError(err)
		is.Equal(value, string(raw))
	})

	t.Run("Delete", func(t *testing.T) {
		err := kv.Del(ctx, key)
		is.NoError(err)
	})

	t.Run("Get-Deleted", func(t *testing.T) {
		_, err := kv.Get(ctx, key)
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Close", func(t *testing.T) {
		err := kv.Close()
		is.NoError(err)
	})
}
