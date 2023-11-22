package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LouisHatton/menu-link-up/internal/bandwidth"
	internal_context "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/LouisHatton/menu-link-up/internal/users"
	"github.com/google/uuid"
)

type BandwidthSvc struct {
	logger  *log.Logger
	repo    bandwidth.Repository
	userSvc users.Service
}

func New(logger *log.Logger, repo bandwidth.Repository, userSvc users.Service) (*BandwidthSvc, error) {
	svc := BandwidthSvc{
		logger:  logger,
		repo:    repo,
		userSvc: userSvc,
	}

	err := repo.Ping()
	if err != nil {
		return nil, err
	}

	return &svc, nil
}

// RecordDocumentUpload implements bandwidth.Service.
func (svc *BandwidthSvc) RecordDocumentUpload(ctx context.Context, userId string, fileSize int) error {
	ctxUserId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.Context(ctx), log.UserId(ctxUserId), log.RequestedId(userId))

	record, err := svc.getUsersMonthsRecord(ctx, userId)
	if err != nil {
		msg := "attempting to get users current bandwidth record"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	if record.BytesUploaded >= record.BytesUploadedLimit {
		return bandwidth.ErrUploadLimitReached
	}

	err = svc.repo.IncreaseBytesUploaded(ctx, record.ID, int64(fileSize))
	if err != nil {
		msg := "attempting to increase bytes uploaded"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	return nil
}

// RecordDocumentView implements bandwidth.Service.
func (svc *BandwidthSvc) RecordDocumentView(ctx context.Context, userId string, fileSize int) error {
	logger := svc.logger.With(log.Context(ctx), log.UserId(userId))

	record, err := svc.getUsersMonthsRecord(ctx, userId)
	if err != nil {
		msg := "attempting to get users current bandwidth record"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	if record.BytesTransferred >= record.BytesTransferredLimit {
		return bandwidth.ErrBytesTransferredLimitReached
	}

	err = svc.repo.IncreaseBytesTransferred(ctx, record.ID, int64(fileSize))
	if err != nil {
		msg := "attempting to increase bytes transferred"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	return nil
}

func (svc *BandwidthSvc) getUsersMonthsRecord(ctx context.Context, userId string) (*bandwidth.MonthlyBandwidth, error) {
	month := time.Now().Month()
	year := time.Now().Year()

	var record *bandwidth.MonthlyBandwidth

	record, err := svc.repo.GetByUserIdMonthYear(ctx, userId, int(month), year)
	switch err {
	case nil:
	case bandwidth.ErrRecordNotFound:
		record, err = svc.createUsersMonthsRecord(ctx, userId)
		if err != nil {
			return nil, fmt.Errorf("encountered error attempting create new monthly record: %w", err)
		}
	default:
		return nil, fmt.Errorf("encountered error attempting to get users current monthly record: %w", err)
	}

	return record, nil
}

func (svc *BandwidthSvc) createUsersMonthsRecord(ctx context.Context, userId string) (*bandwidth.MonthlyBandwidth, error) {
	month := time.Now().Month()
	year := time.Now().Year()

	bandwidthLimits, err := svc.userSvc.GetBandwidthLimits(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("attempting to get users bandwidth limits to create new record: %w", err)
	}

	newRecord := bandwidth.MonthlyBandwidth{
		ID:                    uuid.NewString(),
		UserId:                userId,
		Month:                 int(month),
		Year:                  year,
		BytesTransferred:      0,
		BytesUploaded:         0,
		BytesTransferredLimit: bandwidthLimits.BytesTransferredLimit,
		BytesUploadedLimit:    bandwidthLimits.BytesUploadedLimit,
	}

	err = svc.repo.Create(ctx, &newRecord)
	if err != nil {
		return nil, err
	}

	return &newRecord, nil
}

func (svc *BandwidthSvc) DeleteAllUserRecords(ctx context.Context, userId string) error {
	return svc.repo.DeleteByUserId(ctx, userId)
}

var _ bandwidth.Service = &BandwidthSvc{}
