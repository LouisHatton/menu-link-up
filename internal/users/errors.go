package users

import "errors"

var (
	ErrUserNotFound = errors.New("the requested user was not found")
)
