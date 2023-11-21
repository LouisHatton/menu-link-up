package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/users"
)

var selectAll string = "SELECT `id`, `email`, `stripe_customer_id`, `stripe_subscription_id`, `subscription_status`, `bytes_transferred_limit`, `bytes_uploaded_limit`, `file_size_limit`, `file_upload_limit` FROM `users` "

func selectAllScan(rows *sql.Rows) (users.User, error) {
	var user users.User
	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.StripeCustomerId,
		&user.StripeSubscriptionId,
		&user.SubscriptionStatus,
		&user.BytesTransferredLimit,
		&user.BytesUploadedLimit,
		&user.FileSizeLimit,
		&user.FileUploadLimit,
	)
	return user, err
}

// Count implements users.Repository.
func (r *UserRepo) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(`id`) FROM `users`"

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

// DeleteById implements users.Repository.
func (r *UserRepo) DeleteById(ctx context.Context, id string) error {
	query := "DELETE FROM `users` WHERE `id` = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}

// GetById implements users.Repository.
func (r *UserRepo) GetById(ctx context.Context, id string) (*users.User, error) {
	query := selectAll + "WHERE `id` = ?"

	results, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	if !results.Next() {
		return nil, users.ErrUserNotFound
	}

	user, err := selectAllScan(results)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update implements users.Repository.
func (r *UserRepo) Create(ctx context.Context, user *users.User) error {
	query := "INSERT INTO `users` (`id`, `email`, `stripe_customer_id`, `stripe_subscription_id`, `subscription_status`, `bytes_transferred_limit`, `bytes_uploaded_limit`, `file_size_limit`, `file_upload_limit`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.StripeCustomerId,
		user.StripeSubscriptionId,
		user.SubscriptionStatus,
		user.BytesTransferredLimit,
		user.BytesUploadedLimit,
		user.FileSizeLimit,
		user.FileUploadLimit,
	)
	if err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}

	return nil
}

// Update implements users.Repository.
func (r *UserRepo) Update(ctx context.Context, user *users.User) error {
	query := "UPDATE `users` SET `email` = ?, `stripe_customer_id` = ?,  `stripe_subscription_id` = ?, `subscription_status` = ?,  `bytes_transferred_limit` = ?, `bytes_uploaded_limit` = ?, `file_size_limit` = ?, `file_upload_limit` = ? WHERE `id` = ?"

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.StripeCustomerId,
		user.StripeSubscriptionId,
		user.SubscriptionStatus,
		user.BytesTransferredLimit,
		user.BytesUploadedLimit,
		user.FileSizeLimit,
		user.FileUploadLimit,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}

	return nil
}
