package service

import (
	"context"

	internal_context "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/LouisHatton/menu-link-up/internal/subscriptions"
	"github.com/LouisHatton/menu-link-up/internal/users"
	"github.com/stripe/stripe-go/v76"
)

func (svc *UserService) GetBilling(ctx context.Context, id string) (*users.Billing, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(id))

	user, err := svc.GetById(ctx, id)
	if err != nil {
		logger.Error("attempting to get user from service", log.Error(err))
		return nil, err
	}
	logger = svc.logger.With(log.SubscriptionId(user.StripeSubscriptionId))

	subscription, err := svc.subscriptionSvc.GetSubscription(ctx, user.StripeSubscriptionId)
	if err != nil {
		logger.Error("attempting to get users subscription", log.Error(err))
		return nil, err
	}

	price := subscription.Items.Data[0].Price

	product, err := svc.subscriptionSvc.GetProduct(ctx, price.Product.ID)
	if err != nil {
		logger.Error("attempting to get subscription product", log.Error(err))
		return nil, err
	}

	var card *stripe.PaymentMethodCard

	// default_payment_method
	if subscription.DefaultPaymentMethod != nil {
		if subscription.DefaultPaymentMethod.Type == stripe.PaymentMethodTypeCard {
			card = subscription.DefaultPaymentMethod.Card
		}

	} else {
		customer, err := svc.subscriptionSvc.GetCustomer(ctx, user.StripeCustomerId)
		if err != nil {
			logger.Error("attempting to get customer", log.Error(err))
			return nil, err
		}
		if customer.InvoiceSettings.DefaultPaymentMethod != nil {
			if customer.InvoiceSettings.DefaultPaymentMethod.Type == stripe.PaymentMethodTypeCard {
				card = customer.InvoiceSettings.DefaultPaymentMethod.Card
			}
		}
	}

	var defaultPayment users.BillingDefaultPayment
	if card != nil {
		defaultPayment.Brand = string(card.Brand)
		defaultPayment.ExpiresMonth = int(card.ExpMonth)
		defaultPayment.ExpiresYear = int(card.ExpYear)
		defaultPayment.LastFour = card.Last4
	}

	return &users.Billing{
		PlanName:          product.Name,
		BillingInterval:   price.Recurring.Interval,
		Price:             price.UnitAmount,
		CurrentPeriodEnd:  subscriptions.StripeTime(subscription.CurrentPeriodEnd),
		CancelAtPeriodEnd: subscription.CancelAtPeriodEnd,
		DefaultPayment:    &defaultPayment,
	}, nil
}

func (svc *UserService) UpdateBillingLink(ctx context.Context, id string) (*users.CustomerPortalLink, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(id))

	user, err := svc.GetById(ctx, id)
	if err != nil {
		logger.Error("attempting to get user from service", log.Error(err))
		return nil, err
	}

	url, err := svc.subscriptionSvc.PortalUpdateBilling(ctx, user.StripeCustomerId)
	if err != nil {
		logger.Error("attempting to create customer billing portal link", log.Error(err))
		return nil, err
	}

	return &users.CustomerPortalLink{
		Url: url,
	}, nil
}
