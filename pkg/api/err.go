package api

import (
	dbsql "database/sql"
	"errors"
	"github.com/acoshift/arpc/v2"
)

var (
	ErrRecordNotFound = arpc.NewError("record not found")

	ErrUsernameExist = arpc.NewError("username already exists")
	ErrLoginFailed   = arpc.NewError("credentials are invalid")
)

func WrapError(err error) error {
	switch {
	case errors.Is(err, dbsql.ErrNoRows):
		return ErrRecordNotFound
	default:
		return err
	}
}
