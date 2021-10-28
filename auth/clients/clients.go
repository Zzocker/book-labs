package clients

import "context"

type UserProfileClient interface {
	CheckCredentails(ctx context.Context, userID, password string) error
}
