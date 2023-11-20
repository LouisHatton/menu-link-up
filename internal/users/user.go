package users

import (
	"firebase.google.com/go/v4/auth"
)

type User struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	StripeCustomerId string `json:"stripeCustomerId"`
}

func AuthUserRecordToUser(user *auth.UserRecord) User {
	return User{
		ID:    user.UID,
		Email: user.Email,
	}
}

type BandwidthLimits struct {
	BytesTransferredLimit int64 `json:"bytesTransferredLimit"`
	BytesUploadedLimit    int64 `json:"bytesUploadedLimit"`
}
