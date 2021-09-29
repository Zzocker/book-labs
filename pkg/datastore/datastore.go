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
