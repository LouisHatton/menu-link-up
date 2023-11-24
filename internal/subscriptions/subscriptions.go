package subscriptions

import (
	"github.com/stripe/stripe-go/v76"
)

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewCustomer struct {
	UserId string
	Name   string
	Email  string
}

func StripeCustomerToLocal(c *stripe.Customer) Customer {
	return Customer{
		ID:   c.ID,
		Name: c.Name,
	}
}
