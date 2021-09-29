package datastore

import (
	"context"
)

type File struct {
	ID   string
	Data []byte
}

type FileDataStore interface {
	// Put : insert a file into file server
	Put(ctx context.Context, fl *File) error
	// Get : fetches a file from fileserver
	Get(ctx context.Context, id string) (*File, error)
	// Delete : a file from the fileserver
	Delete(ctx context.Context, id string) error
}

// DumbDataStore : represents a datastore which doesn't support query feature
// this type datastore don't care about value of the key ,only key matters
// eg redis, etcd.
type DumbDataStore interface {
	Put(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

// SmartDS : represnets a datastore which support query feature
// eg : mongo.
type SmartDataStore interface {
	Store(ctx context.Context, in interface{}) error
	Get(ctx context.Context, filter map[string]interface{}) ([]byte, error)
	Update(ctx context.Context, filter map[string]interface{}, in interface{}) error
	UpdateMatching(ctx context.Context, query map[string]interface{}, in interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
	Query(ctx context.Context, sortingKey string, query map[string]interface{}, pageNumber, perPage int64) ([][]byte, error)
	DeleteMatching(ctx context.Context, query map[string]interface{}) error
}
