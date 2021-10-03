package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/Zzocker/book-labs/internal/app/delivery/http/v1/model"
	"github.com/Zzocker/book-labs/internal/app/domain"
	"github.com/Zzocker/book-labs/internal/app/service"
	"github.com/Zzocker/book-labs/pkg/errors"
	pb "github.com/Zzocker/book-labs/protos/profile"
)

type authService struct {
	servs             service.Auth
	profileRPCChannel grpc.ClientConnInterface
}

func newAuthRouters(handler *gin.RouterGroup, servs service.Auth, profileRPCChannel grpc.ClientConnInterface) {
	a := authService{servs: servs, profileRPCChannel: profileRPCChannel}

	h := handler.Group("/auth")
	{
		h.POST("", a.newToken)
		h.GET("", a.getToken)
		h.PATCH("", a.extendTokenTTL)
		h.DELETE("", a.deleteToken)
	}
}

func (a *authService) newToken(c *gin.Context) {
	const op = errors.Op("AuthRouter.newToken")

	var input model.NewTokenBody
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err = errors.E(
			op,
			err,
			errors.CodeInvalidInput,
			errors.SeverityDebug,
		)
		errResponse(c, err)

		return
	}

	loggerWithReqID(c).Debugf("sending request to user-profile service")
	_, err = pb.NewUserProfileClient(a.profileRPCChannel).CheckCred(gRPCCtxWithReqID(c), &pb.CheckCredRequest{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		s := status.Convert(err)
		err = errors.E(
			op,
			fmt.Errorf(s.Message()),
			errors.Code(s.Code()),
			errors.SeverityDebug,
		)
		errResponse(c, err)

		return
	}

	tokenID, err := a.servs.NewToken(c, domain.NewTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		errResponse(c, errors.E(op, err))

		return
	}

	c.String(http.StatusCreated, tokenID)
}

func (a *authService) getToken(c *gin.Context) { //nolint:dupl //no they are not
	const op = errors.Op("AuthRouter.getToken")
	tokenID, err := getTokenID(c, op)
	if err != nil {
		errResponse(c, err)

		return
	}

	token, err := a.servs.GetToken(c, tokenID)
	if err != nil {
		errResponse(c, errors.E(op, err))

		return
	}
	c.JSON(http.StatusOK, model.Token{
		ID:        token.ID,
		Expiry:    token.Expiry,
		IssueTime: token.IssueTime,
		Username:  token.Username,
	})
}

func (a *authService) extendTokenTTL(c *gin.Context) { //nolint:dupl //no they are not
	const op = errors.Op("AuthRouter.extendTokenTTL")
	tokenID, err := getTokenID(c, op)
	if err != nil {
		errResponse(c, err)

		return
	}

	token, err := a.servs.ExtendTTL(c, tokenID)
	if err != nil {
		errResponse(c, errors.E(op, err))

		return
	}
	c.JSON(http.StatusOK, model.Token{
		ID:        token.ID,
		Expiry:    token.Expiry,
		IssueTime: token.IssueTime,
		Username:  token.Username,
	})
}

func (a *authService) deleteToken(c *gin.Context) {
	const op = errors.Op("AuthRouter.deleteToken")
	tokenID, err := getTokenID(c, op)
	if err != nil {
		errResponse(c, err)

		return
	}

	err = a.servs.DelToken(c, tokenID)
	if err != nil {
		errResponse(c, errors.E(op, err))

		return
	}
	c.Status(http.StatusNoContent)
}

func getTokenID(c *gin.Context, op errors.Op) (string, error) {
	tokenString := strings.Split(c.GetHeader("Authorization"), " ")
	if len(tokenString) != 2 { //nolint:gomnd //Authorization: Bearer <tokenID>
		return "", errors.E(
			op,
			fmt.Errorf("invalid auth token"),
			errors.CodeInvalidInput,
			errors.SeverityDebug,
		)
	}

	return tokenString[1], nil
}
