package users

import "context"

type Service interface {
	GetById(ctx context.Context, id string) (*User, error)
	DeleteById(ctx context.Context, id string) error
}
