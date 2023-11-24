package service

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	internal_context "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/LouisHatton/menu-link-up/internal/subscriptions"
	"github.com/LouisHatton/menu-link-up/internal/users"
	"go.uber.org/zap"
)

type UserService struct {
	logger          *log.Logger
	client          *auth.Client
	repo            users.Repository
	subscriptionSvc subscriptions.Service
}

var _ users.Service = &UserService{}

func New(l *log.Logger, authClient *auth.Client, userRepo users.Repository, subscriptionSvc subscriptions.Service) (*UserService, error) {
	svc := UserService{
		logger:          l,
		client:          authClient,
		repo:            userRepo,
		subscriptionSvc: subscriptionSvc,
	}

	err := svc.repo.Ping()
	if err != nil {
		return nil, fmt.Errorf("checking userRepo db connection: %w", err)
	}

	return &svc, nil
}

// GetById implements users.Service.
func (svc *UserService) GetById(ctx context.Context, id string) (*users.User, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(id))

	var repoUser *users.User
	repoUser, err := svc.repo.GetById(ctx, id)
	switch err {
	case nil:
	case users.ErrUserNotFound:
		repoUser, err = svc.createClientUserInRepo(ctx, id)
		msg := "attempting to create the client user in the repository"
		switch err {
		case nil:
		case users.ErrEmailNotVerified:
			logger.Error(msg, zap.Error(err))
			return nil, err
		default:
			logger.Error(msg, zap.Error(err))
			return nil, fmt.Errorf(msg+": %w", err)
		}
	default:
		msg := "attempting to fetch user from repository"
		logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)

	}

	return repoUser, nil
}

// DeleteById implements users.Service.
func (svc *UserService) DeleteById(ctx context.Context, id string) error {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(id))

	logger.Info("deleting user with requested id")

	err := svc.client.DeleteUser(ctx, id)
	if err != nil {
		msg := "attempting to delete user from client"
		logger.Error(msg, zap.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	err = svc.repo.DeleteById(ctx, id)
	if err != nil {
		msg := "attempting to delete user from repo"
		logger.Error(msg, zap.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	return nil
}

func (svc *UserService) createClientUserInRepo(ctx context.Context, id string) (*users.User, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(id))

	logger.Info("creating new user from firebase auth client")

	authRecord, err := svc.client.GetUser(ctx, id)
	if err != nil {
		msg := "attempting to fetch user from client"
		logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	if !authRecord.EmailVerified {
		return nil, users.ErrEmailNotVerified
	}

	user := users.AuthUserRecordToUser(authRecord)

	customer, subscription, err := svc.subscriptionSvc.CreateCustomerWithTrial(ctx, subscriptions.NewCustomer{
		UserId: user.ID,
		Name:   authRecord.DisplayName,
		Email:  user.Email,
	})
	if err != nil {
		msg := "attempting to create customer with trial"
		logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	user.StripeCustomerId = customer.ID
	user.StripeSubscriptionId = subscription.ID
	user.SubscriptionStatus = subscription.Status

	limits, err := svc.subscriptionSvc.GetLimitsForSubscription(ctx, subscription.ID)
	if err != nil {
		msg := "attempting to get limits for trial subscription"
		logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}
	user.AddLimits(limits)

	err = svc.repo.Create(ctx, &user)
	if err != nil {
		msg := "attempting to store user from client into repository"
		logger.Error(msg, zap.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	return &user, nil
}

func (svc *UserService) GetBandwidthLimits(ctx context.Context, id string) (*users.BandwidthLimits, error) {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(id))

	user, err := svc.GetById(ctx, id)
	if err != nil {
		logger.Error("attempting to get user to return bandwidth limits", log.Error(err))
		return nil, err
	}

	return &users.BandwidthLimits{
		BytesTransferredLimit: user.BytesTransferredLimit,
		BytesUploadedLimit:    user.BytesUploadedLimit,
	}, nil
}
