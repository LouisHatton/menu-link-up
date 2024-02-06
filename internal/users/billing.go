package users

import (
	"time"

	"github.com/stripe/stripe-go/v76"
)

type Billing struct {
	PlanName          string                        `json:"planName"`
	BillingInterval   stripe.PriceRecurringInterval `json:"billingInterval"`
	Price             int64                         `json:"price"`
	CurrentPeriodEnd  *time.Time                    `json:"currentPeriodEnd,"`
	CancelAtPeriodEnd bool                          `json:"cancelAtPeriodEnd"`
	DefaultPayment    *BillingDefaultPayment        `json:"defaultPayment,omitempty"`
}

type BillingDefaultPayment struct {
	Brand        string `json:"brand,omitempty"`
	ExpiresMonth int    `json:"expiresMonth,omitempty"`
	ExpiresYear  int    `json:"expiresYear,omitempty"`
	LastFour     string `json:"lastFour,omitempty"`
}

type CustomerPortalLink struct {
	Url string `json:"url"`
}
