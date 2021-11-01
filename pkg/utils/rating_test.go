package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
)

func TestRating(t *testing.T) {
	is := assert.New(t)

	r := NewRatingStore(datastore.NewMockKVStore())

	ctx := context.Background()

	userID := "user1"
	bookID := "book1"

	t.Run("Set", func(t *testing.T) {
		err := r.Set(ctx, userID, bookID, 5)
		is.NoError(err)

		_, err = r.store.Get(ctx, getRatingKey(userID, bookID))
		is.NoError(err)
	})

	t.Run("Get", func(t *testing.T) {
		rating, err := r.Get(ctx, userID, bookID)
		is.NoError(err)

		is.Equal(5, rating.Value)
	})

	t.Run("Delete", func(t *testing.T) {
		err := r.Delete(ctx, userID, bookID)
		is.NoError(err)

		_, err = r.store.Get(ctx, getRatingKey(userID, bookID))
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})
}
