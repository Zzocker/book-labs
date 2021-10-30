package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/Zzocker/book-labs/mediafile/config"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/logger"
	"github.com/Zzocker/book-labs/pkg/utils"
	pb "github.com/Zzocker/book-labs/protos/mediafile"
)

func run(cfg *config.MediaFileServiceConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}
	logger.Setup(cfg.Log.Level, cfg.Log.ServiceName, cfg.Log.ServiceVersion)
	mdService := mediaFileRPC{
		store: datastore.NewS3BlobDatastore(&datastore.S3BlobStoreConfig{
			Endpoint:        cfg.S3Config.Endpoint,
			AccessKeyID:     cfg.S3Config.AccessKeyID,
			SecretAccessKey: cfg.S3Config.SecretAccessKey,
			SessionToken:    cfg.S3Config.SessionToken,
			Region:          cfg.S3Config.Region,
			BucketName:      cfg.S3Config.BucketName,
		}),
	}
	gRPCSrv := grpc.NewServer(
		grpc.UnaryInterceptor(utils.GRPCServerLoggerUnarayInteroceptor),
	)

	pb.RegisterMediaFileServiceServer(gRPCSrv, &mdService)

	serveErrChan := make(chan error)
	go func() {
		logger.Infof("Listing on Port: %s", cfg.Port)
		serveErrChan <- gRPCSrv.Serve(lis)
		close(serveErrChan)
	}()

	sysSignal := make(chan os.Signal, 1)
	signal.Notify(sysSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case sys := <-sysSignal:
		logger.Infof("Run sys signal: %v", sys.String())
	case err := <-serveErrChan:
		logger.Info("Server error: %v", err)
	}
}
