package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
	"github.com/Zzocker/book-labs/protos/common"
	pb "github.com/Zzocker/book-labs/protos/mediafile"
)

type mediaFileRPC struct {
	pb.UnimplementedMediaFileServiceServer
	store datastore.BlobStore
}

const (
	mdExtension string = "Extension"
	mdType      string = "Type"
)

func (m *mediaFileRPC) Upload(ctx context.Context, fi *pb.MediaFile) (*pb.MediaFileID, error) {
	const op = errors.Op("Service.Upload")

	if len(fi.Data) == 0 {
		return nil, errors.E(op, fmt.Errorf("empty data"), errors.CodeInvalidArgument)
	}

	if fi.Type == pb.MediaFileType_UNKNOWN {
		return nil, errors.E(op, fmt.Errorf("unknown data type"), errors.CodeInvalidArgument)
	}

	id := uuid.New().String()
	logger.Debugf("upload id = %s, extension = %s, type = %s", id, fi.GetExtension(), fi.GetType().String())
	err := m.store.Put(ctx, id, fi.Data, map[string]string{
		mdExtension: fi.GetExtension(),
		mdType:      fi.GetType().String(),
	})
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &pb.MediaFileID{
		ID: id,
	}, nil
}

func (m *mediaFileRPC) Get(ctx context.Context, id *pb.MediaFileID) (*pb.MediaFile, error) {
	const op = errors.Op("Service.Get")

	logger.Debugf("get id = %s", id.GetID())
	blob, err := m.store.Get(ctx, id.GetID())
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &pb.MediaFile{
		ID:        id.GetID(),
		Data:      blob.Data,
		Extension: blob.Metadata[mdExtension],
		Type:      pb.MediaFileType(pb.MediaFileType_value[blob.Metadata[mdType]]),
	}, nil
}

func (m *mediaFileRPC) Delete(ctx context.Context, id *pb.MediaFileID) (*common.Empty, error) {
	const op = errors.Op("Service.Delete")

	err := m.store.Del(ctx, id.GetID())
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &common.Empty{}, nil
}
