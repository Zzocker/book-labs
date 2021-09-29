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
