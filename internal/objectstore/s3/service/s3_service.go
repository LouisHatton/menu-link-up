package service

import (
	"context"
	"fmt"
	"time"

	internal_context "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/LouisHatton/menu-link-up/internal/objectstore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type S3Service struct {
	logger        *log.Logger
	s3            *s3.S3
	region        string
	defaultBucket string
}

func New(logger *log.Logger, s *session.Session, defaultBucket string) (*S3Service, error) {
	svc := s3.New(s)

	return &S3Service{
		logger:        logger,
		s3:            svc,
		region:        *s.Config.Region,
		defaultBucket: defaultBucket,
	}, nil
}

// GenerateFileLocation implements objectstore.Service.
func (svc *S3Service) GenerateFileLocation(ctx context.Context) (objectstore.FileLocation, error) {
	var location objectstore.FileLocation
	var key string
	var err error
	var exists bool
	attempts := 0

	keyFound := false
	for !keyFound && attempts < 3 {
		attempts++
		key = uuid.NewString() + ".pdf"
		exists, err = svc.keyExists(svc.defaultBucket, key)
		if err != nil {
			msg := "attempting to check if key exists"
			err = fmt.Errorf(msg+": %w", err)
			svc.logger.Warn(msg, zap.Error(err))
		} else {
			keyFound = !exists
		}
	}

	if !keyFound {
		return location, err
	}

	location.Region = svc.region
	location.Bucket = svc.defaultBucket
	location.Key = key

	return location, err
}

func (svc *S3Service) keyExists(bucket string, key string) (bool, error) {
	_, err := svc.s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if a_err, ok := err.(awserr.Error); ok {
			switch a_err.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

// PresignedGet implements objectstore.Service.
func (svc *S3Service) PresignedGet(ctx context.Context, location objectstore.FileLocation, expires time.Duration) (string, error) {
	err := svc.parseFileLocation(&location)
	if err != nil {
		return "", err
	}

	req, _ := svc.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(location.Bucket),
		Key:    aws.String(location.Key),
	})

	str, err := req.Presign(expires)
	if err != nil {
		return "", fmt.Errorf("attempting to pre-sign get object request: %w", err)
	}

	svc.logger.Info("generated presigned get for object", log.FileLocation(location), log.String("operation", "GET"))
	return str, nil
}

// PresignedPut implements objectstore.Service.
func (svc *S3Service) PresignedPut(ctx context.Context, location objectstore.FileLocation, fileSize int, expires time.Duration) (string, error) {
	userId := internal_context.GetUserIdFromContext(ctx)
	err := svc.parseFileLocation(&location)
	if err != nil {
		return "", err
	}

	req, _ := svc.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(location.Bucket),
		Key:    aws.String(location.Key),
	})

	req.HTTPRequest.Header.Set("Content-Type", "application/pdf")
	req.HTTPRequest.Header.Set("Content-Length", fmt.Sprint(fileSize))
	str, err := req.Presign(expires)
	if err != nil {
		return "", fmt.Errorf("attempting to pre-sign put object request: %w", err)
	}

	svc.logger.Info("generated presigned put for object", log.FileLocation(location), log.UserId(userId), log.String("operation", "PUT"))
	return str, nil
}

func (svc *S3Service) parseFileLocation(location *objectstore.FileLocation) error {
	if location.Region == "" {
		location.Region = svc.region
	}

	if location.Bucket == "" {
		location.Bucket = svc.defaultBucket
	}

	if location.Region != svc.region {
		return objectstore.ErrIncorrectRegion
	}

	return nil
}

func (svc *S3Service) DeleteFile(ctx context.Context, location objectstore.FileLocation) error {
	userId := internal_context.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.UserId(userId), log.FileLocation(location), log.String("operation", "DELETE"), log.Context(ctx))

	_, err := svc.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(location.Bucket),
		Key:    aws.String(location.Key),
	})
	if err != nil {
		logger.Error("attempting to delete s3 file: %s", log.Error(err))
		return err
	}

	logger.Info("deleted object")
	return nil
}

var _ objectstore.Service = &S3Service{}
