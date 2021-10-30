package utils

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
)

func GRPCServerLoggerUnarayInteroceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger.Infof("method = %s", info.FullMethod)

	res, err := handler(ctx, req)
	if err != nil {
		return nil, status.Error(
			codes.Code(errors.ErrCode(err)),
			err.Error(),
		)
	}
	logger.Info("SUCCESS")

	return res, nil
}
