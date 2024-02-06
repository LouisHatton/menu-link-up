package users

import "context"

type Service interface {
	GetById(ctx context.Context, id string) (*User, error)
	DeleteById(ctx context.Context, id string) error
	GetBandwidthLimits(ctx context.Context, id string) (*BandwidthLimits, error)
	GetBilling(ctx context.Context, id string) (*Billing, error)
	UpdateBillingLink(ctx context.Context, id string) (*CustomerPortalLink, error)
}
