package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/Zzocker/book-labs/internal/app/delivery/http/v1/model"
	"github.com/Zzocker/book-labs/pkg/errors"
	pb "github.com/Zzocker/book-labs/protos/profile"
)

type userProfileRouter struct {
	channel grpc.ClientConnInterface
}

func newUserProfileRouter(handler *gin.RouterGroup, channel grpc.ClientConnInterface) {
	u := userProfileRouter{channel: channel}
	h := handler.Group("/user_profile")
	{
		h.POST("", u.newUserProfile)
		h.GET("", u.queryUserProfile)
		h.PATCH("", u.updateUserProfile)
		h.DELETE("", u.deleteUserProfile)
		h.GET("/:username", u.getUserProfile)
		h.GET("/pic/:username", u.getUserProfilePic)
	}
}

func (u *userProfileRouter) newUserProfile(ctx *gin.Context) {
	const op = errors.Op("UserProfileRouter.newUserProfile")
	var req model.CreateProfileRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		err = errors.E(
			op,
			err,
			errors.CodeInvalidInput,
			errors.SeverityDebug,
		)
		errResponse(ctx, err)

		return
	}
	var profilePic []byte
	fh, err := ctx.FormFile("profile_pic")
	if err == nil && fh != nil {
		f, err := fh.Open() //nolint:govet //ignore
		if err != nil {
			err = errors.E(
				op,
				fmt.Errorf("failed to read profile_pic file : %w", err),
				errors.CodeInvalidInput,
				errors.SeverityDebug,
			)
			errResponse(ctx, err)

			return
		}
		defer f.Close()
		profilePic = make([]byte, fh.Size)
		f.Read(profilePic) //nolint:errcheck //not required
	}

	lg := loggerWithReqID(ctx)
	lg.Debugf("sending request to user-profile service")
	_, err = pb.NewUserProfileClient(u.channel).CreateProfile(gRPCCtxWithReqID(ctx), &pb.UpdateProfileRequest{
		Username:   req.Username,
		Email:      req.Email,
		Name:       req.Name,
		Password:   req.Password,
		ProfilePic: profilePic,
	})

	if err != nil {
		s := status.Convert(err)
		err = errors.E(
			op,
			fmt.Errorf(s.Message()),
			errors.Code(s.Code()),
			errors.SeverityDebug,
		)
		errResponse(ctx, err)

		return
	}
	lg.Debugf("success")
	ctx.Status(http.StatusCreated)
}

func (u *userProfileRouter) queryUserProfile(ctx *gin.Context) {
}

func (u *userProfileRouter) updateUserProfile(ctx *gin.Context) {
	const op = errors.Op("UserProfileRouter.updateUserProfile")
	var req model.CreateProfileRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		err = errors.E(
			op,
			err,
			errors.CodeInvalidInput,
			errors.SeverityDebug,
		)
		errResponse(ctx, err)

		return
	}
	var profilePic []byte
	fh, err := ctx.FormFile("profile_pic")
	if err == nil && fh != nil {
		f, err := fh.Open() //nolint:govet //ignore
		if err != nil {
			err = errors.E(
				op,
				fmt.Errorf("failed to read profile_pic file : %w", err),
				errors.CodeInvalidInput,
				errors.SeverityDebug,
			)
			errResponse(ctx, err)

			return
		}
		defer f.Close()
		profilePic = make([]byte, fh.Size)
		f.Read(profilePic) //nolint:errcheck //not required
	}

	lg := loggerWithReqID(ctx)
	lg.Debugf("sending request to user-profile service")
	profile, err := pb.NewUserProfileClient(u.channel).UpdateProfile(gRPCCtxWithReqID(ctx), &pb.UpdateProfileRequest{
		Username:   ctx.GetString(usernameKey),
		Email:      req.Email,
		Name:       req.Name,
		Password:   req.Password,
		ProfilePic: profilePic,
	})
	if err != nil {
		s := status.Convert(err)
		err = errors.E(
			op,
			fmt.Errorf(s.Message()),
			errors.Code(s.Code()),
			errors.SeverityDebug,
		)
		errResponse(ctx, err)

		return
	}
	ctx.JSON(http.StatusOK, model.UserProfile{
		Username: profile.GetUsername(),
		Email:    profile.GetEmail(),
		Name:     profile.GetName(),
	})
}

func (u *userProfileRouter) deleteUserProfile(ctx *gin.Context) {
}

func (u *userProfileRouter) getUserProfile(ctx *gin.Context) {
	const op = errors.Op("UserProfileRouter.newUserProfile")
	username := ctx.Param("username")

	loggerWithReqID(ctx).Debugf("sending request to user-profile service")
	profile, err := pb.NewUserProfileClient(u.channel).GetProfile(gRPCCtxWithReqID(ctx), &pb.UsernameRequest{
		Username: username,
	})
	if err != nil {
		s := status.Convert(err)
		err = errors.E(
			op,
			fmt.Errorf(s.Message()),
			errors.Code(s.Code()),
			errors.SeverityDebug,
		)
		errResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, model.UserProfile{
		Username: profile.GetUsername(),
		Email:    profile.GetEmail(),
		Name:     profile.GetName(),
	})
}

func (u *userProfileRouter) getUserProfilePic(ctx *gin.Context) {
	const op = errors.Op("UserProfileRouter.getUserProfilePic")
	username := ctx.Param("username")

	loggerWithReqID(ctx).Debugf("sending request to user-profile service")
	pic, err := pb.NewUserProfileClient(u.channel).GetProfilePic(gRPCCtxWithReqID(ctx), &pb.UsernameRequest{
		Username: username,
	})
	if err != nil {
		s := status.Convert(err)
		err = errors.E(
			op,
			fmt.Errorf(s.Message()),
			errors.Code(s.Code()),
			errors.SeverityDebug,
		)
		errResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, model.Image{
		Data: pic.GetData(),
	})
}