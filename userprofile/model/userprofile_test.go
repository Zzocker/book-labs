package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Zzocker/book-labs/pkg/errors"
)

func TestPassword(t *testing.T) {
	is := assert.New(t)

	var user Userprofile
	const pass string = "testPassword"
	const op = errors.Op("Helper")
	err := user.SetPassword(op, pass)
	is.NoError(err)

	err = user.CheckPassword(op, pass)
	is.NoError(err)

	err = user.CheckPassword(op, "wrong")
	is.Error(err)
	is.Equal(errors.CodeUnauthenticated, errors.ErrCode(err))
}
