package clients

import (
	"context"
	"fmt"

	"github.com/Zzocker/book-labs/pkg/errors"
)

func NewMockUserProfileClient() UserProfileClient {
	return &upMockClient{backend: map[string]string{
		"user1": "user1pw",
		"user2": "user2pw",
		"user3": "user3pw",
	}}
}

type upMockClient struct {
	backend map[string]string
}

func (u *upMockClient) CheckCredentails(ctx context.Context, userID, password string) error {
	const op = errors.Op("MockUserProfileClient.CheckCredentails")
	pw, ok := u.backend[userID]
	if !ok {
		return errors.E(op, fmt.Errorf("user-notfound"), errors.CodeNotFound)
	}

	if password != pw {
		return errors.E(op, fmt.Errorf("invalid password"), errors.CodeUnauthenticated)
	}

	return nil
}
