package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/LouisHatton/menu-link-up/internal/api/responses"
	internalContext "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/users"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Auth struct {
	client *auth.Client
	logger *zap.Logger
}

func NewAuth(l *zap.Logger, client *auth.Client) (*Auth, error) {
	return &Auth{
		client: client,
		logger: l,
	}, nil
}

func (m *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		providedToken, err := extractBearerToken(r)
		if err != nil {
			m.logger.Info("failed to extract bearer token", zap.Error(err))
			render.Render(w, r, responses.ErrUnauthorised())
			return
		}

		token, err := m.client.VerifyIDToken(ctx, providedToken)
		if err != nil {
			m.logger.Info("token provided is invalid", zap.Error(err))
			render.Render(w, r, responses.ErrUnauthorised())
			return
		}

		userId := token.UID
		logger := m.logger.With(zap.String("userId", userId))

		authUser, err := m.client.GetUser(ctx, userId)
		if err != nil {
			logger.Info("could not get user from id in token", zap.Error(err))
			render.Render(w, r, responses.ErrInternalServerError())
			return
		}

		realUser := users.AuthUserRecordToUser(authUser)
		ctx = internalContext.AddUserToContext(ctx, realUser)

		next.ServeHTTP(w, r.Clone(ctx))
	})
}

func extractBearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("no Authorization header provided")
	}

	splitHeader := strings.Split(header, " ")
	if len(splitHeader) < 1 {
		return "", fmt.Errorf("invalid Authorization header")
	}
	return splitHeader[1], nil
}
