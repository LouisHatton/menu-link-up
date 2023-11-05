package appconfig

import "github.com/LouisHatton/menu-link-up/internal/config/environment"

type Environment struct {
	CurrentEnv environment.Type `env:"ENVIRONMENT" envDefault:"other"`
}

type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}
