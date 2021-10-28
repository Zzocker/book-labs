package datastore

import (
	"context"
	"fmt"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type kvMock struct {
	backend map[string][]byte
}

func NewMockKVStore() KVStore {
	return &kvMock{backend: make(map[string][]byte)}
}

func (k *kvMock) Set(ctx context.Context, key string, value []byte) error {
	k.backend[key] = value

	return nil
}

func (k *kvMock) Get(ctx context.Context, key string) ([]byte, error) {
	const op = errors.Op("MockKV.Get")
	raw, ok := k.backend[key]
	if !ok {
		return nil, errors.E(op, fmt.Errorf("not-found"), errors.CodeNotFound)
	}

	return raw, nil
}

func (k *kvMock) Del(ctx context.Context, key string) error {
	const op = errors.Op("MockKV.Del")
	if _, ok := k.backend[key]; !ok {
		return errors.E(op, fmt.Errorf("not-found"), errors.CodeNotFound)
	}
	delete(k.backend, key)

	return nil
}

func (k *kvMock) Close() error {
	return nil
}
