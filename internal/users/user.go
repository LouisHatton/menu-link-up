package users

import (
	"firebase.google.com/go/v4/auth"
)

type User struct {
	ID            string `json:"id"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	ProviderID    string `json:"providerId"`
}

func AuthUserRecordToUser(user *auth.UserRecord) User {
	return User{
		ID:            user.UID,
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		ProviderID:    user.ProviderID,
	}
}
