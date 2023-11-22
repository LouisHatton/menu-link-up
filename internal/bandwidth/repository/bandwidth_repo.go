package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/bandwidth"
)

type BandwidthRepo struct {
	db *sql.DB
	tx *sql.Tx
}

var _ bandwidth.Repository = &BandwidthRepo{}

func New(db *sql.DB) *BandwidthRepo {
	return &BandwidthRepo{
		db: db,
		tx: nil,
	}
}

func (r *BandwidthRepo) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return err
	}

	_, err = r.db.Exec("select 1 from `monthly_bandwidth` limit 1")
	if err != nil {
		return fmt.Errorf("attempting to check if table `monthly_bandwidth` exists: %w", err)
	}

	return nil
}

func (r *BandwidthRepo) GetByUserIdMonthYear(ctx context.Context, userId string, month int, year int) (*bandwidth.MonthlyBandwidth, error) {
	query := selectAll + "WHERE `user_id` = ? AND `year` = ? AND `month` = ?"

	results, err := r.db.QueryContext(ctx, query, userId, year, month)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	ok := results.Next()
	if !ok {
		return nil, bandwidth.ErrRecordNotFound
	}

	record, err := selectAllScan(results)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *BandwidthRepo) IncreaseBytesTransferred(ctx context.Context, id string, bytesTransferred int64) error {
	query := "UPDATE `monthly_bandwidth` SET `bytes_transferred` = `bytes_transferred` + ? WHERE `id` = ?"

	_, err := r.db.Exec(query, bytesTransferred, id)
	if err != nil {
		return fmt.Errorf("attempting to update `bytes_transferred`: %w", err)
	}

	return nil
}

func (r *BandwidthRepo) IncreaseBytesUploaded(ctx context.Context, id string, bytesTransferred int64) error {
	query := "UPDATE `monthly_bandwidth` SET `bytes_uploaded` = `bytes_uploaded` + ? WHERE `id` = ?"

	_, err := r.db.Exec(query, bytesTransferred, id)
	if err != nil {
		return fmt.Errorf("attempting to update `bytes_transferred`: %w", err)
	}

	return nil
}

func (r *BandwidthRepo) DeleteByUserId(ctx context.Context, userId string) error {
	query := "DELETE FROM `monthly_bandwidth` WHERE `user_id` = ?"

	_, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("unable to delete all monthly_bandwidth records for user: %w", err)
	}

	return nil
}
