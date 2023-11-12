package repository

import "context"

type CrudRepository[T any] interface {
	Ping() error
	Count(ctx context.Context) (int, error)
	GetById(ctx context.Context, id string) (*T, error)
	DeleteById(ctx context.Context, id string) error
	Update(ctx context.Context, T *T) error
	Create(ctx context.Context, T *T) error
}
