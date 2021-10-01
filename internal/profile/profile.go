package profile

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/Zzocker/book-labs/config"
	"github.com/Zzocker/book-labs/internal/common"
	sgrpc "github.com/Zzocker/book-labs/internal/profile/delivery/grpc"
	"github.com/Zzocker/book-labs/internal/profile/service"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/logger"
	pb "github.com/Zzocker/book-labs/protos/profile"
)

func Run(cfg *config.BookSharing) {
	lis, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", cfg.UserProfile.Port))
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(common.LoggerServerInteroceptor),
	)

	// datastore
	fStore := datastore.NewS3FileDatastore(&datastore.S3FileDatastoreConfig{
		Endpoint:        cfg.S3.Endpoint,
		AccessKeyID:     cfg.S3.AccessKeyID,
		SecretAccessKey: cfg.S3.SecretAccessKey,
		SessionToken:    cfg.S3.SessionToken,
		Region:          cfg.S3.Region,
		BucketName:      cfg.UserProfile.ProfileBucketName,
	})
	sStore := datastore.NewMongoSmartDatastore(&datastore.MongoSmartDatastoreConfig{
		Username:   cfg.MongoDB.Username,
		Password:   cfg.MongoDB.Password,
		URL:        cfg.MongoDB.URL,
		Database:   cfg.MongoDB.Database,
		Collection: cfg.UserProfile.CollectionName,
	})
	upService := service.NewUserProfile(fStore, sStore)
	pb.RegisterUserProfileServer(srv, sgrpc.NewService(upService))

	serverNotify := make(chan error)

	go func() {
		logger.Infof("Listing on Port : %s", cfg.UserProfile.Port)
		serverNotify <- srv.Serve(lis)
		close(serverNotify)
	}()

	sysInt := make(chan os.Signal, 1)
	signal.Notify(sysInt, os.Interrupt, syscall.SIGTERM)

	// channel listener for error
	select {
	case sys := <-sysInt:
		logger.Infof("Run sys signal : %s", sys.String())
	case serverErr := <-serverNotify:
		logger.Infof("Server signal : %s", serverErr)
	}

	// close application gracefully
}
