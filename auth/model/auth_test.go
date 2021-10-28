package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/pkg/errors"
)

func TestNew(t *testing.T) {
	is := assert.New(t)

	const expiryS int = 50
	const userID string = "testUserID"
	var token AuthToken
	token.New(userID, expiryS)
	is.Equal(userID, token.UserID)
	is.Equal(int64(expiryS), token.ExpiryTime-token.CreationTime)
}

func TestRefresh(t *testing.T) {
	is := assert.New(t)

	var token AuthToken
	token.Refresh(1)
	is.NotEqual(int64(0), token.ExpiryTime)
}

func TestIsExpired(t *testing.T) {
	is := assert.New(t)

	var token AuthToken
	token.New("", 1)
	is.False(token.IsExpired())

	token.ExpiryTime = token.CreationTime - 1
	is.True(token.IsExpired())
}

func TestToBytes(t *testing.T) {
	is := assert.New(t)

	token := AuthToken{
		ID:           "fakeID",
		UserID:       "fakeUserID",
		CreationTime: 0,
		ExpiryTime:   1,
	}
	raw := token.ToBytes()
	is.NotNil(raw)
}

func TestFromBytes(t *testing.T) {
	is := assert.New(t)

	t.Run("Ok", func(t *testing.T) {
		token := AuthToken{
			ID:           "fakeID",
			UserID:       "fakeUserID",
			CreationTime: 0,
			ExpiryTime:   1,
		}
		raw, _ := json.Marshal(token) //nolint:errcheck //not required in test
		var tk AuthToken
		err := tk.FromBytes(errors.Op("Test"), raw)
		is.NoError(err)
		is.Equal(token, tk)
	})

	t.Run("fail", func(t *testing.T) {
		var token AuthToken
		err := token.FromBytes("test", []byte{})
		is.Error(err)
	})
}
