package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Zzocker/book-labs/internal/app/domain"
	"github.com/Zzocker/book-labs/internal/app/service"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
)

type authService struct {
	store   datastore.DumbDataStore
	expiryS int64
}

func NewAuthService(store datastore.DumbDataStore, expiryS int64) service.Auth {
	return &authService{store: store, expiryS: expiryS}
}

func (a *authService) NewToken(ctx context.Context, in domain.NewTokenInput) (string, error) {
	const op = errors.Op("AuthService.NewToken")
	// authenticate
	issueTime := time.Now()
	token := domain.Token{
		ID:        uuid.New().String(),
		Expiry:    issueTime.Add(time.Second * time.Duration(a.expiryS)).Unix(),
		IssueTime: issueTime.Unix(),
		Username:  in.Username,
	}
	raw, _ := json.Marshal(token) // nolint:errcheck //if thow error means panic
	err := a.store.Put(ctx, token.ID, raw)
	if err != nil {
		return "", errors.E(op, err)
	}

	return token.ID, nil
}

func (a *authService) GetToken(ctx context.Context, tokenID string) (*domain.Token, error) {
	const op = errors.Op("AuthService.GetToken")

	raw, err := a.store.Get(ctx, tokenID)
	if err != nil && errors.ErrCode(err) == errors.CodeNotFound {
		return nil, errors.E(
			op,
			fmt.Errorf("token expired"),
			errors.CodeUnauthorized,
			errors.SeverityDebug,
		)
	} else if err != nil {
		return nil, errors.E(op, err)
	}
	var token domain.Token
	err = json.Unmarshal(raw, &token)
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to unmarshal token data : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return &token, nil
}

func (a *authService) DelToken(ctx context.Context, tokenID string) error {
	const op = errors.Op("AuthService.DelToken")
	err := a.store.Delete(ctx, tokenID)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (a *authService) ExtendTTL(ctx context.Context, tokenID string) (*domain.Token, error) {
	const op = errors.Op("AuthService.ExtendTTL")

	token, err := a.GetToken(ctx, tokenID)
	if err != nil {
		return nil, errors.E(op, err)
	}
	token.Expiry = time.Now().Add(time.Second * time.Duration(a.expiryS)).Unix()

	raw, _ := json.Marshal(token) // nolint:errcheck //if thow error means panic
	err = a.store.Put(ctx, token.ID, raw)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return token, nil
}
