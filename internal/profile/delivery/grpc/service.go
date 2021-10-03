package grpc

import (
	"context"

	"github.com/Zzocker/book-labs/internal/profile/domain"
	"github.com/Zzocker/book-labs/internal/profile/service"
	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/protos/common"
	pb "github.com/Zzocker/book-labs/protos/profile"
)

func NewService(upService service.UserProfile) pb.UserProfileServer {
	return &server{upService: upService}
}

type server struct {
	upService service.UserProfile
	pb.UnimplementedUserProfileServer
}

func (s *server) CreateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*common.EmptyRequest, error) {
	const op = errors.Op("UserProfileGrpcService.CreateProfile")

	err := s.upService.CreateProfile(ctx, &domain.UserProfile{
		ID:       in.GetUsername(),
		Email:    in.GetEmail(),
		Name:     in.GetName(),
		Password: in.GetPassword(),
	}, in.GetProfilePic())
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &common.EmptyRequest{}, nil
}

func (s *server) QueryProfile(ctx context.Context, in *pb.QueryProfileRequest) (*pb.Profiles, error) {
	return nil, nil
}

func (s *server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.Profile, error) {
	const op = errors.Op("UserProfileGrpcService.UpdateProfile")

	profile, err := s.upService.UpdateProfile(ctx, in)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return profile, nil
}

func (s *server) DeleteProfile(ctx context.Context, in *pb.UsernameRequest) (*common.EmptyRequest, error) {
	return nil, nil
}

func (s *server) GetProfile(ctx context.Context, in *pb.UsernameRequest) (*pb.Profile, error) {
	const op = errors.Op("UserProfileGrpcService.GetProfile")

	profile, err := s.upService.GetProfile(ctx, in.Username)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return profile, nil
}

func (s *server) GetProfilePic(ctx context.Context, in *pb.UsernameRequest) (*common.Image, error) {
	const op = errors.Op("UserProfileGrpcService.GetProfile")

	pic, err := s.upService.GetProfilePic(ctx, in.Username)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return pic, nil
}

func (s *server) CheckCred(ctx context.Context, in *pb.CheckCredRequest) (*common.EmptyRequest, error) {
	const op = errors.Op("UserProfileGrpcService.CheckCred")

	err := s.upService.CheckCred(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &common.EmptyRequest{}, nil
}
