package migrate

import "database/sql"

var createUserTable string = `
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    stripe_customer_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);`

var createFilesTable string = `
CREATE TABLE IF NOT EXISTS files (
    id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    slug VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    s3_region VARCHAR(255) NOT NULL,
    s3_bucket VARCHAR(255) NOT NULL,
    s3_key VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);`

var createBandwidthTable string = `
CREATE TABLE IF NOT EXISTS monthly_bandwidth (
    id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    bytes_transferred BIGINT NOT NULL,
    bytes_transferred_limit BIGINT NOT NULL,
    bytes_uploaded BIGINT NOT NULL,
    bytes_uploaded_limit BIGINT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);`

var createBandwidthIndexes string = `
CREATE INDEX IF NOT EXISTS monthly_bandwidth_user_year_month_index ON monthly_bandwidth (user_id, year, month);`

func RunMigrate(db *sql.DB) []error {
	var commands []*string = []*string{
		&createUserTable,
		&createFilesTable,
		&createBandwidthTable,
		&createBandwidthIndexes,
	}
	var errors []error

	for _, command := range commands {
		_, err := db.Exec(*command)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
