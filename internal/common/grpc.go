package common

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
)

var reqID struct{} //nolint:gochecknoglobals //key for context

func LoggerServerInteroceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "could't parse incoming context")
	}
	id := md.Get("request_id")[0]
	lg := logger.WithFields(map[string]interface{}{
		"request_id": id,
	})
	lg.Infof("method = %s", info.FullMethod)

	ctx = context.WithValue(ctx, reqID, lg)

	res, err := handler(ctx, req)
	if err != nil {
		errors.Set(err, errors.ReqID(id))
		logger.SystemErr(err)

		return nil, status.Error(codes.Code(errors.ErrCode(err)), err.Error())
	}

	lg.Infof("SUCCESS")

	return res, nil
}

func GetLoggerWithReqID(ctx context.Context) logger.Logger {
	return ctx.Value(reqID).(logger.Logger)
}
