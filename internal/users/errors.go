package users

import "errors"

var (
	ErrUserNotFound     = errors.New("the requested user was not found")
	ErrEmailNotVerified = errors.New("the users email has not been verified")
)
