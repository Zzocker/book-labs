package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	"github.com/Zzocker/book-labs/pkg/logger"
)

func loggerWithReqID(c *gin.Context) logger.Logger {
	return logger.WithFields(map[string]interface{}{
		"request_id": c.GetString(reqIDKey),
	})
}

func gRPCCtxWithReqID(ctx *gin.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD{
		"request_id": []string{ctx.GetString(reqIDKey)},
	})
}
