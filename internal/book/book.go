package book

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zzocker/book-labs/config"
	"github.com/Zzocker/book-labs/internal/book/delivery"
	"github.com/Zzocker/book-labs/internal/common"
	"github.com/Zzocker/book-labs/pkg/logger"
	pb "github.com/Zzocker/book-labs/protos/book"
	"google.golang.org/grpc"
)

func Run(cfg *config.BookSharing) {
	lis, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", cfg.Book.Port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(common.LoggerServerInteroceptor),
	)

	pb.RegisterBookServiceServer(srv, delivery.NewBookService())

	serverNotify := make(chan error)

	go func() {
		logger.Infof("Listing on Port : %s", cfg.Book.Port)
		serverNotify <- srv.Serve(lis)
		close(serverNotify)
	}()

	sysInt := make(chan os.Signal, 1)

	signal.Notify(sysInt, os.Interrupt, syscall.SIGTERM)

	select {
	case sys := <-sysInt:
		logger.Infof("Run sys signal : %s", sys.String())
	case serverErr := <-serverNotify:
		logger.Infof("Server signal : %s", serverErr)
	}
}
