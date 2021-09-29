package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type redisStore struct {
	expiry int64
	pool   *redis.Pool
}

const redisPingTimeoutS = 5

type RedisDumbDatastoreConfig struct {
	ExpiryS  int64
	URL      string
	Username string
	Password string
	Database int
}

func NewRedisDumbDatastore(cfg *RedisDumbDatastoreConfig) DumbDataStore {
	pool := redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(
				ctx,
				"tcp",
				cfg.URL,
				redis.DialUsername(cfg.Username),
				redis.DialPassword(cfg.Password),
				redis.DialDatabase(cfg.Database),
			)
		},
	}
	// ping timeout of 5 second
	ctx, cancel := context.WithTimeout(context.Background(), redisPingTimeoutS*time.Second)
	defer cancel()
	conn, err := pool.DialContext(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, err = conn.Do("PING")
	if err != nil {
		panic(err)
	}

	return &redisStore{
		pool:   &pool,
		expiry: cfg.ExpiryS,
	}
}

func (r *redisStore) Put(ctx context.Context, key string, value []byte) error {
	const op = errors.Op("RedisDumbDatastore.Put")
	conn, err := r.dial(ctx, op)
	if err != nil {
		return err
	}
	defer conn.Close()
	args := []interface{}{key, value}
	if r.expiry > 0 {
		args = append(args, []interface{}{"EX", r.expiry}...)
	}
	_, err = conn.Do("SET", args...)
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed to store : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return nil
}

func (r *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	const op = errors.Op("RedisDumbDatastore.Put")
	conn, err := r.dial(ctx, op)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, errors.E(
			op,
			fmt.Errorf("key not found"),
			errors.CodeNotFound,
			errors.SeverityDebug,
		)
	} else if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to get key : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return raw, nil
}

func (r *redisStore) Delete(ctx context.Context, key string) error {
	const op = errors.Op("RedisDumbDatastore.Delete")
	conn, err := r.dial(ctx, op)
	if err != nil {
		return err
	}
	defer conn.Close()
	code, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed to delete key : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}
	if code != 1 {
		return errors.E(
			op,
			fmt.Errorf("key not found"),
			errors.CodeNotFound,
			errors.SeverityDebug,
		)
	}

	return nil
}

func (r *redisStore) dial(ctx context.Context, op errors.Op) (redis.Conn, error) {
	conn, err := r.pool.DialContext(ctx)
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to dial : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return conn, err
}
