package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/files"
)

type FileRepo struct {
	db *sql.DB
}

var _ files.Repository = &FileRepo{}

func New(db *sql.DB) *FileRepo {
	return &FileRepo{
		db: db,
	}
}

func (r *FileRepo) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return err
	}

	_, err = r.db.Exec("select 1 from `files` limit 1")
	if err != nil {
		return fmt.Errorf("attempting to check if table `files` exists: %w", err)
	}

	return nil
}

// GetBySlug implements files.Repository.
func (r *FileRepo) GetBySlug(ctx context.Context, slug string) (*files.File, error) {
	query := selectAll + "WHERE `slug` = ?"

	results, err := r.db.QueryContext(ctx, query, slug)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, fmt.Errorf("failed to get file by slug")
	}

	file, err := selectAllScan(results)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetByUserId implements files.Repository.
func (r *FileRepo) GetByUserId(ctx context.Context, userId string) (*[]files.File, error) {
	query := selectAll + "WHERE `user_id` = ? ORDER BY `created_at` DESC"

	results, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var output []files.File = []files.File{}
	for results.Next() {
		file, err := selectAllScan(results)
		if err != nil {
			return nil, err
		}

		output = append(output, file)
	}

	return &output, nil
}
