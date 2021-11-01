package datastore

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type redisKVStore struct {
	pool   *redis.Pool
	expiry int
}

type RedisKVConfig struct {
	URL string
	// if <=0, expiry is ignored
	ExpiryS  int
	Username string
	Password string
	Database int
}

func NewRedisKVStore(cfg *RedisKVConfig) KVStore {
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

	return &redisKVStore{
		pool:   &pool,
		expiry: cfg.ExpiryS,
	}
}

func (r *redisKVStore) Set(ctx context.Context, key string, value []byte) error {
	const op = errors.Op("RedisKV.Set")
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
		return errors.E(op, fmt.Errorf("failed to set: %w", err), errors.CodeInternal)
	}

	return nil
}

func (r *redisKVStore) Get(ctx context.Context, key string) ([]byte, error) {
	const op = errors.Op("RedisKV.Get")
	conn, err := r.dial(ctx, op)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, errors.E(op, fmt.Errorf("key not found: %w", err), errors.CodeNotFound)
	} else if err != nil {
		return nil, err
	}

	return raw, nil
}

func (r *redisKVStore) Del(ctx context.Context, key string) error {
	const op = errors.Op("RedisKV.Del")
	conn, err := r.dial(ctx, op)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		return errors.E(op, fmt.Errorf("failed to delete: %w", err), errors.CodeInternal)
	}

	return nil
}

func (r *redisKVStore) Close() error {
	return r.pool.Close()
}

func (r *redisKVStore) dial(ctx context.Context, op errors.Op) (redis.Conn, error) {
	conn, err := r.pool.DialContext(ctx)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("failed to dial: %w", err), errors.CodeInternal)
	}

	return conn, err
}
