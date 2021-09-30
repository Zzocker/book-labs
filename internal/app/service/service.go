package service

import (
	"context"

	"github.com/Zzocker/book-labs/internal/app/domain"
)

type (
	Auth interface {
		NewToken(ctx context.Context, in domain.NewTokenInput) (string, error)
		GetToken(ctx context.Context, tokenID string) (*domain.Token, error)
		DelToken(ctx context.Context, tokenID string) error
		ExtendTTL(ctx context.Context, tokenID string) (*domain.Token, error)
	}
)
