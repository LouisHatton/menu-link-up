package users

import (
	"firebase.google.com/go/v4/auth"
	"github.com/stripe/stripe-go/v76"
)

type User struct {
	ID                    string                    `json:"id"`
	Email                 string                    `json:"email"`
	StripeCustomerId      string                    `json:"stripeCustomerId"`
	StripeSubscriptionId  string                    `json:"stripeSubscriptionId"`
	SubscriptionStatus    stripe.SubscriptionStatus `json:"subscriptionStatus"`
	BytesTransferredLimit int64                     `json:"-"`
	BytesUploadedLimit    int64                     `json:"-"`
	FileSizeLimit         int64                     `json:"fileSizeLimit"`
	FileUploadLimit       int                       `json:"fileUploadLimit"`
}

func AuthUserRecordToUser(user *auth.UserRecord) User {
	return User{
		ID:    user.UID,
		Email: user.Email,
	}
}

func (user *User) AddLimits(limits *BandwidthLimits) {
	user.BytesTransferredLimit = limits.BytesTransferredLimit
	user.BytesUploadedLimit = limits.BytesUploadedLimit
	user.FileSizeLimit = limits.FileSizeLimit
	user.FileUploadLimit = limits.FileUploadLimit
}

type BandwidthLimits struct {
	BytesTransferredLimit int64
	BytesUploadedLimit    int64
	FileSizeLimit         int64
	FileUploadLimit       int
}
