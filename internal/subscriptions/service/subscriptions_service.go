package service

import (
	"context"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	internal_context "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/LouisHatton/menu-link-up/internal/subscriptions"
	"github.com/LouisHatton/menu-link-up/internal/users"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
)

type SubscriptionSvc struct {
	logger       *log.Logger
	stripeConfig appconfig.Stripe
}

func New(logger *log.Logger, stripeConfig appconfig.Stripe) (*SubscriptionSvc, error) {
	stripe.Key = stripeConfig.Key
	return &SubscriptionSvc{
		logger:       logger,
		stripeConfig: stripeConfig,
	}, nil
}

// CreateCustomer implements subscriptions.Service.
func (svc *SubscriptionSvc) CreateCustomerWithTrial(ctx context.Context, newCustomer subscriptions.NewCustomer) (*subscriptions.Customer, *stripe.Subscription, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(newCustomer.UserId))

	logger.Info("creating new stripe customer with default trial")

	newCustomerParams := stripe.CustomerParams{
		Name:  &newCustomer.Name,
		Email: &newCustomer.Email,
		Metadata: map[string]string{
			"userId": newCustomer.UserId,
		},
	}

	customer, err := customer.New(&newCustomerParams)
	if err != nil {
		msg := "attempting to save new customer in stripe"
		logger.Error(msg, log.Error(err))
		return nil, nil, fmt.Errorf(msg+": %w", err)
	}

	newSubscriptionParams := &stripe.SubscriptionParams{
		Customer: &customer.ID,
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: &svc.stripeConfig.DefaultSubscriptionPrice,
			},
		},
		TrialPeriodDays: &svc.stripeConfig.DefaultTrialLength,
		TrialSettings: &stripe.SubscriptionTrialSettingsParams{
			EndBehavior: &stripe.SubscriptionTrialSettingsEndBehaviorParams{
				MissingPaymentMethod: stripe.String("cancel"),
			},
		},
	}

	subscription, err := subscription.New(newSubscriptionParams)
	if err != nil {
		msg := "attempting to create new subscription for new customer in stripe"
		logger.Error(msg, log.Error(err))
		return nil, nil, fmt.Errorf(msg+": %w", err)
	}

	logger.Info("new trial created", log.SubscriptionId(subscription.ID))

	localCustomer := subscriptions.StripeCustomerToLocal(customer)
	return &localCustomer, subscription, nil
}

func (svc *SubscriptionSvc) GetLimitsForSubscription(ctx context.Context, subscriptionId string) (*users.BandwidthLimits, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.SubscriptionId(subscriptionId))

	sub, err := svc.GetSubscription(ctx, subscriptionId)
	if err != nil {
		msg := "attempting to get customers subscription"
		logger.Error(msg, log.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	product, err := svc.GetProduct(ctx, sub.Items.Data[0].Price.Product.ID)
	if err != nil {
		msg := "attempting to get product for subscription"
		logger.Error(msg, log.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	limits, err := subscriptions.ExtractLimitsFromProduct(product)
	if err != nil {
		msg := "attempting extract limits from new product"
		logger.Error(msg, log.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	return limits, nil
}

func (svc *SubscriptionSvc) GetSubscription(ctx context.Context, subscriptionId string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionId, nil)
}

func (svc *SubscriptionSvc) GetProduct(ctx context.Context, productId string) (*stripe.Product, error) {
	return product.Get(productId, nil)
}

func (svc *SubscriptionSvc) GetCustomer(ctx context.Context, customerId string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{}
	params.AddExpand("invoice_settings.default_payment_method")

	return customer.Get(customerId, params)
}

var _ subscriptions.Service = &SubscriptionSvc{}
