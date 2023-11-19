package service

import (
	"context"
	"fmt"
	"time"

	internalContext "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/files"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/LouisHatton/menu-link-up/internal/objectstore"
	"github.com/google/uuid"
)

type FileSvc struct {
	fileRepo    files.Repository
	objStoreSvc objectstore.Service
	logger      *log.Logger
}

func New(logger *log.Logger, fileRepo files.Repository, objStoreSvc objectstore.Service) (*FileSvc, error) {
	err := fileRepo.Ping()
	if err != nil {
		return nil, err
	}

	return &FileSvc{
		fileRepo:    fileRepo,
		logger:      logger,
		objStoreSvc: objStoreSvc,
	}, nil
}

// Create implements files.Service.
func (svc *FileSvc) Create(ctx context.Context, userId string, newFile files.NewFile) (*files.FileUpload, error) {
	logger := svc.logger.With(log.UserId(userId), log.Context(ctx))

	exists := svc.ExistsWithSlug(ctx, newFile.Slug)
	if exists {
		return nil, files.ErrSlugAlreadyInUse
	}

	location, err := svc.objStoreSvc.GenerateFileLocation(ctx, "u/"+userId, &newFile.FileName)
	if err != nil {
		msg := "attempting to generate objectstore location"
		logger.Error(msg, log.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	url, err := svc.objStoreSvc.PresignedPut(ctx, location, newFile.FileSize, 15*time.Minute)
	if err != nil {
		msg := "attempting to create put url"
		logger.Error(msg, log.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}

	id := uuid.NewString()
	file := files.File{
		ID:        id,
		UserId:    userId,
		Name:      newFile.Name,
		Slug:      newFile.Slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FileSize:  newFile.FileSize,
		S3Region:  location.Region,
		S3Bucket:  location.Bucket,
		S3Key:     location.Key,
	}

	err = svc.fileRepo.Create(ctx, &file)
	if err != nil {
		msg := "failed to create new file"
		logger.Error(msg, log.Error(err))
		return nil, fmt.Errorf(msg+": %w", err)
	}
	logger.Info("new file created", log.FileId(id))

	return &files.FileUpload{
		Url: url,
	}, nil
}

// Delete implements files.Service.
func (svc *FileSvc) Delete(ctx context.Context, id string) error {
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.UserId(userId), log.FileId(id), log.Context(ctx))

	file, err := svc.fileRepo.GetById(ctx, id)
	if err != nil {
		logger.Error("unable to get file to delete", log.Error(err))
		return err
	}

	if file.UserId != userId {
		logger.Warn("user attempting to delete file they do not own")
		return files.ErrNotUsersFile
	}

	err = svc.fileRepo.DeleteById(ctx, id)
	if err != nil {
		msg := "failed to delete file"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	err = svc.objStoreSvc.DeleteFile(ctx, objectstore.FileLocation{
		Bucket: file.S3Bucket,
		Region: file.S3Region,
		Key:    file.S3Key,
	})
	if err != nil {
		logger.Warn("failed to delete file from object store", log.Error(err))
	}

	logger.Info("file deleted")
	return nil
}

func (svc *FileSvc) DeleteByUserId(ctx context.Context, userId string) error {
	requestingUserId := internalContext.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.UserId(requestingUserId), log.Context(ctx), log.RequestedId(userId))

	files, err := svc.GetByUserId(ctx, userId)
	if err != nil {
		msg := "unable to get users files to delete"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	for _, file := range *files {
		logger := logger.With(log.FileId(file.ID))

		err := svc.Delete(ctx, file.ID)
		if err != nil {
			logger.Warn("unable to delete file, continuing to not block other deletes")
		}
	}

	return nil
}

// Edit implements files.Service.
func (svc *FileSvc) Edit(ctx context.Context, id string, newFile files.NewFile) error {
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := svc.logger.With(log.UserId(userId), log.FileId(id), log.Context(ctx))

	file, err := svc.fileRepo.GetById(ctx, id)
	if err != nil {
		logger.Error("unable to get file to edit", log.Error(err))
		return err
	}

	if file.UserId != userId {
		logger.Warn("user attempting to edit file they do not own")
		return files.ErrNotUsersFile
	}

	file.Name = newFile.Name
	file.Slug = newFile.Slug

	err = svc.fileRepo.Update(ctx, file)
	if err != nil {
		msg := "failed to update file"
		logger.Error(msg, log.Error(err))
		return fmt.Errorf(msg+": %w", err)
	}

	logger.Info("file updated")
	return nil
}

// GetById implements files.Service.
func (svc *FileSvc) GetById(ctx context.Context, id string) (*files.File, error) {
	return svc.fileRepo.GetById(ctx, id)
}

// GetByUserId implements files.Service.
func (svc *FileSvc) GetByUserId(ctx context.Context, userId string) (*[]files.File, error) {
	return svc.fileRepo.GetByUserId(ctx, userId)
}

// GetByUserId implements files.Service.
func (svc *FileSvc) ExistsWithSlug(ctx context.Context, slug string) bool {
	_, err := svc.fileRepo.GetBySlug(ctx, slug)
	if err == nil { // No error - file found
		return true
	}
	return false
}

func (svc *FileSvc) GetLinkFromSlug(ctx context.Context, slug string) (string, error) {
	logger := svc.logger.With(log.Context(ctx))
	msg := "attempting to generate presigned get url for file"

	file, err := svc.fileRepo.GetBySlug(ctx, slug)
	if err != nil {
		logger.Warn(msg, log.Error(err))
		return "", err
	}

	logger = logger.With(log.FileId(file.ID))

	url, err := svc.objStoreSvc.PresignedGet(
		ctx,
		objectstore.FileLocation{
			Region: file.S3Region,
			Bucket: file.S3Bucket,
			Key:    file.S3Key,
		},
		time.Minute*30,
	)
	if err != nil {
		logger.Error(msg, log.Error(err))
		return "", fmt.Errorf(msg+": %w", err)
	}

	return url, nil
}

var _ files.Service = &FileSvc{}
