package datastore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/pkg/errors"
)

func TestS3Blob(t *testing.T) {
	is := assert.New(t)

	store := NewS3BlobDatastore(&S3BlobStoreConfig{
		Endpoint:        "http://localhost:4566",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Region:          "ap-southeast-1",
		BucketName:      "dev",
	})
	ctx := context.Background()
	id := "testS3Blob"
	data := "testS3Blob_data"
	extension := "txt"
	mt := map[string]string{
		"Extension": extension,
	}
	t.Run("Put", func(t *testing.T) {
		err := store.Put(ctx, id, []byte(data), mt)
		is.NoError(err)
	})

	t.Run("Get", func(t *testing.T) {
		bfile, err := store.Get(ctx, id)
		is.NoError(err)
		is.Equal(data, string(bfile.Data))
		is.Equal(extension, bfile.Metadata["Extension"])
	})

	t.Run("Del", func(t *testing.T) {
		err := store.Del(ctx, id)
		is.NoError(err)
	})

	t.Run("Get-Deleted", func(t *testing.T) {
		_, err := store.Get(ctx, id)
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})
}
