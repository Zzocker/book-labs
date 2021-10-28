package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	ID           string
	UserID       string
	CreationTime int64
	ExpiryTime   int64
}

func (a *AuthToken) New(userID string, expiryS int) {
	a.UserID = userID
	a.ID = uuid.New().String()
	now := time.Now()
	a.CreationTime = now.Unix()
	a.ExpiryTime = now.Add(time.Second * time.Duration(expiryS)).Unix()
}

func (a *AuthToken) Refresh(expiryS int) {
	a.ExpiryTime = time.Now().Add(time.Second * time.Duration(expiryS)).Unix()
}

func (a *AuthToken) IsExpired() bool {
	return time.Now().Unix() > a.ExpiryTime
}
