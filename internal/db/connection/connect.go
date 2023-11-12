package connection

import (
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	"github.com/LouisHatton/menu-link-up/internal/db/migrate"
	"github.com/LouisHatton/menu-link-up/internal/log"
)

func Connect(logger *log.Logger, cfg *appconfig.Database) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.UserName, cfg.Password, cfg.Host, cfg.Name))
	if err != nil {
		return nil, err
	}

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	logger.Info(fmt.Sprintf("Connected to database: '%s' (%s), at %s as %s", cfg.Name, version, cfg.Host, cfg.UserName))

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = migrate.RunMigrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
