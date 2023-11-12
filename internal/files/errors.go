package files

import "errors"

var (
	ErrFileNotFound = errors.New("the file requested was not found")
	ErrNotUsersFile = errors.New("the requesting user does not have permission to edit this file")
)
