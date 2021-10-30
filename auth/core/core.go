package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Zzocker/book-labs/auth/clients"
	"github.com/Zzocker/book-labs/auth/model"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
)

type AuthCore struct {
	store   datastore.KVStore
	up      clients.UserProfileClient
	expiryS int
}

const (
	tokenPrefix  string = "AUTH_TOKEN:"
	logoutPrefix string = "LOGOUT:"
)

func (a *AuthCore) Create(ctx context.Context, userID, password string) (*model.AuthToken, error) {
	const op = errors.Op("AuthCore.Create")
	logger.Debugf("checking credential for userID = %s", userID)
	err := a.up.CheckCredentails(ctx, userID, password)
	if err != nil {
		return nil, errors.E(op, err)
	}

	logger.Debug("deleting existing logout blocks")
	err = a.store.Del(ctx, getLogoutKey(userID))
	if err != nil && errors.ErrCode(err) != errors.CodeNotFound {
		return nil, errors.E(op, err)
	}

	var token model.AuthToken
	logger.Debug("storing auth token")
	token.New(userID, a.expiryS)
	err = a.store.Set(ctx, getKey(token.ID), token.ToBytes())
	if err != nil {
		return nil, errors.E(op, err)
	}
	logger.Debugf("new auth token created for %s", userID)

	return &token, nil
}

func (a *AuthCore) Get(ctx context.Context, id string) (*model.AuthToken, error) {
	const op = errors.Op("AuthCore.Get")

	raw, err := a.store.Get(ctx, getKey(id))
	if err != nil {
		return nil, errors.E(op, err)
	}

	var token model.AuthToken
	err = token.FromBytes(op, raw)
	if err != nil {
		return nil, err
	}

	logger.Debug("get logout block for userID = %s", token.UserID)
	raw, err = a.store.Get(ctx, getLogoutKey(token.UserID))
	if err != nil && errors.ErrCode(err) != errors.CodeNotFound {
		return nil, errors.E(op, err)
	}
	if len(raw) != 0 {
		logger.Debug("checking logout block on user")
		blockTime, err := strconv.Atoi(string(raw))
		if err != nil {
			return nil, errors.E(op, err, errors.CodeInternal)
		}
		if token.CreationTime <= int64(blockTime) {
			return nil, errors.E(op, fmt.Errorf("token expired"), errors.CodeUnauthenticated)
		}
	}

	logger.Debug("checking is token expired")
	if token.IsExpired() {
		return nil, errors.E(op, fmt.Errorf("token expired"), errors.CodeUnauthenticated)
	}

	return &token, nil
}

func (a *AuthCore) Refresh(ctx context.Context, id string) (*model.AuthToken, error) {
	const op = errors.Op("AuthCore.Create")

	token, err := a.Get(ctx, id)
	if err != nil {
		return nil, errors.E(op, err)
	}

	logger.Debugf("refresh token for userID = %s", token.UserID)
	token.Refresh(a.expiryS)
	err = a.store.Set(ctx, getKey(id), token.ToBytes())
	if err != nil {
		return nil, errors.E(op, err)
	}

	logger.Debug("token refreshed")

	return token, nil
}

func (a *AuthCore) Delete(ctx context.Context, id string) error {
	const op = errors.Op("AuthCore.Delete")
	err := a.store.Del(ctx, getKey(id))
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (a *AuthCore) Logout(ctx context.Context, userID string) error {
	const op = errors.Op("AuthCore.Logout")
	err := a.store.Set(ctx, getLogoutKey(userID), []byte(fmt.Sprintf("%d", time.Now().Unix())))
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func getKey(id string) string {
	return tokenPrefix + id
}

func getLogoutKey(userID string) string {
	return logoutPrefix + userID
}
