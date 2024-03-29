package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/files"
)

var selectAll string = "SELECT `id`, `user_id`, `name`, `updated_at`, `created_at`, `slug`, `file_size`, `s3_region`, `s3_bucket`, `s3_key` FROM `files` "

func selectAllScan(rows *sql.Rows) (files.File, error) {
	var file files.File
	err := rows.Scan(
		&file.ID,
		&file.UserId,
		&file.Name,
		&file.UpdatedAt,
		&file.CreatedAt,
		&file.Slug,
		&file.FileSize,
		&file.S3Region,
		&file.S3Bucket,
		&file.S3Key,
	)
	return file, err
}

// Count implements files.Repository.
func (r *FileRepo) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(`id`) FROM `files`"

	results, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer results.Close()

	var count int
	if !results.Next() {
		return 0, fmt.Errorf("failed to get result")
	}

	err = results.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// DeleteById implements files.Repository.
func (r *FileRepo) DeleteById(ctx context.Context, id string) error {
	query := "DELETE FROM `files` WHERE `id` = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("unable to delete file: %w", err)
	}

	return nil
}

// GetById implements files.Repository.
func (r *FileRepo) GetById(ctx context.Context, id string) (*files.File, error) {
	query := selectAll + "WHERE `id` = ?"

	results, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	if !results.Next() {
		return nil, files.ErrFileNotFound
	}

	file, err := selectAllScan(results)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// Update implements files.Repository.
func (r *FileRepo) Create(ctx context.Context, file *files.File) error {
	query := "INSERT INTO `files` (`id`, `user_id`, `name`, `updated_at`, `created_at`, `slug`, `file_size`, `s3_region`, `s3_bucket`, `s3_key`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query,
		file.ID,
		file.UserId,
		file.Name,
		file.UpdatedAt,
		file.CreatedAt,
		file.Slug,
		file.FileSize,
		file.S3Region,
		file.S3Bucket,
		file.S3Key,
	)
	if err != nil {
		return fmt.Errorf("unable to insert file: %w", err)
	}

	return nil
}

// Update implements files.Repository.
func (r *FileRepo) Update(ctx context.Context, file *files.File) error {
	query := "UPDATE `files` SET `user_id` = ?, `name` = ?, `updated_at` = ?, `created_at` = ?, `slug` = ?, `file_size` = ?, `s3_region` = ?, `s3_bucket` = ?, `s3_key` = ? WHERE `id` = ?"

	_, err := r.db.ExecContext(ctx, query,
		file.UserId,
		file.Name,
		file.UpdatedAt,
		file.CreatedAt,
		file.Slug,
		file.FileSize,
		file.S3Region,
		file.S3Bucket,
		file.S3Key,
		file.ID,
	)
	if err != nil {
		return fmt.Errorf("unable to insert file: %w", err)
	}

	return nil
}
