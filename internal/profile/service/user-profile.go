package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Zzocker/book-labs/internal/common"
	"github.com/Zzocker/book-labs/internal/profile/domain"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
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
		lg.Debugf("storing user's profile picture")
		id := uuid.New().String()
		err = u.fileStore.Put(ctx, &datastore.File{
			ID:   id,
			Data: pppic,
		})

		if err != nil {
			return errors.E(op, err)
		}
		profile.ProfilePicID = id
	}
	lg.Debugf("storing user's profile")
	err = u.smartStore.Store(ctx, profile)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
