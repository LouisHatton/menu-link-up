package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/bandwidth"
)

var selectAll string = "SELECT `id`, `user_id`, `month`, `year`, `bytes_transferred`, `bytes_transferred_limit`, `bytes_uploaded`, `bytes_uploaded_limit` FROM `monthly_bandwidth` "

func selectAllScan(rows *sql.Rows) (bandwidth.MonthlyBandwidth, error) {
	var doc bandwidth.MonthlyBandwidth
	err := rows.Scan(
		&doc.ID,
		&doc.UserId,
		&doc.Month,
		&doc.Year,
		&doc.BytesTransferred,
		&doc.BytesTransferredLimit,
		&doc.BytesUploaded,
		&doc.BytesUploadedLimit,
	)
	return doc, err
}

// Count implements bandwidth.Repository.
func (r *BandwidthRepo) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(`id`) FROM `monthly_bandwidth`"

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

// DeleteById implements bandwidth.Repository.
func (r *BandwidthRepo) DeleteById(ctx context.Context, id string) error {
	query := "DELETE FROM `monthly_bandwidth` WHERE `id` = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("unable to delete monthly_bandwidth record: %w", err)
	}

	return nil
}

// GetById implements bandwidth.Repository.
func (r *BandwidthRepo) GetById(ctx context.Context, id string) (*bandwidth.MonthlyBandwidth, error) {
	query := selectAll + "WHERE `id` = ?"

	results, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	if !results.Next() {
		return nil, bandwidth.ErrRecordNotFound
	}

	record, err := selectAllScan(results)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// Update implements bandwidth.Repository.
func (r *BandwidthRepo) Create(ctx context.Context, record *bandwidth.MonthlyBandwidth) error {
	query := "INSERT INTO `monthly_bandwidth` (`id`, `user_id`, `month`, `year`, `bytes_transferred`, `bytes_transferred_limit`, `bytes_uploaded`, `bytes_uploaded_limit`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query,
		record.ID,
		record.UserId,
		record.Month,
		record.Year,
		record.BytesTransferred,
		record.BytesTransferredLimit,
		record.BytesUploaded,
		record.BytesUploadedLimit,
	)
	if err != nil {
		return fmt.Errorf("unable to insert monthly_bandwidth record: %w", err)
	}

	return nil
}

// Update implements bandwidth.Repository.
func (r *BandwidthRepo) Update(ctx context.Context, record *bandwidth.MonthlyBandwidth) error {
	query := "UPDATE `monthly_bandwidth` SET `user_id` = ?, `month` = ?, `year` = ?, `bytes_transferred_limit` = ?, `bytes_uploaded_limit` = ?, WHERE `id` = ?"

	_, err := r.db.ExecContext(ctx, query,
		record.UserId,
		record.Month,
		record.Year,
		record.BytesTransferredLimit,
		record.BytesUploadedLimit,
		record.ID,
	)
	if err != nil {
		return fmt.Errorf("unable to insert monthly_bandwidth record: %w", err)
	}

	return nil
}
