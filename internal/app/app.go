package app

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Zzocker/book-labs/config"
	v1 "github.com/Zzocker/book-labs/internal/app/delivery/http/v1"
	"github.com/Zzocker/book-labs/internal/app/service/auth"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/logger"
)

func Run(cfg *config.BookSharing) {
	// setup
	redisStore := datastore.NewRedisDumbDatastore(&datastore.RedisDumbDatastoreConfig{
		ExpiryS:  cfg.App.OAuthRedis.ExpiryS,
		URL:      cfg.App.OAuthRedis.URL,
		Username: cfg.App.OAuthRedis.Username,
		Password: cfg.App.OAuthRedis.Password,
		Database: cfg.App.OAuthRedis.Database,
	})
	// services
	authService := auth.NewAuthService(redisStore, cfg.App.OAuthRedis.ExpiryS)
	//
	engine := gin.New()
	v1.NewRouter(engine, authService)

	server := http.Server{
		Handler:      engine,
		Addr:         net.JoinHostPort("127.0.0.1", cfg.App.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// error chan for server error
	serverNotify := make(chan error)

	go func() {
		logger.Infof("Listing on Port : %s", cfg.App.Port)
		serverNotify <- server.ListenAndServe()
		close(serverNotify)
	}()

	// signal channel for other sys interrupts
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
