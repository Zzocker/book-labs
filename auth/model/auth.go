package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Zzocker/book-labs/pkg/errors"
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

func (a *AuthToken) ToBytes() []byte {
	raw, _ := json.Marshal(a) //nolint:errcheck //not required

	return raw
}

func (a *AuthToken) FromBytes(op errors.Op, raw []byte) error {
	err := json.Unmarshal(raw, a)
	if err != nil {
		return errors.E(op, fmt.Errorf("failed to unmarshal token: %w", err), errors.CodeInternal)
	}

	return nil
}
