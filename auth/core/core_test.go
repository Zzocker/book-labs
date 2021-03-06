package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/auth/clients"
	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/logger"
)

func TestCreate(t *testing.T) {
	logger.Setup("debug", "AuthCode_Test", "1.0.0")
	is := assert.New(t)

	ctx := context.Background()
	aCore := AuthCore{
		store:   datastore.NewMockKVStore(),
		up:      clients.NewMockUserProfileClient(),
		expiryS: 1,
	}

	t.Run("Happy Flow", func(t *testing.T) {
		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)
		is.Equal("user1", token.UserID)
		_, err = aCore.store.Get(ctx, getKey(token.ID))
		is.NoError(err)
	})

	t.Run("wrong cred", func(t *testing.T) {
		_, err := aCore.Create(ctx, "user1", "wrongPW")
		is.Error(err)
	})

	t.Run("logged-out", func(t *testing.T) {
		err := aCore.Logout(ctx, "user1")
		is.NoError(err)

		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)
		is.NotNil(token)

		_, err = aCore.store.Get(ctx, getLogoutKey("user1"))
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})
}

func TestGet(t *testing.T) {
	logger.Setup("debug", "AuthCode_Test", "1.0.0")
	is := assert.New(t)

	ctx := context.Background()
	aCore := AuthCore{
		store:   datastore.NewMockKVStore(),
		up:      clients.NewMockUserProfileClient(),
		expiryS: 1,
	}
	t.Run("Happy Flow", func(t *testing.T) {
		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)

		tk, err := aCore.Get(ctx, token.ID)
		is.NoError(err)
		is.Equal(*tk, *token)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := aCore.Get(ctx, "not-found")
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Expired", func(t *testing.T) {
		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)

		time.Sleep(2 * time.Second)
		_, err = aCore.Get(ctx, token.ID)
		is.Error(err)
		is.Equal(errors.CodeUnauthenticated, errors.ErrCode(err))
	})

	t.Run("logged-out", func(t *testing.T) {
		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)

		err = aCore.Logout(ctx, "user1")
		is.NoError(err)

		_, err = aCore.Get(ctx, token.ID)
		is.Error(err)
		is.Equal(errors.CodeUnauthenticated, errors.ErrCode(err))
	})
}

func TestRefresh(t *testing.T) {
	logger.Setup("debug", "AuthCode_Test", "1.0.0")
	is := assert.New(t)

	ctx := context.Background()
	aCore := AuthCore{
		store:   datastore.NewMockKVStore(),
		up:      clients.NewMockUserProfileClient(),
		expiryS: 1,
	}

	t.Run("Happy Flow", func(t *testing.T) {
		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)

		time.Sleep(time.Second)
		tk, err := aCore.Refresh(ctx, token.ID)
		is.NoError(err)
		is.NotEqual(token.ExpiryTime, tk.ExpiryTime)
		is.Equal(token.UserID, tk.UserID)
		is.Equal(token.CreationTime, tk.CreationTime)
		_ = tk
	})

	t.Run("not found", func(t *testing.T) {
		_, err := aCore.Refresh(ctx, "not-found")
		is.Error(err)
	})
}

func TestDelete(t *testing.T) {
	logger.Setup("debug", "AuthCode_Test", "1.0.0")
	is := assert.New(t)

	ctx := context.Background()
	aCore := AuthCore{
		store:   datastore.NewMockKVStore(),
		up:      clients.NewMockUserProfileClient(),
		expiryS: 1,
	}

	t.Run("Happy Flow", func(t *testing.T) {
		token, err := aCore.Create(ctx, "user1", "user1pw")
		is.NoError(err)

		err = aCore.Delete(ctx, token.ID)
		is.NoError(err)

		_, err = aCore.store.Get(ctx, getKey(token.ID))
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("not found", func(t *testing.T) {
		err := aCore.Delete(ctx, "not-found")
		is.Error(err)
	})
}

func TestLogout(t *testing.T) {
	logger.Setup("debug", "AuthCode_Test", "1.0.0")
	is := assert.New(t)

	ctx := context.Background()
	aCore := AuthCore{
		store:   datastore.NewMockKVStore(),
		up:      clients.NewMockUserProfileClient(),
		expiryS: 1,
	}

	err := aCore.Logout(ctx, "user1")
	is.NoError(err)

	raw, err := aCore.store.Get(ctx, getLogoutKey("user1"))
	is.NoError(err)
	is.NotNil(raw)
}
