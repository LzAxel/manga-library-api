package errors

import "errors"

var (
	ErrUsernameExists = errors.New("username already exists")
)
