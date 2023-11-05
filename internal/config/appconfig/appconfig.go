package appconfig

import "github.com/LouisHatton/menu-link-up/internal/config/environment"

type Environment struct {
	CurrentEnv environment.Type `env:"ENVIRONMENT" envDefault:"other"`
}

type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type Database struct {
	ProjectsCollectionName string `env:"PROJECTS_COLLECTION_NAME" envDefault:"projects"`
	FilesCollectionName    string `env:"FILES_COLLECTION_NAME" envDefault:"files"`
}
