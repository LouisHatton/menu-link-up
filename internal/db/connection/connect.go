package connection

import (
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	"github.com/LouisHatton/menu-link-up/internal/db/migrate"
)

func Connect(cfg *appconfig.Database) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.UserName, cfg.Password, cfg.Host, cfg.Name))
	if err != nil {
		return nil, err
	}

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Printf("Connected to database: '%s' (%s), at %s as %s\n", cfg.Name, version, cfg.Host, cfg.UserName)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	db.Ping()

	err = migrate.RunMigrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
