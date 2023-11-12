package appconfig

import "github.com/LouisHatton/menu-link-up/internal/config/environment"

type Environment struct {
	CurrentEnv environment.Type `env:"ENVIRONMENT" envDefault:"other"`
}

type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type Database struct {
	Name         string `env:"DATABASE_NAME" envDefault:"menulinkup"`
	UserName     string `env:"DATABASE_USER_NAME" envDefault:"adminer"`
	Password     string `env:"DATABASE_PASSWORD" envDefault:"adminer"`
	Host         string `env:"DATABASE_HOST" envDefault:"localhost:5555"`
	MaxIdleConns int    `env:"DATABASE_MAX_IDLE_CONNS" envDefault:"5"`
	MaxOpenConns int    `env:"DATABASE_MAX_OPEN_CONNS" envDefault:"10"`
}

type AWS_S3 struct {
	Region        string `env:"AWS_S3_REGION" envDefault:"eu-west-1"`
	DefaultBucket string `env:"AWS_S3_DEFAULT_BUCKET" envDefault:"menulink-dev-pdf-bucket"`
}
