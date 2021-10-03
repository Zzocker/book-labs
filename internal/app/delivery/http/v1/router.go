package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/Zzocker/book-labs/internal/app/service"
	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
)

const (
	reqIDKey    = "REQUEST_ID_KEY"
	usernameKey = "USERNAME_KEY"
)

type RPCEndpoints struct {
	UserProfile string
}

type GRPCChannel struct {
	Userprofile grpc.ClientConnInterface
}

// services...
func NewRouter(engine *gin.Engine, auth service.Auth, channels GRPCChannel) {
	engine.Use(gin.Recovery())
	engine.Use(func(c *gin.Context) {
		reqID := uuid.New().String()
		logger.WithFields(map[string]interface{}{
			"request_id": reqID,
		}).Infof("[%s] %s", c.Request.Method, c.Request.URL.Path)
		c.Set(reqIDKey, reqID)
		c.Next()
	})

	h := engine.Group("/api/v1")
	{
		newAuthRouters(h, auth, channels.Userprofile)
	}

	authorized := engine.Group("/api/v1")
	authorized.Use(func(c *gin.Context) {
		const op = errors.Op("Middleware.Auth")
		if c.Request.Method == http.MethodPost && c.Request.URL.String() == "/api/v1/user_profile" {
			c.Next()

			return
		}
		tokenID, err := getTokenID(c, op)
		if err != nil {
			errResponse(c, err)

			return
		}
		token, err := auth.GetToken(c, tokenID)
		if err != nil {
			errResponse(c, errors.E(op, err))

			return
		}
		c.Set(usernameKey, token.Username)
		c.Next()
	})
	{
		newUserProfileRouter(authorized, channels.Userprofile)
	}
}
