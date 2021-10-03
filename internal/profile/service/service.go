package service

import (
	"context"

	"github.com/Zzocker/book-labs/internal/profile/domain"
	"github.com/Zzocker/book-labs/protos/common"
	pb "github.com/Zzocker/book-labs/protos/profile"
)

type (
	UserProfile interface {
		CreateProfile(ctx context.Context, profile *domain.UserProfile, ppic []byte) error
		CheckCred(ctx context.Context, username, password string) error
		GetProfile(ctx context.Context, username string) (*pb.Profile, error)
		GetProfilePic(ctx context.Context, username string) (*common.Image, error)
		UpdateProfile(ctx context.Context, profile *pb.UpdateProfileRequest) (*pb.Profile, error)
	}
)
