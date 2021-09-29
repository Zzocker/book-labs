package datastore

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3(t *testing.T) {
	is := assert.New(t)
	fileStore := NewS3FileDatastore(&S3FileDatastoreConfig{
		// localstack s3
		Endpoint:        "http://localhost:4566",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Region:          "ap-southeast-1",
		BucketName:      "dev",
		// aws --endpoint-url=http://localhost:4566 s3 mb s3://dev
	})

	f := File{
		ID: "testS3.txt",
		Data: []byte(`
Hello This is just for testing`),
	}
	t.Run("Put", func(t *testing.T) {
		err := fileStore.Put(context.Background(), &f)
		is.NoError(err)
	})

	t.Run("Get", func(t *testing.T) {
		fl, err := fileStore.Get(context.Background(), f.ID)
		is.NoError(err)
		is.True(reflect.DeepEqual(fl.Data, f.Data))
	})

	t.Run("Delete", func(t *testing.T) {
		err := fileStore.Delete(context.Background(), f.ID)
		is.NoError(err)
	})
}
