package objectstore

import "errors"

var (
	ErrIncorrectRegion = errors.New("the provided region does not match the assigned")
)
