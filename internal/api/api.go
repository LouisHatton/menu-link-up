package api

import (
	"net/http"
	"time"

	"github.com/LouisHatton/menu-link-up/internal/api/middleware"
	"github.com/LouisHatton/menu-link-up/internal/api/routes"
	"github.com/LouisHatton/menu-link-up/internal/config/environment"
	"github.com/LouisHatton/menu-link-up/internal/files"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Config struct {
	env environment.Type
}

type API struct {
	l              *log.Logger
	config         *Config
	fileSvc        files.Service
	born           time.Time
	authMiddleware middleware.Auth
}

func New(logger *log.Logger, env environment.Type, authMiddleware *middleware.Auth, fileSvc files.Service) (*API, error) {

	cfg := &Config{
		env: env,
	}
	api := API{
		l:              logger,
		config:         cfg,
		fileSvc:        fileSvc,
		born:           time.Now(),
		authMiddleware: *authMiddleware,
	}

	return &api, nil
}

func (api API) Register(r chi.Router) error {

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, map[string]string{
			"born":        api.born.Format(time.RFC3339),
			"environment": string(api.config.env),
		})
	})

	r.Route("/v1", func(r chi.Router) {

		r.Use(api.authMiddleware.Middleware)

		r.Get(routes.FileIdPath, api.GetFile)
		r.Post(routes.FileIdPath, api.EditFile)
		r.Delete(routes.FileIdPath, api.DeleteFile)
		r.Get(routes.CreateFilesPath, api.ListFiles)
		r.Post(routes.CreateFilesPath, api.CreateFile)

	})

	return nil
}
