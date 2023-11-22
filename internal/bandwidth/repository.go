package bandwidth

import (
	"context"

	"github.com/LouisHatton/menu-link-up/internal/common/repository"
)

type Repository interface {
	repository.CrudRepository[MonthlyBandwidth]

	IncreaseBytesTransferred(ctx context.Context, id string, bytesTransferred int64) error
	IncreaseBytesUploaded(ctx context.Context, id string, bytesTransferred int64) error
	GetByUserIdMonthYear(ctx context.Context, userId string, month int, year int) (*MonthlyBandwidth, error)
	DeleteByUserId(ctx context.Context, userId string) error
}
