package datastore

import "context"

type KVStore interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	Close() error
}

type BlobFile struct {
	Data     []byte
	Metadata map[string]string
}

type BlobStore interface {
	Put(ctx context.Context, ID string, data []byte, metadata map[string]string) error
	Get(ctx context.Context, ID string) (*BlobFile, error)
	Del(ctx context.Context, ID string) error
}

type RichStore interface {
	Put(ctx context.Context, in interface{}) error
	Get(ctx context.Context, filter map[string]interface{}) ([]byte, error)
	Update(ctx context.Context, filter map[string]interface{}, in interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
	Query(ctx context.Context, query map[string]interface{}) ([][]byte, error)
}
