package core

import (
	"context"
	"fmt"

	"github.com/Zzocker/book-labs/auth/clients"
	"github.com/Zzocker/book-labs/auth/model"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
)

type AuthCore struct {
	store   datastore.KVStore
	up      clients.UserProfileClient
	expiryS int
}

const (
	tokenPrefix string = "AUTH_TOKEN:"
)

func (a *AuthCore) Create(ctx context.Context, userID, password string) (*model.AuthToken, error) {
	const op = errors.Op("AuthCore.Create")
	err := a.up.CheckCredentails(ctx, userID, password)
	if err != nil {
		return nil, errors.E(op, err)
	}

	var token model.AuthToken
	token.New(userID, a.expiryS)
	err = a.store.Set(ctx, getKey(token.ID), token.ToBytes())
	if err != nil {
		return nil, errors.E(op, err)
	}

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

	token.Refresh(a.expiryS)
	err = a.store.Set(ctx, getKey(id), token.ToBytes())
	if err != nil {
		return nil, errors.E(op, err)
	}

	return token, nil
}

func (a *AuthCore) Delete(ctx context.Context, id string) error {
	const op = errors.Op("AuthCode.Delete")
	err := a.store.Del(ctx, getKey(id))
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func getKey(id string) string {
	return tokenPrefix + id
}
