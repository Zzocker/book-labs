package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/Zzocker/book-labs/internal/app/delivery/http/v1/model"
	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
)

func errResponse(c *gin.Context, err error) {
	reqID := c.GetString(reqIDKey)
	code := errors.ErrCodeToHTTP(errors.ErrCode(err))
	errors.Set(err, errors.ReqID(reqID))
	logger.SystemErr(err)
	c.JSON(code, model.Error{Code: code, Message: err.Error()})
}
