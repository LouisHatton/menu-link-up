package files

import (
	"context"

	"github.com/LouisHatton/menu-link-up/internal/common/repository"
)

type Repository interface {
	repository.CrudRepository[File]

	GetBySlug(ctx context.Context, slug string) (*File, error)
	GetByUserId(ctx context.Context, userId string) (*[]File, error)
}
