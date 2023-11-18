package repository

import (
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/users"
)

type UserRepo struct {
	db *sql.DB
}

var _ users.Repository = &UserRepo{}

func New(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return err
	}

	_, err = r.db.Exec("select 1 from `users` limit 1")
	if err != nil {
		return fmt.Errorf("attempting to check if table `users` exists: %w", err)
	}

	return nil
}
