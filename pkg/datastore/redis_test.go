package datastore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/pkg/errors"
)

func TestRedisStore(t *testing.T) {
	is := assert.New(t)

	dumbDS := NewRedisDumbDatastore(&RedisDumbDatastoreConfig{
		ExpiryS:  int64(time.Minute.Seconds()), // sec
		URL:      "localhost:6379",
		Database: 5,
	})
	const key = "redis-test-key"
	const value = "redis-test-value"
	t.Run("Put", func(t *testing.T) {
		err := dumbDS.Put(context.Background(), key, []byte(value))
		is.NoError(err)
	})

	t.Run("Get", func(t *testing.T) {
		raw, err := dumbDS.Get(context.Background(), key)
		is.NoError(err)
		is.Equal(value, string(raw))
	})

	t.Run("Delete", func(t *testing.T) {
		err := dumbDS.Delete(context.Background(), key)
		is.NoError(err)
	})

	t.Run("NotFound-Get", func(t *testing.T) {
		raw, err := dumbDS.Get(context.Background(), key)
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
		is.Nil(raw)
	})

	t.Run("NotFound-Delete", func(t *testing.T) {
		err := dumbDS.Delete(context.Background(), key)
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})
}
