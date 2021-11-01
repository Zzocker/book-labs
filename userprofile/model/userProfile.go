package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/Zzocker/book-labs/pkg/errors"
	"github.com/Zzocker/book-labs/pkg/utils"
)

type Userprofile struct {
	ID           string
	Name         string
	Email        string
	ProfilePicID string
	EthPublicKey string
	HashedSecret string
	Rating       utils.NetRating
	RegisteredAt int64
}

func (u *Userprofile) SetPassword(op errors.Op, pass string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return errors.E(op, fmt.Errorf("error hasing password: %w", err), errors.CodeInternal)
	}
	u.HashedSecret = string(h)

	return nil
}

func (u *Userprofile) CheckPassword(op errors.Op, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedSecret), []byte(pass))
	if err != nil {
		return errors.E(op, fmt.Errorf("wrong password: %w", err), errors.CodeUnauthenticated)
	}

	return nil
}
