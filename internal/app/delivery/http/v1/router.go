package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Zzocker/book-labs/internal/app/service"
	"github.com/Zzocker/book-labs/pkg/logger"
)

const reqIDKey = "REQUEST_ID_KEY"

// services...
func NewRouter(engine *gin.Engine, auth service.Auth) {
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
		newAuthRouters(h, auth)
	}
}
