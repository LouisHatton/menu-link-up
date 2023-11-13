package migrate

import "database/sql"

var createUserTable string = `
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
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

func RunMigrate(db *sql.DB) error {
	_, err := db.Exec(createUserTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createFilesTable)
	return err
}
