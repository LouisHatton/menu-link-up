package service

import (
	"context"

	"github.com/LouisHatton/menu-link-up/internal/files"
)

type FileSvc struct {
	repo files.Repository
}

// Create implements files.Service.
func (*FileSvc) Create(ctx context.Context, newFile files.NewFile) (*files.File, error) {
	panic("unimplemented")
}

// Delete implements files.Service.
func (*FileSvc) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// Edit implements files.Service.
func (*FileSvc) Edit(ctx context.Context, id string, newFile files.NewFile) error {
	panic("unimplemented")
}

// GetById implements files.Service.
func (svc *FileSvc) GetById(ctx context.Context, id string) (*files.File, error) {
	return svc.repo.GetById(ctx, id)
}

// GetByUserId implements files.Service.
func (svc *FileSvc) GetByUserId(ctx context.Context, userId string) (*[]files.File, error) {
	return svc.repo.GetByUserId(ctx, userId)
}

func (svc *FileSvc) GetS3LinkFromSlug(ctx context.Context, slug string) (string, error) {
	panic("unimplemented")
}

var _ files.Service = &FileSvc{}
