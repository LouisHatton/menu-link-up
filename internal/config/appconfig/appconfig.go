package appconfig

import "github.com/LouisHatton/menu-link-up/internal/config/environment"

type Environment struct {
	CurrentEnv environment.Type `env:"ENVIRONMENT" envDefault:"other"`
}

type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type Database struct {
	Name     string `env:"DATABASE_NAME" envDefault:"menulinkup"`
	UserName string `env:"DATABASE_USER_NAME" envDefault:"adminer"`
	Password string `env:"DATABASE_PASSWORD" envDefault:"adminer"`
	Host     string `env:"DATABASE_HOST" envDefault:"localhost:5555"`
}
