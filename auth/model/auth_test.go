package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
