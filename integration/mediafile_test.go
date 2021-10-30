package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Zzocker/book-labs/protos/mediafile"
)

func TestMediafileService(t *testing.T) {
	is := assert.New(t)

	client, err := grpc.Dial(epMediaFile, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer client.Close()
	stub := pb.NewMediaFileServiceClient(client)

	ctx := context.Background()

	var ID string
	var fileCheckSum string
	t.Run("Upload", func(t *testing.T) {
		data, err := ioutil.ReadFile(fullMetalFront)
		if err != nil {
			panic(err)
		}
		hasher := sha256.New()
		hasher.Write(data)
		fileCheckSum = hex.EncodeToString(hasher.Sum(nil))
		id, err := stub.Upload(ctx, &pb.MediaFile{
			Data:      data,
			Extension: filepath.Ext(fullMetalFront),
			Type:      pb.MediaFileType_BOOK,
		})

		is.NoError(err)
		ID = id.GetID()
	})

	t.Run("Get", func(t *testing.T) {
		blob, err := stub.Get(ctx, &pb.MediaFileID{
			ID: ID,
		})
		is.NoError(err)
		hasher := sha256.New()
		hasher.Write(blob.Data)
		is.Equal(fileCheckSum, hex.EncodeToString(hasher.Sum(nil)))
		is.Equal(filepath.Ext(fullMetalFront), blob.Extension)
		is.Equal(pb.MediaFileType_BOOK, blob.Type)
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := stub.Delete(ctx, &pb.MediaFileID{
			ID: ID,
		})
		is.NoError(err)
	})

	t.Run("Get-Deleted", func(t *testing.T) {
		_, err := stub.Get(ctx, &pb.MediaFileID{
			ID: ID,
		})
		is.Error(err)
		is.Equal(codes.NotFound, status.Code(err))
	})
}
