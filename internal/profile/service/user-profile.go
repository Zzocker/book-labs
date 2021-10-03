package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/Zzocker/book-labs/internal/common"
	"github.com/Zzocker/book-labs/internal/profile/domain"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
	cpb "github.com/Zzocker/book-labs/protos/common"
	pb "github.com/Zzocker/book-labs/protos/profile"
)

type userProfile struct {
	fileStore  datastore.FileDataStore
	smartStore datastore.SmartDataStore
}

func NewUserProfile(fStore datastore.FileDataStore, sStore datastore.SmartDataStore) UserProfile {
	return &userProfile{fileStore: fStore, smartStore: sStore}
}

func (u *userProfile) CreateProfile(ctx context.Context, profile *domain.UserProfile, pppic []byte) error {
	const op = errors.Op("UserProfileService.CreateProfile")
	lg := common.GetLoggerWithReqID(ctx)

	hPass, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed to hash password : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}
	profile.Password = string(hPass)

	if len(pppic) != 0 {
		profile.ProfilePicID = uuid.New().String()
	}
	lg.Debugf("storing user's profile")
	err = u.smartStore.Store(ctx, profile)
	if err != nil {
		return errors.E(op, err)
	}

	if len(pppic) != 0 {
		lg.Debugf("storing user's profile picture")
		err = u.fileStore.Put(ctx, &datastore.File{
			ID:   profile.ProfilePicID,
			Data: pppic,
		})

		if err != nil {
			return errors.E(op, err)
		}
	}

	return nil
}

func (u *userProfile) CheckCred(ctx context.Context, username, password string) error {
	const op = errors.Op("UserProfileService.CheckCred")
	lg := common.GetLoggerWithReqID(ctx)

	lg.Debugf("fetching user-profile for %s", username)
	raw, err := u.smartStore.Get(ctx, map[string]interface{}{
		"_id": username,
	})
	if err != nil {
		return errors.E(op, err)
	}

	var profile domain.UserProfile
	err = bson.Unmarshal(raw, &profile)
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed to parse user-profile : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	lg.Debugf("comparing password")
	err = bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(password))
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("wrong credentials"),
			errors.CodeUnauthorized,
			errors.SeverityDebug,
		)
	}

	return nil
}

func (u *userProfile) GetProfile(ctx context.Context, username string) (*pb.Profile, error) {
	const op = errors.Op("UserProfileService.GetProfile")
	lg := common.GetLoggerWithReqID(ctx)

	lg.Debugf("fetching user-profile of %s", username)
	raw, err := u.smartStore.Get(ctx, map[string]interface{}{
		"_id": username,
	})
	if err != nil {
		return nil, errors.E(op, err)
	}
	var profile domain.UserProfile
	err = bson.Unmarshal(raw, &profile)
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to parse user-profile : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return &pb.Profile{
		Username: profile.ID,
		Name:     profile.Name,
		Email:    profile.Email,
	}, nil
}

func (u *userProfile) GetProfilePic(ctx context.Context, username string) (*cpb.Image, error) {
	const op = errors.Op("UserProfileService.GetProfilePic")
	lg := common.GetLoggerWithReqID(ctx)

	lg.Debugf("fetching user-profile of %s", username)
	raw, err := u.smartStore.Get(ctx, map[string]interface{}{
		"_id": username,
	})
	if err != nil {
		return nil, errors.E(op, err)
	}
	var profile domain.UserProfile
	err = bson.Unmarshal(raw, &profile)
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to parse user-profile : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	lg.Debugf("fetching user's profile pic")
	file, err := u.fileStore.Get(ctx, profile.ProfilePicID)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &cpb.Image{
		Data: file.Data,
	}, nil
}

func (u *userProfile) UpdateProfile(ctx context.Context, profile *pb.UpdateProfileRequest) (*pb.Profile, error) {
	const op = errors.Op("UserProfileService.CreateProfile")
	lg := common.GetLoggerWithReqID(ctx)

	lg.Debugf("fetching %s's profile", profile.Username)
	raw, err := u.smartStore.Get(ctx, map[string]interface{}{
		"_id": profile.Username,
	})
	if err != nil {
		return nil, errors.E(op, err)
	}
	var old domain.UserProfile
	err = bson.Unmarshal(raw, &old)
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to parse user-profile : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	if len(profile.GetProfilePic()) != 0 {
		lg.Debugf("updating profile pic")
		err = u.fileStore.Put(ctx, &datastore.File{
			ID:   old.ProfilePicID,
			Data: profile.ProfilePic,
		})
		if err != nil {
			return nil, errors.E(op, err)
		}
	}

	if profile.Email != "" {
		old.Email = profile.Email
	}
	if profile.Name != "" {
		old.Name = profile.Name
	}
	if profile.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost) //nolint:errcheck //no required
		old.Password = string(hashed)
	}
	lg.Debugf("updating user profile")
	err = u.smartStore.Update(ctx, map[string]interface{}{
		"_id": profile.Username,
	}, old)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &pb.Profile{
		Username: old.ID,
		Name:     old.Name,
		Email:    old.Email,
	}, nil
}
