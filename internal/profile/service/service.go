package service

import (
	"context"

	"github.com/Zzocker/book-labs/internal/profile/domain"
)

type (
	UserProfile interface {
		CreateProfile(ctx context.Context, profile *domain.UserProfile, ppic []byte) error
	}
)
