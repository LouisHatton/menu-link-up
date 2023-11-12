package files

import (
	"context"
)

type Service interface {
	GetById(ctx context.Context, id string) (*File, error)
	GetByUserId(ctx context.Context, userId string) (*[]File, error)
	Create(ctx context.Context, userId string, newFile NewFile) (*FileUpload, error)
	Edit(ctx context.Context, id string, newFile NewFile) error
	Delete(ctx context.Context, id string) error
	GetS3LinkFromSlug(ctx context.Context, slug string) (string, error)
}
