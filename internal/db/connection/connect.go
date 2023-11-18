package connection

import (
	"database/sql"
	"fmt"

	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	"github.com/LouisHatton/menu-link-up/internal/db/migrate"
	"github.com/LouisHatton/menu-link-up/internal/log"
)

func Connect(logger *log.Logger, cfg *appconfig.Database) (*sql.DB, error) {
	logger = logger.With(log.String("dbName", cfg.Name), log.String("dbHost", cfg.Host), log.String("dbUsername", cfg.UserName))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.UserName, cfg.Password, cfg.Host, cfg.Name))
	if err != nil {
		logger.Error("error opening mysql connection", log.Error(err))
		return nil, err
	}

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	logger.Info(fmt.Sprintf("Connected to database (Version: %s)", version))

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	err = db.Ping()
	if err != nil {
		logger.Error("error pinging mysql db", log.Error(err))
		return nil, err
	}

	err = migrate.RunMigrate(db)
	if err != nil {
		logger.Error("error running db migration", log.Error(err))
		return nil, err
	}

	return db, nil
}
