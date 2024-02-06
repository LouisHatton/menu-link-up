package subscriptions

import (
	"context"

	"github.com/LouisHatton/menu-link-up/internal/users"
	"github.com/stripe/stripe-go/v76"
)

type Service interface {
	CreateCustomerWithTrial(ctx context.Context, newCustomer NewCustomer) (*Customer, *stripe.Subscription, error)
	GetLimitsForSubscription(ctx context.Context, subscriptionId string) (*users.BandwidthLimits, error)
	GetSubscription(ctx context.Context, subscriptionId string) (*stripe.Subscription, error)
	GetProduct(ctx context.Context, productId string) (*stripe.Product, error)
	GetCustomer(ctx context.Context, customerId string) (*stripe.Customer, error)
	PortalUpdateBilling(ctx context.Context, customerId string) (portalUrl string, err error)
}
